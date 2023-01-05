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

package agentcli

import (
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/openziti/agent"
	"github.com/openziti/channel/v2"
	"github.com/openziti/fabric/pb/mgmt_pb"
	"github.com/openziti/identity"
	"github.com/openziti/ziti/ziti/cmd/api"
	"github.com/openziti/ziti/ziti/cmd/common"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/helpers"
	"github.com/spf13/cobra"
)

type AgentCtrlRaftListOptions struct {
	AgentOptions
}

func NewAgentCtrlRaftList(p common.OptionsProvider) *cobra.Command {
	options := &AgentCtrlRaftListOptions{
		AgentOptions: AgentOptions{
			CommonOptions: p(),
		},
	}

	cmd := &cobra.Command{
		Args: cobra.RangeArgs(0, 1),
		Use:  "raft-list <optional-target>",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			cmdhelper.CheckErr(err)
		},
	}

	return cmd
}

// Run implements the command
func (o *AgentCtrlRaftListOptions) Run() error {
	var addr string
	var err error

	if len(o.Args) == 1 {
		addr, err = agent.ParseGopsAddress(o.Args)
		if err != nil {
			return err
		}
	}

	return agent.MakeRequestF(addr, agent.CustomOpAsync, []byte{byte(AgentAppController)}, o.makeRequest)
}

func (o *AgentCtrlRaftListOptions) makeRequest(conn net.Conn) error {
	options := channel.DefaultOptions()
	options.ConnectTimeout = time.Second
	dialer := channel.NewExistingConnDialer(&identity.TokenId{Token: "agent"}, conn, nil)
	ch, err := channel.NewChannel("agent", dialer, nil, options)
	if err != nil {
		return err
	}

	msg := channel.NewMessage(int32(mgmt_pb.ContentType_RaftListMembersRequestType), nil)
	reply, err := msg.WithTimeout(5 * time.Second).SendForReply(ch)
	if err != nil {
		return err
	}
	if reply.ContentType == channel.ContentTypeResultType {
		result := channel.UnmarshalResult(reply)
		if result.Success {
			fmt.Println("success")
		} else {
			fmt.Printf("error: %v\n", result.Message)
		}
	} else if reply.ContentType == int32(mgmt_pb.ContentType_RaftListMembersResponseType) {
		resp := &mgmt_pb.RaftMemberListResponse{}
		if err = proto.Unmarshal(reply.Body, resp); err != nil {
			return err
		}
		t := table.NewWriter()
		t.SetStyle(table.StyleRounded)
		t.AppendHeader(table.Row{"Id", "Address", "Voter", "Leader", "Version", "Connected"})
		for _, m := range resp.Members {
			t.AppendRow(table.Row{m.Id, m.Addr, m.IsVoter, m.IsLeader, m.Version, m.IsConnected})
		}
		api.RenderTable(&api.Options{
			CommonOptions: o.CommonOptions,
		}, t, nil)
	}
	return nil
}
