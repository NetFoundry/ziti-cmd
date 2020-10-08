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

package identity

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/openziti/edge/rest_model"
)

// UpdateIdentityReader is a Reader for the UpdateIdentity structure.
type UpdateIdentityReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateIdentityReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateIdentityOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateIdentityBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewUpdateIdentityUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateIdentityNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewUpdateIdentityOK creates a UpdateIdentityOK with default headers values
func NewUpdateIdentityOK() *UpdateIdentityOK {
	return &UpdateIdentityOK{}
}

/*UpdateIdentityOK handles this case with default header values.

The update request was successful and the resource has been altered
*/
type UpdateIdentityOK struct {
	Payload *rest_model.Empty
}

func (o *UpdateIdentityOK) Error() string {
	return fmt.Sprintf("[PUT /identities/{id}][%d] updateIdentityOK  %+v", 200, o.Payload)
}

func (o *UpdateIdentityOK) GetPayload() *rest_model.Empty {
	return o.Payload
}

func (o *UpdateIdentityOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.Empty)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateIdentityBadRequest creates a UpdateIdentityBadRequest with default headers values
func NewUpdateIdentityBadRequest() *UpdateIdentityBadRequest {
	return &UpdateIdentityBadRequest{}
}

/*UpdateIdentityBadRequest handles this case with default header values.

The supplied request contains invalid fields or could not be parsed (json and non-json bodies). The error's code, message, and cause fields can be inspected for further information
*/
type UpdateIdentityBadRequest struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *UpdateIdentityBadRequest) Error() string {
	return fmt.Sprintf("[PUT /identities/{id}][%d] updateIdentityBadRequest  %+v", 400, o.Payload)
}

func (o *UpdateIdentityBadRequest) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *UpdateIdentityBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateIdentityUnauthorized creates a UpdateIdentityUnauthorized with default headers values
func NewUpdateIdentityUnauthorized() *UpdateIdentityUnauthorized {
	return &UpdateIdentityUnauthorized{}
}

/*UpdateIdentityUnauthorized handles this case with default header values.

The currently supplied session does not have the correct access rights to request this resource
*/
type UpdateIdentityUnauthorized struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *UpdateIdentityUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /identities/{id}][%d] updateIdentityUnauthorized  %+v", 401, o.Payload)
}

func (o *UpdateIdentityUnauthorized) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *UpdateIdentityUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateIdentityNotFound creates a UpdateIdentityNotFound with default headers values
func NewUpdateIdentityNotFound() *UpdateIdentityNotFound {
	return &UpdateIdentityNotFound{}
}

/*UpdateIdentityNotFound handles this case with default header values.

The requested resource does not exist
*/
type UpdateIdentityNotFound struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *UpdateIdentityNotFound) Error() string {
	return fmt.Sprintf("[PUT /identities/{id}][%d] updateIdentityNotFound  %+v", 404, o.Payload)
}

func (o *UpdateIdentityNotFound) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *UpdateIdentityNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
