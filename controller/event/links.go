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

package event

import (
	"fmt"
	"time"
)

type LinkEventType string

const (
	LinkEventNS = "fabric.links"

	LinkDialed                     LinkEventType = "dialed"
	LinkFault                      LinkEventType = "fault"
	LinkDuplicate                  LinkEventType = "duplicate"
	LinkConnected                  LinkEventType = "connected"
	LinkFromRouterNew              LinkEventType = "routerLinkNew"
	LinkFromRouterKnown            LinkEventType = "routerLinkKnown"
	LinkFromRouterDisconnectedDest LinkEventType = "routerLinkDisconnectedDest"
)

// A LinkConnection describes a physical connection that forms a link. A Link may be made
// up of multiple LinkConnections.
//
// Link ids currently have three valid values:
//   - single - meaning the link has a single connection
//   - payload - a connection used only for payloads
//   - ack - a connection used only for acks
type LinkConnection struct {
	// The connection identifier.
	Id string `json:"id"`

	// The connection address on dialing router side.
	LocalAddr string `json:"local_addr"`

	// The connection address on accepting router side.
	RemoteAddr string `json:"remote_addr"`
}

// A LinkEvent will be emitted for various link lifecycle events.
//
// Valid values for link event types are:
//   - dialed
//   - fault
//   - duplicate
//   - connected
//   - routerLinkNew
//   - routerLinkKnown
//   - routerLinkDisconnectedDest
//
// Examples:
//
//	{
//	 "namespace": "fabric.links",
//	 "event_src_id": "ctrl1",
//	 "timestamp": "2022-07-15T18:10:19.752766075-04:00",
//	 "event_type": "dialed",
//	 "link_id": "47kGIApCXI29VQoCA1xXWI",
//	 "src_router_id": "niY.XmLArx",
//	 "dst_router_id": "YPpTEd8JP",
//	 "protocol": "tls",
//	 "dial_address": "tls:127.0.0.1:4024",
//	 "cost": 1
//	}
//
//	{
//	 "namespace": "fabric.links",
//	 "event_src_id": "ctrl1",
//	 "timestamp": "2022-07-15T18:10:19.973626185-04:00",
//	 "event_type": "connected",
//	 "link_id": "47kGIApCXI29VQoCA1xXWI",
//	 "src_router_id": "niY.XmLArx",
//	 "dst_router_id": "YPpTEd8JP",
//	 "protocol": "tls",
//	 "dial_address": "tls:127.0.0.1:4024",
//	 "cost": 1,
//	 "connections": [
//	   {
//	     "id": "ack",
//	     "local_addr": "tcp:127.0.0.1:49138",
//	     "remote_addr": "tcp:127.0.0.1:4024"
//	   },
//	   {
//	     "id": "payload",
//	     "local_addr": "tcp:127.0.0.1:49136",
//	     "remote_addr": "tcp:127.0.0.1:4024"
//	   }
//	 ]
//	}
//
//	{
//	  "namespace": "fabric.links",
//	  "event_src_id": "ctrl1",
//	  "timestamp": "2022-07-15T18:10:19.973867809-04:00",
//	  "event_type": "fault",
//	  "link_id": "6slUYCqOB85YTfdiD8I5pl",
//	  "src_router_id": "YPpTEd8JP",
//	  "dst_router_id": "niY.XmLArx",
//	  "protocol": "tls",
//	  "dial_address": "tls:127.0.0.1:4023",
//	  "cost": 1
//	}
//
//	{
//	 "namespace": "fabric.links",
//	 "event_src_id": "ctrl1",
//	 "timestamp": "2022-07-15T18:10:19.974177638-04:00",
//	 "event_type": "routerLinkKnown",
//	 "link_id": "47kGIApCXI29VQoCA1xXWI",
//	 "src_router_id": "niY.XmLArx",
//	 "dst_router_id": "YPpTEd8JP",
//	 "protocol": "tls",
//	 "dial_address": "tls:127.0.0.1:4024",
//	 "cost": 1
//	}
type LinkEvent struct {
	Namespace  string    `json:"namespace"`
	EventSrcId string    `json:"event_src_id"`
	Timestamp  time.Time `json:"timestamp"`

	// The link event type. See above for valid values.
	EventType LinkEventType `json:"event_type"`

	// The link identifier.
	LinkId string `json:"link_id"`

	// The id of the dialing router.
	SrcRouterId string `json:"src_router_id"`

	// The id of the accepting router.
	DstRouterId string `json:"dst_router_id"`

	// The link protocol.
	Protocol string `json:"protocol"`

	// The address dialed.
	DialAddress string `json:"dial_address"`

	// The link cost.
	Cost int32 `json:"cost"`

	// The connections making up the link.
	Connections []*LinkConnection `json:"connections,omitempty"`
}

func (event *LinkEvent) String() string {
	return fmt.Sprintf("%v.%v time=%v linkId=%v srcRouterId=%v dstRouterId=%v",
		event.Namespace, event.EventType, event.Timestamp, event.LinkId, event.SrcRouterId, event.DstRouterId)
}

type LinkEventHandler interface {
	AcceptLinkEvent(event *LinkEvent)
}
