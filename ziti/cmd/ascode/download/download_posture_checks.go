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

package download

import (
	"github.com/openziti/edge-api/rest_management_api_client/posture_checks"
)

func (d Download) GetPostureChecks() ([]map[string]interface{}, error) {

	return d.getEntities(
		"PostureChecks",

		func() (int64, error) {
			limit := int64(1)
			resp, err := d.client.PostureChecks.ListPostureChecks(&posture_checks.ListPostureChecksParams{Limit: &limit}, nil)
			if err != nil {
				return -1, err
			}
			return *resp.GetPayload().Meta.Pagination.TotalCount, nil
		},

		func(offset *int64, limit *int64) ([]interface{}, error) {
			resp, _ := d.client.PostureChecks.ListPostureChecks(&posture_checks.ListPostureChecksParams{Limit: limit, Offset: offset}, nil)
			entities := make([]interface{}, len(resp.GetPayload().Data()))
			for i, c := range resp.GetPayload().Data() {
				entities[i] = interface{}(c)
			}
			return entities, nil
		},

		func(entity interface{}) (map[string]interface{}, error) {

			// convert to a map of values
			m := d.ToMap(entity)

			// filter unwanted properties
			d.Filter(m, []string{"id", "_links", "createdAt", "updatedAt",
				"version"})

			return m, nil
		})
}
