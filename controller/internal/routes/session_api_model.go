/*
	Copyright NetFoundry, Inc.

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

package routes

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/edge/controller/env"
	"github.com/openziti/edge/controller/model"
	"github.com/openziti/edge/controller/response"
	"github.com/openziti/edge/rest_model"
	"github.com/openziti/fabric/controller/models"
	"github.com/openziti/foundation/util/stringz"
)

const EntityNameSession = "sessions"

var SessionLinkFactory = NewBasicLinkFactory(EntityNameApiSession)

func MapCreateSessionToModel(apiSessionId string, session *rest_model.SessionCreate) *model.Session {
	ret := &model.Session{
		BaseEntity: models.BaseEntity{
			Tags: session.Tags,
		},
		Token:        uuid.New().String(),
		ApiSessionId: apiSessionId,
		ServiceId:    session.ServiceID,
		Type:         string(session.Type),
		SessionCerts: nil,
	}

	return ret
}

func MapSessionToRestEntity(ae *env.AppEnv, _ *response.RequestContext, e models.Entity) (interface{}, error) {
	session, ok := e.(*model.Session)

	if !ok {
		err := fmt.Errorf("entity is not a Session \"%s\"", e.GetId())
		log := pfxlog.Logger()
		log.Error(err)
		return nil, err
	}

	restModel, err := MapSessionToRestModel(ae, session)

	if err != nil {
		err := fmt.Errorf("could not convert to API entity \"%s\": %s", e.GetId(), err)
		log := pfxlog.Logger()
		log.Error(err)
		return nil, err
	}
	return restModel, nil
}

func MapSessionToRestModel(ae *env.AppEnv, sessionModel *model.Session) (*rest_model.SessionDetail, error) {
	service, err := ae.Handlers.EdgeService.Read(sessionModel.ServiceId)
	if err != nil {
		return nil, err
	}

	edgeRouters, err := getSessionEdgeRouters(ae, sessionModel)
	if err != nil {
		return nil, err
	}

	apiSession, err := ae.Handlers.ApiSession.Read(sessionModel.ApiSessionId)
	if err != nil {
		return nil, err
	}

	ret := &rest_model.SessionDetail{
		BaseEntity:   BaseEntityToRestModel(sessionModel, SessionLinkFactory),
		APISession:   ToEntityRef("", apiSession, ApiSessionLinkFactory),
		APISessionID: &apiSession.Id,
		Service:      ToEntityRef(service.Name, service, ServiceLinkFactory),
		ServiceID:    &service.Id,
		EdgeRouters:  edgeRouters,
		Type:         rest_model.DialBind(sessionModel.Type),
		Token:        &sessionModel.Token,
	}

	return ret, nil
}

func MapSessionsToRestEntities(ae *env.AppEnv, rc *response.RequestContext, sessions []*model.Session) ([]interface{}, error) {
	var ret []interface{}
	for _, session := range sessions {
		restEntity, err := MapSessionToRestEntity(ae, rc, session)

		if err != nil {
			return nil, err
		}

		ret = append(ret, restEntity)
	}

	return ret, nil
}

func getSessionEdgeRouters(ae *env.AppEnv, ns *model.Session) ([]*rest_model.SessionEdgeRouter, error) {
	var edgeRouters []*rest_model.SessionEdgeRouter

	edgeRoutersForSession, err := ae.Handlers.EdgeRouter.ListForSession(ns.Id)
	if err != nil {
		return nil, err
	}

	for _, edgeRouter := range edgeRoutersForSession.EdgeRouters {
		onlineEdgeRouter := ae.Broker.GetOnlineEdgeRouter(edgeRouter.Id)

		if onlineEdgeRouter != nil {
			restModel := &rest_model.SessionEdgeRouter{
				Hostname: stringz.OrEmpty(onlineEdgeRouter.Hostname),
				Name:     edgeRouter.Name,
				Urls:     map[string]string{},
			}

			for p, url := range onlineEdgeRouter.EdgeRouterProtocols {
				restModel.Urls[p] = url
			}

			pfxlog.Logger().Debugf("Returning %+v to %+v, with urls: %+v", edgeRouter, restModel, restModel.Urls)
			edgeRouters = append(edgeRouters, restModel)
		}
	}

	return edgeRouters, nil
}
