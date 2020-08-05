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
	"github.com/openziti/ziti/common/version"
	"github.com/spf13/cobra"
)

var (
	upgradeZitiFabricLong = templates.LongDesc(`
		Upgrades the Ziti Fabric app if there is a newer release
`)

	upgradeZitiFabricExample = templates.Examples(`
		# Upgrades the ziti-fabric app 
		ziti upgrade ziti-fabric
	`)
)

// UpgradeZitiFabricOptions the options for the upgrade ziti-fabric command
type UpgradeZitiFabricOptions struct {
	CreateOptions

	Version string
}

// NewCmdUpgradeZitiFabric defines the command
func NewCmdUpgradeZitiFabric(f cmdutil.Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &UpgradeZitiFabricOptions{
		CreateOptions: CreateOptions{
			CommonOptions: CommonOptions{
				Factory: f,
				Out:     out,
				Err:     errOut,
			},
		},
	}

	cmd := &cobra.Command{
		Use:     "ziti-fabric",
		Short:   "Upgrades the Ziti Fabric app - if there is a new version available",
		Aliases: []string{"fabric", "fab"},
		Long:    upgradeZitiFabricLong,
		Example: upgradeZitiFabricExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			cmdhelper.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&options.Version, "version", "v", "", "The specific version to upgrade to")
	options.addCommonFlags(cmd)
	return cmd
}

// Run implements the command
func (o *UpgradeZitiFabricOptions) Run() error {
	newVersion, err := o.getLatestZitiAppVersion(version.GetBranch(), c.ZITI_FABRIC)
	if err != nil {
		return err
	}

	newVersionStr := newVersion.String()

	if o.Version != "" {
		newVersionStr = o.Version
	}

	o.deleteInstalledBinary(c.ZITI_FABRIC)

	return o.installZitiApp(version.GetBranch(), c.ZITI_FABRIC, true, newVersionStr)
}
