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

package enroll

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/ziti/router"
	"github.com/openziti/ziti/router/enroll"
	"github.com/openziti/ziti/ziti/cmd/common"
	"github.com/spf13/cobra"
	"os"
)

var jwtPath *string
var engine *string
var keyAlg ziti.KeyAlgVar

func NewEnrollGwCmd() *cobra.Command {
	var enrollEdgeRouterCmd = &cobra.Command{
		Use:   "enroll <config>",
		Short: "Enroll a router as an edge router",
		Args:  cobra.ExactArgs(1),
		Run:   enrollGw,
	}

	jwtPath = enrollEdgeRouterCmd.Flags().StringP("jwt", "j", "", "The path to a JWT file")
	engine = enrollEdgeRouterCmd.Flags().StringP("engine", "e", "", "An engine")
	if err := keyAlg.Set("RSA"); err != nil { // set default
		panic(err)
	}
	enrollEdgeRouterCmd.Flags().VarP(&keyAlg, "keyAlg", "a", "Crypto algorithm to use when generating private key")

	return enrollEdgeRouterCmd
}

func enrollGw(cmd *cobra.Command, args []string) {
	log := pfxlog.Logger()
	if cfg, err := router.LoadConfigWithOptions(args[0], false); err == nil {
		cfg.SetFlags(common.GetFlags(cmd))

		enroller := enroll.NewRestEnroller(cfg)

		jwtBuf, err := os.ReadFile(*jwtPath)
		if err != nil {
			log.Panicf("could not load JWT file from path [%s]", *jwtPath)
		}

		if err := enroller.Enroll(jwtBuf, true, *engine, keyAlg); err != nil {
			log.Fatalf("enrollment failure: (%v)", err)
		}
	} else {
		panic(err)
	}
}
