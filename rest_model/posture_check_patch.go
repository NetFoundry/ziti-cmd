// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright NetFoundry, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// __          __              _
// \ \        / /             (_)
//  \ \  /\  / /_ _ _ __ _ __  _ _ __   __ _
//   \ \/  \/ / _` | '__| '_ \| | '_ \ / _` |
//    \  /\  / (_| | |  | | | | | | | | (_| | : This file is generated, do not edit it.
//     \/  \/ \__,_|_|  |_| |_|_|_| |_|\__, |
//                                      __/ |
//                                     |___/

package rest_model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostureCheckPatch posture check patch
//
// swagger:discriminator PostureCheckPatch typeId
type PostureCheckPatch interface {
	runtime.Validatable

	// description
	Description() string
	SetDescription(string)

	// name
	Name() string
	SetName(string)

	// tags
	Tags() Tags
	SetTags(Tags)

	// AdditionalProperties in base type shoud be handled just like regular properties
	// At this moment, the base type property is pushed down to the subtype
}

type postureCheckPatch struct {
	descriptionField string

	nameField string

	tagsField Tags
}

// Description gets the description of this polymorphic type
func (m *postureCheckPatch) Description() string {
	return m.descriptionField
}

// SetDescription sets the description of this polymorphic type
func (m *postureCheckPatch) SetDescription(val string) {
	m.descriptionField = val
}

// Name gets the name of this polymorphic type
func (m *postureCheckPatch) Name() string {
	return m.nameField
}

// SetName sets the name of this polymorphic type
func (m *postureCheckPatch) SetName(val string) {
	m.nameField = val
}

// Tags gets the tags of this polymorphic type
func (m *postureCheckPatch) Tags() Tags {
	return m.tagsField
}

// SetTags sets the tags of this polymorphic type
func (m *postureCheckPatch) SetTags(val Tags) {
	m.tagsField = val
}

// UnmarshalPostureCheckPatchSlice unmarshals polymorphic slices of PostureCheckPatch
func UnmarshalPostureCheckPatchSlice(reader io.Reader, consumer runtime.Consumer) ([]PostureCheckPatch, error) {
	var elements []json.RawMessage
	if err := consumer.Consume(reader, &elements); err != nil {
		return nil, err
	}

	var result []PostureCheckPatch
	for _, element := range elements {
		obj, err := unmarshalPostureCheckPatch(element, consumer)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
	}
	return result, nil
}

// UnmarshalPostureCheckPatch unmarshals polymorphic PostureCheckPatch
func UnmarshalPostureCheckPatch(reader io.Reader, consumer runtime.Consumer) (PostureCheckPatch, error) {
	// we need to read this twice, so first into a buffer
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return unmarshalPostureCheckPatch(data, consumer)
}

func unmarshalPostureCheckPatch(data []byte, consumer runtime.Consumer) (PostureCheckPatch, error) {
	buf := bytes.NewBuffer(data)
	buf2 := bytes.NewBuffer(data)

	// the first time this is read is to fetch the value of the typeId property.
	var getType struct {
		TypeID string `json:"typeId"`
	}
	if err := consumer.Consume(buf, &getType); err != nil {
		return nil, err
	}

	if err := validate.RequiredString("typeId", "body", getType.TypeID); err != nil {
		return nil, err
	}

	// The value of typeId is used to determine which type to create and unmarshal the data into
	switch getType.TypeID {
	case "PostureCheckDomainPatch":
		var result PostureCheckDomainPatch
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case "PostureCheckMACAddressPatch":
		var result PostureCheckMACAddressPatch
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case "PostureCheckOperatingSystemPatch":
		var result PostureCheckOperatingSystemPatch
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case "PostureCheckPatch":
		var result postureCheckPatch
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil
	case "PostureCheckProcessPatch":
		var result PostureCheckProcessPatch
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, errors.New(422, "invalid typeId value: %q", getType.TypeID)
}

// Validate validates this posture check patch
func (m *postureCheckPatch) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTags(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *postureCheckPatch) validateTags(formats strfmt.Registry) error {

	if swag.IsZero(m.Tags()) { // not required
		return nil
	}

	if err := m.Tags().Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("tags")
		}
		return err
	}

	return nil
}
