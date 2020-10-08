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

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostureCheckDomainCreate posture check domain create
//
// swagger:model PostureCheckDomainCreate
type PostureCheckDomainCreate struct {
	descriptionField *string

	nameField *string

	tagsField Tags

	// domains
	// Required: true
	Domains []string `json:"domains"`
}

// Description gets the description of this subtype
func (m *PostureCheckDomainCreate) Description() *string {
	return m.descriptionField
}

// SetDescription sets the description of this subtype
func (m *PostureCheckDomainCreate) SetDescription(val *string) {
	m.descriptionField = val
}

// Name gets the name of this subtype
func (m *PostureCheckDomainCreate) Name() *string {
	return m.nameField
}

// SetName sets the name of this subtype
func (m *PostureCheckDomainCreate) SetName(val *string) {
	m.nameField = val
}

// Tags gets the tags of this subtype
func (m *PostureCheckDomainCreate) Tags() Tags {
	return m.tagsField
}

// SetTags sets the tags of this subtype
func (m *PostureCheckDomainCreate) SetTags(val Tags) {
	m.tagsField = val
}

// TypeID gets the type Id of this subtype
func (m *PostureCheckDomainCreate) TypeID() PostureCheckType {
	return "PostureCheckDomainCreate"
}

// SetTypeID sets the type Id of this subtype
func (m *PostureCheckDomainCreate) SetTypeID(val PostureCheckType) {
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *PostureCheckDomainCreate) UnmarshalJSON(raw []byte) error {
	var data struct {

		// domains
		// Required: true
		Domains []string `json:"domains"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var base struct {
		/* Just the base type fields. Used for unmashalling polymorphic types.*/

		Description *string `json:"description"`

		Name *string `json:"name"`

		Tags Tags `json:"tags"`

		TypeID PostureCheckType `json:"typeId"`
	}
	buf = bytes.NewBuffer(raw)
	dec = json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&base); err != nil {
		return err
	}

	var result PostureCheckDomainCreate

	result.descriptionField = base.Description

	result.nameField = base.Name

	result.tagsField = base.Tags

	if base.TypeID != result.TypeID() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid typeId value: %q", base.TypeID)
	}

	result.Domains = data.Domains

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m PostureCheckDomainCreate) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// domains
		// Required: true
		Domains []string `json:"domains"`
	}{

		Domains: m.Domains,
	})
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		Description *string `json:"description"`

		Name *string `json:"name"`

		Tags Tags `json:"tags"`

		TypeID PostureCheckType `json:"typeId"`
	}{

		Description: m.Description(),

		Name: m.Name(),

		Tags: m.Tags(),

		TypeID: m.TypeID(),
	})
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this posture check domain create
func (m *PostureCheckDomainCreate) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTags(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDomains(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostureCheckDomainCreate) validateDescription(formats strfmt.Registry) error {

	if err := validate.Required("description", "body", m.Description()); err != nil {
		return err
	}

	return nil
}

func (m *PostureCheckDomainCreate) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name()); err != nil {
		return err
	}

	return nil
}

func (m *PostureCheckDomainCreate) validateTags(formats strfmt.Registry) error {

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

func (m *PostureCheckDomainCreate) validateDomains(formats strfmt.Registry) error {

	if err := validate.Required("domains", "body", m.Domains); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostureCheckDomainCreate) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostureCheckDomainCreate) UnmarshalBinary(b []byte) error {
	var res PostureCheckDomainCreate
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
