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

package model

import (
	"github.com/openziti/edge/controller/persistence"
	"github.com/openziti/fabric/controller/models"
	"github.com/openziti/foundation/storage/boltz"
	"go.etcd.io/bbolt"
	"strings"
)

func NewPostureCheckHandler(env Env) *PostureCheckHandler {
	handler := &PostureCheckHandler{
		baseHandler: newBaseHandler(env, env.GetStores().PostureCheck),
	}
	handler.impl = handler
	return handler
}

type PostureCheckHandler struct {
	baseHandler
}

func (handler *PostureCheckHandler) newModelEntity() boltEntitySink {
	return &PostureCheck{}
}

func (handler *PostureCheckHandler) Create(postureCheckModel *PostureCheck) (string, error) {
	return handler.createEntity(postureCheckModel)
}

func (handler *PostureCheckHandler) Read(id string) (*PostureCheck, error) {
	modelEntity := &PostureCheck{}
	if err := handler.readEntity(id, modelEntity); err != nil {
		return nil, err
	}
	return modelEntity, nil
}

func (handler *PostureCheckHandler) readInTx(tx *bbolt.Tx, id string) (*PostureCheck, error) {
	modelEntity := &PostureCheck{}
	if err := handler.readEntityInTx(tx, id, modelEntity); err != nil {
		return nil, err
	}
	return modelEntity, nil
}

func (handler *PostureCheckHandler) IsUpdated(field string) bool {
	return strings.EqualFold(field, persistence.FieldName) ||
		strings.EqualFold(field, boltz.FieldTags)
}

func (handler *PostureCheckHandler) Update(ca *PostureCheck) error {
	return handler.updateEntity(ca, handler)
}

func (handler *PostureCheckHandler) Patch(ca *PostureCheck, checker boltz.FieldChecker) error {
	combinedChecker := &AndFieldChecker{first: handler, second: checker}
	return handler.patchEntity(ca, combinedChecker)
}

func (handler *PostureCheckHandler) Delete(id string) error {
	return handler.deleteEntity(id)
}

func (handler *PostureCheckHandler) Query(query string) (*PostureCheckListResult, error) {
	result := &PostureCheckListResult{handler: handler}
	if err := handler.list(query, result.collect); err != nil {
		return nil, err
	}
	return result, nil
}

type PostureCheckListResult struct {
	handler       *PostureCheckHandler
	PostureChecks []*PostureCheck
	models.QueryMetaData
}

func (result *PostureCheckListResult) collect(tx *bbolt.Tx, ids []string, queryMetaData *models.QueryMetaData) error {
	result.QueryMetaData = *queryMetaData
	for _, key := range ids {
		entity, err := result.handler.readInTx(tx, key)
		if err != nil {
			return err
		}
		result.PostureChecks = append(result.PostureChecks, entity)
	}
	return nil
}
