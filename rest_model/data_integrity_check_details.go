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
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DataIntegrityCheckDetails data integrity check details
//
// swagger:model dataIntegrityCheckDetails
type DataIntegrityCheckDetails struct {

	// end time
	// Required: true
	// Format: date-time
	EndTime *strfmt.DateTime `json:"endTime"`

	// error
	// Required: true
	Error *string `json:"error"`

	// fixing errors
	// Required: true
	FixingErrors *bool `json:"fixingErrors"`

	// in progress
	// Required: true
	InProgress *bool `json:"inProgress"`

	// results
	// Required: true
	Results DataIntegrityCheckDetailList `json:"results"`

	// start time
	// Required: true
	// Format: date-time
	StartTime *strfmt.DateTime `json:"startTime"`

	// too many errors
	// Required: true
	TooManyErrors *bool `json:"tooManyErrors"`
}

// Validate validates this data integrity check details
func (m *DataIntegrityCheckDetails) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEndTime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateError(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFixingErrors(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInProgress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResults(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartTime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTooManyErrors(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DataIntegrityCheckDetails) validateEndTime(formats strfmt.Registry) error {

	if err := validate.Required("endTime", "body", m.EndTime); err != nil {
		return err
	}

	if err := validate.FormatOf("endTime", "body", "date-time", m.EndTime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *DataIntegrityCheckDetails) validateError(formats strfmt.Registry) error {

	if err := validate.Required("error", "body", m.Error); err != nil {
		return err
	}

	return nil
}

func (m *DataIntegrityCheckDetails) validateFixingErrors(formats strfmt.Registry) error {

	if err := validate.Required("fixingErrors", "body", m.FixingErrors); err != nil {
		return err
	}

	return nil
}

func (m *DataIntegrityCheckDetails) validateInProgress(formats strfmt.Registry) error {

	if err := validate.Required("inProgress", "body", m.InProgress); err != nil {
		return err
	}

	return nil
}

func (m *DataIntegrityCheckDetails) validateResults(formats strfmt.Registry) error {

	if err := validate.Required("results", "body", m.Results); err != nil {
		return err
	}

	if err := m.Results.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("results")
		}
		return err
	}

	return nil
}

func (m *DataIntegrityCheckDetails) validateStartTime(formats strfmt.Registry) error {

	if err := validate.Required("startTime", "body", m.StartTime); err != nil {
		return err
	}

	if err := validate.FormatOf("startTime", "body", "date-time", m.StartTime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *DataIntegrityCheckDetails) validateTooManyErrors(formats strfmt.Registry) error {

	if err := validate.Required("tooManyErrors", "body", m.TooManyErrors); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this data integrity check details based on the context it is used
func (m *DataIntegrityCheckDetails) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DataIntegrityCheckDetails) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Results.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("results")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DataIntegrityCheckDetails) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataIntegrityCheckDetails) UnmarshalBinary(b []byte) error {
	var res DataIntegrityCheckDetails
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
