/*
	Copyright NetFoundry Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package xgress_edge

import (
	"encoding/binary"
	"fmt"
	"github.com/openziti/ziti/common/capabilities"
	"github.com/openziti/ziti/common/cert"
	fabricMetrics "github.com/openziti/ziti/common/metrics"
	"github.com/openziti/ziti/common/pb/ctrl_pb"
	"github.com/openziti/ziti/common/pb/edge_ctrl_pb"
	"github.com/openziti/ziti/controller/idgen"
	"github.com/openziti/ziti/router/state"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"math/big"
	"strings"
	"time"

	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/channel/v2"
	"github.com/openziti/channel/v2/protobufs"
	"github.com/openziti/identity"
	"github.com/openziti/sdk-golang/ziti/edge"
	"github.com/openziti/transport/v2"
	"github.com/openziti/ziti/router/xgress"
	"github.com/openziti/ziti/router/xgress_common"
)

type listener struct {
	id               *identity.TokenId
	factory          *Factory
	options          *Options
	bindHandler      xgress.BindHandler
	underlayListener channel.UnderlayListener
	headers          map[int32][]byte
}

// newListener creates a new xgress edge listener
func newListener(id *identity.TokenId, factory *Factory, options *Options, headers map[int32][]byte) xgress.Listener {
	return &listener{
		id:      id,
		factory: factory,
		options: options,
		headers: headers,
	}
}

func (listener *listener) Listen(address string, bindHandler xgress.BindHandler) error {
	if address == "" {
		return errors.New("address must be specified for edge listeners")
	}
	listener.bindHandler = bindHandler
	addr, err := transport.ParseAddress(address)

	if err != nil {
		return fmt.Errorf("cannot listen on invalid address [%s] (%s)", address, err)
	}

	pfxlog.Logger().WithField("address", addr).Info("starting channel listener")

	listenerConfig := channel.ListenerConfig{
		ConnectOptions:   listener.options.channelOptions.ConnectOptions,
		TransportConfig:  listener.factory.edgeRouterConfig.Tcfg,
		Headers:          listener.headers,
		PoolConfigurator: fabricMetrics.GoroutinesPoolMetricsConfigF(listener.factory.metricsRegistry, "pool.listener.xgress_edge"),
	}

	listener.underlayListener = channel.NewClassicListener(listener.id, addr, listenerConfig)

	if err := listener.underlayListener.Listen(); err != nil {
		return err
	}
	accepter := NewAcceptor(listener, listener.underlayListener, nil)
	go accepter.Run()

	return nil
}

func (listener *listener) Close() error {
	return listener.underlayListener.Close()
}

type edgeClientConn struct {
	msgMux       edge.MsgMux
	listener     *listener
	fingerprints cert.Fingerprints
	ch           channel.Channel
	idSeq        uint32
}

func (self *edgeClientConn) HandleClose(_ channel.Channel) {
	log := pfxlog.ContextLogger(self.ch.Label())
	log.Debugf("closing")
	self.listener.factory.hostedServices.cleanupServices(self)
	self.msgMux.Close()
}

func (self *edgeClientConn) ContentType() int32 {
	return edge.ContentTypeData
}

func (self *edgeClientConn) processConnect(manager state.Manager, req *channel.Message, ch channel.Channel) {
	sessionToken := string(req.Body)
	log := pfxlog.ContextLogger(ch.Label()).WithField("token", sessionToken).WithFields(edge.GetLoggerFields(req))
	connId, found := req.GetUint32Header(edge.ConnIdHeader)
	if !found {
		pfxlog.Logger().Errorf("connId not set. unable to process connect message")
		return
	}

	ctrlCh := self.listener.factory.ctrls.AnyCtrlChannel()
	if ctrlCh == nil {
		errStr := "no controller available, cannot create circuit"
		log.Error(errStr)
		self.sendStateClosedReply(errStr, req)
		return
	}

	conn := &edgeXgressConn{
		mux:        self.msgMux,
		MsgChannel: *edge.NewEdgeMsgChannel(self.ch, connId),
		seq:        NewMsgQueue(4),
	}

	// need to remove session remove listener on close
	conn.onClose = self.listener.factory.stateManager.AddEdgeSessionRemovedListener(sessionToken, func(token string) {
		conn.close(true, "session closed")
	})

	// We can't fix conn id, since it's provided by the client
	if err := self.msgMux.AddMsgSink(conn); err != nil {
		log.WithError(err).WithField("token", sessionToken).Error("error adding to msg mux")
		self.sendStateClosedReply(err.Error(), req)
		return
	}

	// fabric connect
	log.Debug("dialing fabric")
	peerData := make(map[uint32][]byte)

	for _, key := range []uint32{edge.PublicKeyHeader, edge.CallerIdHeader, edge.AppDataHeader} {
		if pk, found := req.Headers[int32(key)]; found {
			peerData[key] = pk
		}
	}

	terminatorIdentity, _ := req.GetStringHeader(edge.TerminatorIdentityHeader)

	request := &edge_ctrl_pb.CreateCircuitRequest{
		SessionToken:         sessionToken,
		Fingerprints:         self.fingerprints.Prints(),
		TerminatorInstanceId: terminatorIdentity,
		PeerData:             peerData,
	}

	if strings.HasPrefix(sessionToken, JwtTokenPrefix) {
		apiSession := manager.GetApiSessionFromCh(ch)

		if apiSession == nil {
			pfxlog.Logger().Errorf("could not find api session for channel, unable to process bind message")
			return
		}

		request.ApiSessionToken = apiSession.Token
	}

	response := &edge_ctrl_pb.CreateCircuitResponse{}
	timeout := self.listener.options.Options.GetCircuitTimeout
	responseMsg, err := protobufs.MarshalTyped(request).WithTimeout(timeout).SendForReply(ctrlCh)

	if err = getResultOrFailure(responseMsg, err, response); err != nil {
		log.WithError(err).Warn("failed to dial fabric")
		self.sendStateClosedReply(err.Error(), req)
		conn.close(false, "failed to dial fabric")
		return
	}

	x := xgress.NewXgress(response.CircuitId, ctrlCh.Id(), xgress.Address(response.Address), conn, xgress.Initiator, &self.listener.options.Options, response.Tags)
	self.listener.bindHandler.HandleXgressBind(x)
	conn.ctrlRx = x
	// send the state_connected before starting the xgress. That way we can't get a state_closed before we get state_connected
	self.sendStateConnectedReply(req, response.PeerData)
	x.Start()
}

func (self *edgeClientConn) processBind(manager state.Manager, req *channel.Message, ch channel.Channel) {
	ctrlCh := self.listener.factory.ctrls.AnyCtrlChannel()
	if ctrlCh == nil {
		errStr := "no controller available, cannot create terminator"
		pfxlog.ContextLogger(ch.Label()).
			WithField("sessionToken", string(req.Body)).
			WithFields(edge.GetLoggerFields(req)).
			WithField("routerId", self.listener.id.Token).
			Error(errStr)
		self.sendStateClosedReply(errStr, req)
		return
	}

	supportsCreateTerminatorV2 := false
	headers := ctrlCh.Underlay().Headers()
	if val, found := headers[int32(ctrl_pb.ContentType_CapabilitiesHeader)]; found {
		capabilitiesMask := &big.Int{}
		capabilitiesMask.SetBytes(val)
		supportsCreateTerminatorV2 = capabilitiesMask.Bit(capabilities.ControllerCreateTerminatorV2) == 1
	}

	if supportsCreateTerminatorV2 {
		self.processBindV2(manager, req, ch, ctrlCh)
	} else {
		self.processBindV1(manager, req, ch, ctrlCh)
	}
}

func (self *edgeClientConn) processBindV1(manager state.Manager, req *channel.Message, ch channel.Channel, ctrlCh channel.Channel) {
	sessionToken := string(req.Body)

	log := pfxlog.ContextLogger(ch.Label()).
		WithField("sessionToken", sessionToken).
		WithFields(edge.GetLoggerFields(req)).
		WithField("routerId", self.listener.id.Token)

	connId, found := req.GetUint32Header(edge.ConnIdHeader)
	if !found {
		pfxlog.Logger().Errorf("connId not set. unable to process bind message")
		return
	}

	log.Debug("binding service")

	hostData := make(map[uint32][]byte)
	pubKey, hasKey := req.Headers[edge.PublicKeyHeader]
	if hasKey {
		hostData[edge.PublicKeyHeader] = pubKey
	}

	cost := uint16(0)
	if costBytes, hasCost := req.Headers[edge.CostHeader]; hasCost {
		cost = binary.LittleEndian.Uint16(costBytes)
	}

	precedence := edge_ctrl_pb.TerminatorPrecedence_Default
	if precedenceData, hasPrecedence := req.Headers[edge.PrecedenceHeader]; hasPrecedence && len(precedenceData) > 0 {
		edgePrecedence := edge.Precedence(precedenceData[0])
		if edgePrecedence == edge.PrecedenceRequired {
			precedence = edge_ctrl_pb.TerminatorPrecedence_Required
		} else if edgePrecedence == edge.PrecedenceFailed {
			precedence = edge_ctrl_pb.TerminatorPrecedence_Failed
		}
	}

	assignIds, _ := req.GetBoolHeader(edge.RouterProvidedConnId)
	log.Debugf("client requested router provided connection ids: %v", assignIds)

	log.Debug("establishing listener")

	messageSink := &edgeTerminator{
		MsgChannel:     *edge.NewEdgeMsgChannel(self.ch, connId),
		edgeClientConn: self,
		token:          sessionToken,
		assignIds:      assignIds,
	}

	// need to remove session remove listener on close
	messageSink.onClose = self.listener.factory.stateManager.AddEdgeSessionRemovedListener(sessionToken, func(token string) {
		messageSink.close(true, "session ended")
	})

	self.listener.factory.hostedServices.Put(sessionToken, messageSink)

	terminatorIdentity, _ := req.GetStringHeader(edge.TerminatorIdentityHeader)
	var terminatorIdentitySecret []byte
	if terminatorIdentity != "" {
		terminatorIdentitySecret = req.Headers[edge.TerminatorIdentitySecretHeader]
	}

	request := &edge_ctrl_pb.CreateTerminatorRequest{
		SessionToken:   sessionToken,
		Fingerprints:   self.fingerprints.Prints(),
		PeerData:       hostData,
		Cost:           uint32(cost),
		Precedence:     precedence,
		InstanceId:     terminatorIdentity,
		InstanceSecret: terminatorIdentitySecret,
	}

	if strings.HasPrefix(sessionToken, JwtTokenPrefix) {
		apiSession := manager.GetApiSessionFromCh(ch)

		if apiSession == nil {
			pfxlog.Logger().Errorf("could not find api session for channel, unable to process bind message")
			return
		}

		request.ApiSessionToken = apiSession.Token
	}

	timeout := self.listener.factory.ctrls.DefaultRequestTimeout()
	responseMsg, err := protobufs.MarshalTyped(request).WithTimeout(timeout).SendForReply(ctrlCh)
	if err = xgress_common.CheckForFailureResult(responseMsg, err, edge_ctrl_pb.ContentType_CreateTerminatorResponseType); err != nil {
		log.WithError(err).Warn("error creating terminator")
		messageSink.close(false, "") // don't notify here, as we're notifying next line with a response
		self.sendStateClosedReply(err.Error(), req)
		return
	}

	terminatorId := string(responseMsg.Body)
	messageSink.terminatorId.Store(terminatorId)
	log = log.WithField("terminatorId", terminatorIdentity)

	if messageSink.MsgChannel.IsClosed() {
		log.Warn("edge channel closed while setting up terminator. cleaning up terminator now")
		messageSink.close(false, "edge channel closed")
		return
	}

	log.Debug("registered listener for terminator")
	log.Debug("returning connection state CONNECTED to client")
	self.sendStateConnectedReply(req, nil)

	log.Info("created terminator")
}

func (self *edgeClientConn) processBindV2(manager state.Manager, req *channel.Message, ch channel.Channel, ctrlCh channel.Channel) {
	sessionToken := string(req.Body)

	log := pfxlog.ContextLogger(ch.Label()).
		WithField("sessionToken", sessionToken).
		WithFields(edge.GetLoggerFields(req)).
		WithField("routerId", self.listener.id.Token)

	connId, found := req.GetUint32Header(edge.ConnIdHeader)
	if !found {
		pfxlog.Logger().Errorf("connId not set. unable to process bind message")
		return
	}

	terminatorInstance, _ := req.GetStringHeader(edge.TerminatorIdentityHeader)

	assignIds, _ := req.GetBoolHeader(edge.RouterProvidedConnId)
	log.Debugf("client requested router provided connection ids: %v", assignIds)

	terminatorId := idgen.NewUUIDString()

	terminator := &edgeTerminator{
		MsgChannel:     *edge.NewEdgeMsgChannel(self.ch, connId),
		edgeClientConn: self,
		token:          sessionToken,
		instance:       terminatorInstance,
		assignIds:      assignIds,
		v2:             true,
	}

	log = log.WithField("bindConnId", terminator.MsgChannel.Id())
	terminator.terminatorId.Store(terminatorId)

	log = log.WithField("terminatorId", terminatorId)
	log.Debug("binding service")

	hostData := make(map[uint32][]byte)
	pubKey, hasKey := req.Headers[edge.PublicKeyHeader]
	if hasKey {
		hostData[edge.PublicKeyHeader] = pubKey
	}

	cost := uint16(0)
	if costBytes, hasCost := req.Headers[edge.CostHeader]; hasCost {
		cost = binary.LittleEndian.Uint16(costBytes)
	}

	precedence := edge_ctrl_pb.TerminatorPrecedence_Default
	if precedenceData, hasPrecedence := req.Headers[edge.PrecedenceHeader]; hasPrecedence && len(precedenceData) > 0 {
		edgePrecedence := edge.Precedence(precedenceData[0])
		if edgePrecedence == edge.PrecedenceRequired {
			precedence = edge_ctrl_pb.TerminatorPrecedence_Required
		} else if edgePrecedence == edge.PrecedenceFailed {
			precedence = edge_ctrl_pb.TerminatorPrecedence_Failed
		}
	}

	log.Debug("establishing listener")

	// need to remove session remove listener on close
	terminator.onClose = self.listener.factory.stateManager.AddEdgeSessionRemovedListener(sessionToken, func(token string) {
		terminator.close(true, "session ended")
	})

	self.listener.factory.hostedServices.Put(terminatorId, terminator)

	var terminatorInstanceSecret []byte
	if terminatorInstance != "" {
		terminatorInstanceSecret = req.Headers[edge.TerminatorIdentitySecretHeader]
	}

	request := &edge_ctrl_pb.CreateTerminatorV2Request{
		Address:        terminatorId,
		SessionToken:   sessionToken,
		Fingerprints:   self.fingerprints.Prints(),
		PeerData:       hostData,
		Cost:           uint32(cost),
		Precedence:     precedence,
		InstanceId:     terminatorInstance,
		InstanceSecret: terminatorInstanceSecret,
	}

	if strings.HasPrefix(sessionToken, JwtTokenPrefix) {
		apiSession := manager.GetApiSessionFromCh(ch)

		if apiSession == nil {
			pfxlog.Logger().Errorf("could not find api session for channel, unable to process bind message")
			return
		}

		request.ApiSessionToken = apiSession.Token
	}

	timeout := self.listener.factory.ctrls.DefaultRequestTimeout()
	responseMsg, err := protobufs.MarshalTyped(request).WithTimeout(timeout).SendForReply(ctrlCh)
	resp := &edge_ctrl_pb.CreateTerminatorV2Response{}
	err = protobufs.TypedResponse(resp).Unmarshall(responseMsg, err)
	if err == nil && resp.Result != edge_ctrl_pb.CreateTerminatorResult_Success {
		err = errors.Errorf("terminator create failed: %s", resp.Msg)
	}

	if err != nil {
		log.WithError(err).Warn("error creating terminator")
		terminator.close(false, "") // don't notify here, as we're notifying next line with a response
		self.sendStateClosedReply(err.Error(), req)
		return
	}

	if terminator.MsgChannel.IsClosed() {
		log.Warn("edge channel closed while setting up terminator. cleaning up terminator now")
		terminator.close(false, "edge channel closed")
		return
	}

	log.Debug("registered listener for terminator")
	log.Debug("returning connection state CONNECTED to client")
	self.sendStateConnectedReply(req, nil)

	self.listener.factory.hostedServices.cleanupDuplicates(terminator)

	log.Info("created terminator")
}

func (self *edgeClientConn) processUnbind(manager state.Manager, req *channel.Message, _ channel.Channel) {
	token := string(req.Body)
	atLeastOneTerminatorRemoved := self.listener.factory.hostedServices.unbindSession(token, self)

	if !atLeastOneTerminatorRemoved {
		self.sendStateClosedReply(fmt.Sprintf("no terminator found for token '%s'", token), req)
	}
}

func (self *edgeClientConn) removeTerminator(ctrlCh channel.Channel, token, terminatorId string) error {
	request := &edge_ctrl_pb.RemoveTerminatorRequest{
		SessionToken: token,
		Fingerprints: self.fingerprints.Prints(),
		TerminatorId: terminatorId,
	}

	timeout := self.listener.factory.ctrls.DefaultRequestTimeout()
	responseMsg, err := protobufs.MarshalTyped(request).WithTimeout(timeout).SendForReply(ctrlCh)
	return xgress_common.CheckForFailureResult(responseMsg, err, edge_ctrl_pb.ContentType_RemoveTerminatorResponseType)
}

func (self *edgeClientConn) processUpdateBind(manager state.Manager, req *channel.Message, ch channel.Channel) {
	sessionToken := string(req.Body)

	log := pfxlog.ContextLogger(ch.Label()).WithField("sessionToken", sessionToken).WithFields(edge.GetLoggerFields(req))
	terminators := self.listener.factory.hostedServices.getRelatedTerminators(sessionToken, self)

	if len(terminators) == 0 {
		log.Error("failed to update bind, no listener found")
		return
	}

	ctrlCh := self.listener.factory.ctrls.AnyCtrlChannel()
	if ctrlCh == nil {
		log.Error("no controller available, cannot update terminator")
		return
	}

	for _, terminator := range terminators {
		request := &edge_ctrl_pb.UpdateTerminatorRequest{
			SessionToken: sessionToken,
			Fingerprints: self.fingerprints.Prints(),
			TerminatorId: terminator.terminatorId.Load(),
		}

		if strings.HasPrefix(sessionToken, JwtTokenPrefix) {
			apiSession := manager.GetApiSessionFromCh(ch)
			request.ApiSessionToken = apiSession.Token
		}

		if costVal, hasCost := req.GetUint16Header(edge.CostHeader); hasCost {
			request.UpdateCost = true
			request.Cost = uint32(costVal)
		}

		if precedenceData, hasPrecedence := req.Headers[edge.PrecedenceHeader]; hasPrecedence && len(precedenceData) > 0 {
			edgePrecedence := edge.Precedence(precedenceData[0])
			request.Precedence = edge_ctrl_pb.TerminatorPrecedence_Default
			request.UpdatePrecedence = true
			if edgePrecedence == edge.PrecedenceRequired {
				request.Precedence = edge_ctrl_pb.TerminatorPrecedence_Required
			} else if edgePrecedence == edge.PrecedenceFailed {
				request.Precedence = edge_ctrl_pb.TerminatorPrecedence_Failed
			}
		}

		log = log.WithField("terminator", terminator.terminatorId.Load()).
			WithField("precedence", request.Precedence).
			WithField("cost", request.Cost).
			WithField("updatingPrecedence", request.UpdatePrecedence).
			WithField("updatingCost", request.UpdateCost)

		log.Debug("updating terminator")

		timeout := self.listener.factory.ctrls.DefaultRequestTimeout()
		responseMsg, err := protobufs.MarshalTyped(request).WithTimeout(timeout).SendForReply(ctrlCh)
		if err := xgress_common.CheckForFailureResult(responseMsg, err, edge_ctrl_pb.ContentType_UpdateTerminatorResponseType); err != nil {
			log.WithError(err).Error("terminator update failed")
		} else {
			log.Debug("terminator updated successfully")
		}
	}
}

func (self *edgeClientConn) processHealthEvent(manager state.Manager, req *channel.Message, ch channel.Channel) {
	sessionToken := string(req.Body)
	log := pfxlog.ContextLogger(ch.Label()).WithField("sessionId", sessionToken).WithFields(edge.GetLoggerFields(req))

	ctrlCh := self.listener.factory.ctrls.AnyCtrlChannel()
	if ctrlCh == nil {
		log.Error("no controller available, cannot forward health event")
		return
	}

	terminator, ok := self.listener.factory.hostedServices.Get(sessionToken)

	if !ok {
		log.Error("failed to update bind, no listener found")
		return
	}

	checkPassed, _ := req.GetBoolHeader(edge.HealthStatusHeader)

	request := &edge_ctrl_pb.HealthEventRequest{
		SessionToken: sessionToken,
		Fingerprints: self.fingerprints.Prints(),
		TerminatorId: terminator.terminatorId.Load(),
		CheckPassed:  checkPassed,
	}

	if strings.HasPrefix(sessionToken, JwtTokenPrefix) {
		apiSession := manager.GetApiSessionFromCh(ch)
		request.ApiSessionToken = apiSession.Token
	}

	log = log.WithField("terminator", terminator.terminatorId.Load()).WithField("checkPassed", checkPassed)
	log.Debug("sending health event")

	if err := protobufs.MarshalTyped(request).Send(ctrlCh); err != nil {
		log.WithError(err).Error("send of health event failed")
	}
}

func (self *edgeClientConn) processTraceRoute(msg *channel.Message, ch channel.Channel) {
	log := pfxlog.ContextLogger(ch.Label()).WithFields(edge.GetLoggerFields(msg))

	hops, _ := msg.GetUint32Header(edge.TraceHopCountHeader)
	if hops > 0 {
		hops--
		msg.PutUint32Header(edge.TraceHopCountHeader, hops)
	}

	log.WithField("hops", hops).Debug("traceroute received")
	if hops > 0 {
		self.msgMux.HandleReceive(msg, ch)
	} else {
		ts, _ := msg.GetUint64Header(edge.TimestampHeader)
		connId, _ := msg.GetUint32Header(edge.ConnIdHeader)
		resp := edge.NewTraceRouteResponseMsg(connId, hops, ts, "xgress/edge", "")
		resp.ReplyTo(msg)
		if msgUUID := msg.Headers[edge.UUIDHeader]; msgUUID != nil {
			resp.Headers[edge.UUIDHeader] = msgUUID
		}

		if err := ch.Send(resp); err != nil {
			log.WithError(err).Error("failed to send hop response")
		}
	}
}

func (self *edgeClientConn) sendStateConnectedReply(req *channel.Message, hostData map[uint32][]byte) {
	connId, _ := req.GetUint32Header(edge.ConnIdHeader)
	msg := edge.NewStateConnectedMsg(connId)

	if assignIds, _ := req.GetBoolHeader(edge.RouterProvidedConnId); assignIds {
		msg.PutBoolHeader(edge.RouterProvidedConnId, true)
	}

	for k, v := range hostData {
		msg.Headers[int32(k)] = v
	}
	msg.ReplyTo(req)

	err := msg.WithPriority(channel.High).WithTimeout(5 * time.Second).SendAndWaitForWire(self.ch)
	if err != nil {
		pfxlog.Logger().WithFields(edge.GetLoggerFields(msg)).WithError(err).Error("failed to send state response")
		return
	}
}

func (self *edgeClientConn) sendStateClosedReply(message string, req *channel.Message) {
	connId, _ := req.GetUint32Header(edge.ConnIdHeader)
	msg := edge.NewStateClosedMsg(connId, message)
	msg.ReplyTo(req)

	if errorCode, found := req.GetUint32Header(edge.ErrorCodeHeader); found {
		msg.PutUint32Header(edge.ErrorCodeHeader, errorCode)
	}

	err := msg.WithPriority(channel.High).WithTimeout(5 * time.Second).SendAndWaitForWire(self.ch)
	if err != nil {
		pfxlog.Logger().WithFields(edge.GetLoggerFields(msg)).WithError(err).Error("failed to send state response")
	}
}

func getResultOrFailure(msg *channel.Message, err error, result protobufs.TypedMessage) error {
	if err != nil {
		return err
	}

	if msg.ContentType == int32(edge_ctrl_pb.ContentType_ErrorType) {
		msg := string(msg.Body)
		if msg == "" {
			msg = "error state returned from controller with no message"
		}
		return errors.New(msg)
	}

	if msg.ContentType != result.GetContentType() {
		return errors.Errorf("unexpected response type %v to request. expected %v", msg.ContentType, result.GetContentType())
	}

	return proto.Unmarshal(msg.Body, result)
}
