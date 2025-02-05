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

package exporter

import (
	"encoding/json"
	"errors"
	"github.com/judedaryl/go-arrayutils"
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/edge-api/rest_management_api_client"
	"github.com/openziti/ziti/internal"
	ziticobra "github.com/openziti/ziti/internal/cobra"
	"github.com/openziti/ziti/ziti/cmd/api"
	"github.com/openziti/ziti/ziti/cmd/common"
	"github.com/openziti/ziti/ziti/cmd/edge"
	"github.com/openziti/ziti/ziti/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"slices"
	"strings"
)

var log = pfxlog.Logger()

type Exporter struct {
	Out              io.Writer
	Err              io.Writer
	configCache      map[string]any
	configTypeCache  map[string]any
	authPolicyCache  map[string]any
	externalJwtCache map[string]any
	verbose          bool
}

var output *Output

func NewExportCmd(out io.Writer, errOut io.Writer) *cobra.Command {

	exporter := &Exporter{
		Out: out,
		Err: errOut,
	}

	var outputFormat string
	var outputFile string
	var loginOpts = edge.LoginOptions{
		Options: api.Options{
			CommonOptions: common.CommonOptions{
				Out: os.Stdout,
				Err: os.Stderr,
			},
		},
	}

	cmd := &cobra.Command{
		Use:   "export [entity]",
		Short: "Export entities",
		Long: "Export all or comma separated list of selected entities.\n" +
			"Valid entities are: [all|ca/certificate-authority|identity|edge-router|service|config|config-type|service-policy|edge-router-policy|service-edge-router-policy|external-jwt-signer|auth-policy|posture-check] (default all)",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			exporter.verbose = loginOpts.Verbose

			var parsedOutputFormat OutputFormat
			if strings.ToUpper(outputFormat) == "JSON" {
				parsedOutputFormat = JSON
			} else if strings.ToUpper(outputFormat) == "YAML" {
				parsedOutputFormat = YAML
			} else {
				log.Fatalf("Invalid output format: %s", outputFormat)
			}

			client, err := loginOpts.NewMgmtClient()
			if err != nil {
				log.Fatal(err)
			}

			result, err := exporter.Execute(client, args)
			if err != nil {
				log.Fatal(err)
			}

			if outputFile != "" {
				o, err := NewOutputToFile(loginOpts.Verbose, parsedOutputFormat, outputFile, exporter.Err)
				if err != nil {
					log.Fatal(err)
				}
				output = o
			} else {
				o, err := NewOutputToWriter(loginOpts.Verbose, parsedOutputFormat, exporter.Out, exporter.Err)
				if err != nil {
					log.Fatal(err)
				}
				output = o
			}
			err = output.Write(result)
			if err != nil {
				log.Fatal(err)
			}
		},
		Hidden: true,
	}

	v := viper.New()

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, d.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	viper.SetEnvPrefix(constants.ZITI) // All env vars we seek will be prefixed with "ZITI_"

	// Environment variables can't have dashes in them, so bind them to their equivalent
	// keys with underscores, d.g. --favorite-color to STING_FAVORITE_COLOR
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	v.AutomaticEnv()

	edge.AddLoginFlags(cmd, &loginOpts)
	cmd.Flags().SetInterspersed(true)
	cmd.Flags().StringVar(&outputFormat, "output-format", "JSON", "Output data as either JSON or YAML (default JSON)")
	cmd.Flags().StringVar(&loginOpts.ControllerUrl, "controller-url", "", "The url of the controller")
	ziticobra.SetHelpTemplate(cmd)

	cmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write output to local file")

	return cmd
}

func (exporter *Exporter) Execute(client *rest_management_api_client.ZitiEdgeManagement, input []string) (map[string]interface{}, error) {

	logLvl := logrus.InfoLevel
	if exporter.verbose {
		logLvl = logrus.DebugLevel
	}

	pfxlog.GlobalInit(logLvl, pfxlog.DefaultOptions().Color())
	internal.ConfigureLogFormat(logLvl)

	args := arrayutils.Map(input, strings.ToLower)

	exporter.authPolicyCache = map[string]any{}
	exporter.configCache = map[string]any{}
	exporter.configTypeCache = map[string]any{}
	exporter.externalJwtCache = map[string]any{}

	result := map[string]interface{}{}

	if exporter.IsCertificateAuthorityExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Certificate Authorities\r")
		cas, err := exporter.GetCertificateAuthorities(client)
		if err != nil {
			return nil, err
		}
		result["certificateAuthorities"] = cas
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Certificate Authorities\r\n", len(cas))
	}
	if exporter.IsIdentityExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Identities")
		identities, err := exporter.GetIdentities(client)
		if err != nil {
			return nil, err
		}
		result["identities"] = identities
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Identities\r\n", len(identities))
	}

	if exporter.IsEdgeRouterExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Edge Routers\r")
		routers, err := exporter.GetEdgeRouters(client)
		if err != nil {
			return nil, err
		}
		result["edgeRouters"] = routers
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Edge Routers\r\n", len(routers))
	}
	if exporter.IsServiceExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Services\r")
		services, err := exporter.GetServices(client)
		if err != nil {
			return nil, err
		}
		result["services"] = services
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Services\r\n", len(services))
	}
	if exporter.IsConfigExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Configs\r")
		configs, err := exporter.GetConfigs(client)
		if err != nil {
			return nil, err
		}
		result["configs"] = configs
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Configs\r\n", len(configs))
	}
	if exporter.IsConfigTypeExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Config Types\r")
		configTypes, err := exporter.GetConfigTypes(client)
		if err != nil {
			return nil, err
		}
		result["configTypes"] = configTypes
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d ConfigTypes\r\n", len(configTypes))
	}
	if exporter.IsServicePolicyExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Service Policies")
		servicePolicies, err := exporter.GetServicePolicies(client)
		if err != nil {
			return nil, err
		}
		result["servicePolicies"] = servicePolicies
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Service Policies\r\n", len(servicePolicies))
	}
	if exporter.IsEdgeRouterExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Edge Router Policies\r")
		routerPolicies, err := exporter.GetEdgeRouterPolicies(client)
		if err != nil {
			return nil, err
		}
		result["edgeRouterPolicies"] = routerPolicies
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Edge Router Policies\r\n", len(routerPolicies))
	}
	if exporter.IsServiceEdgeRouterPolicyExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Service Edge Router Policies")
		serviceRouterPolicies, err := exporter.GetServiceEdgeRouterPolicies(client)
		if err != nil {
			return nil, err
		}
		result["serviceEdgeRouterPolicies"] = serviceRouterPolicies
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Service Edge Router Policies\r\n", len(serviceRouterPolicies))
	}
	if exporter.IsExtJwtSignerExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting External JWT Signers")
		externalJwtSigners, err := exporter.GetExternalJwtSigners(client)
		if err != nil {
			return nil, err
		}
		result["externalJwtSigners"] = externalJwtSigners
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d External JWT Signers\r\n", len(externalJwtSigners))
	}
	if exporter.IsAuthPolicyExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Auth Policies")
		authPolicies, err := exporter.GetAuthPolicies(client)
		if err != nil {
			return nil, err
		}
		result["authPolicies"] = authPolicies
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Auth Policies\r\n", len(authPolicies))
	}
	if exporter.IsPostureCheckExportRequired(args) {
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exporting Posture Checks")
		postureChecks, err := exporter.GetPostureChecks(client)
		if err != nil {
			return nil, err
		}
		result["postureChecks"] = postureChecks
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Exported %d Posture Checks\r\n", len(postureChecks))
	}

	_, _ = internal.FPrintfReusingLine(exporter.Err, "Export complete\rn")

	return result, nil
}

type ClientCount func() (int64, error)
type ClientList func(offset *int64, limit *int64) ([]interface{}, error)
type EntityProcessor func(item interface{}) (map[string]interface{}, error)

func (exporter *Exporter) getEntities(entityName string, count ClientCount, list ClientList, processor EntityProcessor) ([]map[string]interface{}, error) {

	totalCount, countErr := count()
	if countErr != nil {
		return nil, errors.Join(errors.New("error reading total number of "+entityName), countErr)
	}

	result := []map[string]interface{}{}

	offset := int64(0)
	limit := int64(500)
	more := true
	for more {
		resp, err := list(&offset, &limit)
		_, _ = internal.FPrintfReusingLine(exporter.Err, "Reading %d/%d %s", offset, totalCount, entityName)
		if err != nil {
			return nil, errors.Join(errors.New("error reading "+entityName), err)
		}

		for _, item := range resp {
			m, err := processor(item)
			if err != nil {
				return nil, err
			}
			if m != nil {
				result = append(result, m)
			}
		}

		more = offset < totalCount
		offset += limit
	}

	_, _ = internal.FPrintflnReusingLine(exporter.Err, "Read %d %s", len(result), entityName)

	return result, nil

}

func (exporter *Exporter) ToMap(input interface{}) (map[string]interface{}, error) {
	jsonData, _ := json.MarshalIndent(input, "", "")
	m := map[string]interface{}{}
	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (exporter *Exporter) defaultRoleAttributes(m map[string]interface{}) {
	if m["roleAttributes"] == nil {
		m["roleAttributes"] = []string{}
	}
}

func (exporter *Exporter) Filter(m map[string]interface{}, properties []string) {

	// remove any properties that are not requested
	for k := range m {
		if slices.Contains(properties, k) {
			delete(m, k)
		}
	}
}

type OutputFormat string

const (
	JSON OutputFormat = "JSON"
	YAML OutputFormat = "YAML"
)
