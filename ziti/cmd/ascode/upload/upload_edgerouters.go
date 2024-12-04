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

package upload

import (
	"errors"
	"fmt"
	"github.com/openziti/edge-api/rest_management_api_client/edge_router"
	"github.com/openziti/edge-api/rest_model"
	"github.com/openziti/edge-api/rest_util"
	common "github.com/openziti/ziti/internal/ascode"
	"github.com/openziti/ziti/internal/rest/mgmt"
)

func (u *Upload) ProcessEdgeRouters(input map[string][]interface{}) (map[string]string, error) {

	var result = map[string]string{}
	for _, data := range input["edgeRouters"] {
		create := FromMap(data, rest_model.EdgeRouterCreate{})

		// see if the router already exists
		existing := mgmt.EdgeRouterFromFilter(u.client, mgmt.NameFilter(*create.Name))
		if existing != nil {
			log.WithFields(map[string]interface{}{
				"name":         *create.Name,
				"edgeRouterId": *existing.ID,
			}).
				Info("Found existing EdgeRouter, skipping create")
			_, _ = fmt.Fprintf(u.Err, "\u001B[2KSkipping EdgeRouter %s\r", *create.Name)
			continue
		}

		// do the actual create since it doesn't exist
		_, _ = fmt.Fprintf(u.Err, "\u001B[2KCreating EdgeRouterPolicy %s\r", *create.Name)
		if u.verbose {
			log.WithField("name", *create.Name).Debug("Creating EdgeRouter")
		}
		created, createErr := u.client.EdgeRouter.CreateEdgeRouter(&edge_router.CreateEdgeRouterParams{EdgeRouter: create}, nil)
		if createErr != nil {
			if payloadErr, ok := createErr.(rest_util.ApiErrorPayload); ok {
				log.WithFields(map[string]interface{}{
					"field":  payloadErr.GetPayload().Error.Cause.APIFieldError.Field,
					"reason": payloadErr.GetPayload().Error.Cause.APIFieldError.Reason,
				}).Error("Unable to create EdgeRouter")
			} else {
				log.WithField("err", createErr).Error("Unable to create EdgeRouter")
				return nil, createErr
			}
		}
		if u.verbose {
			log.WithFields(map[string]interface{}{
				"name":         *create.Name,
				"edgeRouterId": created.Payload.Data.ID,
			}).
				Info("Created EdgeRouter")
		}

		result[*create.Name] = created.Payload.Data.ID
	}

	return result, nil
}

func (u *Upload) lookupEdgeRouters(roles []string) ([]string, error) {
	edgeRouterRoles := []string{}
	for _, role := range roles {
		if role[0:1] == "@" {
			value := role[1:]
			edgeRouter, _ := common.GetItemFromCache(u.edgeRouterCache, value, func(name string) (interface{}, error) {
				return mgmt.EdgeRouterFromFilter(u.client, mgmt.NameFilter(name)), nil
			})
			if edgeRouter == nil {
				return nil, errors.New("error reading EdgeRouter: " + value)
			}
			edgeRouterId := edgeRouter.(*rest_model.EdgeRouterDetail).ID
			edgeRouterRoles = append(edgeRouterRoles, "@"+*edgeRouterId)
		} else {
			edgeRouterRoles = append(edgeRouterRoles, role)
		}
	}
	return edgeRouterRoles, nil
}
