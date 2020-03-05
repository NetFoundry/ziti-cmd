/*
	Copyright 2020 NetFoundry, Inc.

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

package subcmd

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/netfoundry/ziti-fabric/pb/mgmt_pb"
	"github.com/netfoundry/ziti-foundation/channel2"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	getEndpointClient = NewMgmtClient(getEndpoint)
	getCmd.AddCommand(getEndpoint)
}

var getEndpoint = &cobra.Command{
	Use:   "endpoint <endpointId>",
	Short: "Retrieve a endpoint definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if ch, err := getEndpointClient.Connect(); err == nil {
			request := &mgmt_pb.GetEndpointRequest{
				EndpointId: args[0],
			}
			body, err := proto.Marshal(request)
			if err != nil {
				panic(err)
			}
			requestMsg := channel2.NewMessage(int32(mgmt_pb.ContentType_GetEndpointRequestType), body)
			responseMsg, err := ch.SendAndWaitWithTimeout(requestMsg, 5*time.Second)
			if err != nil {
				panic(err)
			}
			if responseMsg.ContentType == int32(mgmt_pb.ContentType_GetEndpointResponseType) {
				response := &mgmt_pb.GetEndpointResponse{}
				err := proto.Unmarshal(responseMsg.Body, response)
				if err == nil {
					endpoint := response.Endpoint
					fmt.Printf("\n%10s | %-12s|  %-12s| %v\n", "Id", "Service", "Binding", "Destination")
					fmt.Printf("%-10s | %-12s | %-12s | %v\n", endpoint.Id, endpoint.ServiceId, endpoint.Binding,
						fmt.Sprintf("%-12s -> %s", endpoint.RouterId, endpoint.Address))
				} else {
					fmt.Printf("Id not found\n")
				}
			} else if responseMsg.ContentType == channel2.ContentTypeResultType {
				result := channel2.UnmarshalResult(responseMsg)
				if result.Success {
					fmt.Printf("\nsuccess\n\n")
				} else {
					fmt.Printf("\nfailure [%s]\n\n", result.Message)
				}
			} else {
				panic(fmt.Errorf("unexpected response type %v", responseMsg.ContentType))
			}
		} else {
			panic(err)
		}
	},
}
var getEndpointClient *mgmtClient
