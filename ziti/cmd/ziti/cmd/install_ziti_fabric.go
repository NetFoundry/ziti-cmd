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

package cmd

import (
	"io"

	cmdutil "github.com/openziti/ziti/ziti/cmd/ziti/cmd/factory"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/templates"
	c "github.com/openziti/ziti/ziti/cmd/ziti/constants"
	"github.com/openziti/ziti/ziti/cmd/ziti/internal/log"
	"github.com/openziti/ziti/common/version"
	"github.com/blang/semver"
	"github.com/spf13/cobra"
)

var (
	installZitiFabricLong = templates.LongDesc(`
		Installs the Ziti Fabric app if it has not been installed already
`)

	installZitiFabricExample = templates.Examples(`
		# Install the Ziti Fabric app 
		ziti install ziti-fabric
	`)
)

// InstallZitiFabricOptions the options for the upgrade ziti-fabric command
type InstallZitiFabricOptions struct {
	InstallOptions

	Version string
}

// NewCmdInstallZitiFabric defines the command
func NewCmdInstallZitiFabric(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &InstallZitiFabricOptions{
		InstallOptions: InstallOptions{
			CommonOptions: CommonOptions{
				Factory: f,
				Out:     out,
				Err:     errOut,
			},
		},
	}

	cmd := &cobra.Command{
		Use:     "ziti-fabric",
		Short:   "Installs the Ziti Fabric app - if it has not been installed already",
		Aliases: []string{"fabric"},
		Long:    installZitiFabricLong,
		Example: installZitiFabricExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			cmdhelper.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&options.Version, "version", "v", "", "The specific version to install")
	options.addCommonFlags(cmd)
	return cmd
}

// Run implements the command
func (o *InstallZitiFabricOptions) Run() error {
	newVersion, err := o.getLatestZitiAppVersion(version.GetBranch(), c.ZITI_FABRIC)
	if err != nil {
		return err
	}

	if o.Version != "" {
		newVersion, err = semver.Make(o.Version)
	}

	log.Infoln("Attempting to install '" + c.ZITI_FABRIC + "' version: " + newVersion.String())

	return o.installZitiApp(version.GetBranch(), c.ZITI_FABRIC, false, newVersion.String())
}
