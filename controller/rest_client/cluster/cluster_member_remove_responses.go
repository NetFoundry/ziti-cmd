// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright NetFoundry Inc.
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

package cluster

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/openziti/ziti/controller/rest_model"
)

// ClusterMemberRemoveReader is a Reader for the ClusterMemberRemove structure.
type ClusterMemberRemoveReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ClusterMemberRemoveReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewClusterMemberRemoveOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewClusterMemberRemoveBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewClusterMemberRemoveUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewClusterMemberRemoveNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 429:
		result := NewClusterMemberRemoveTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewClusterMemberRemoveOK creates a ClusterMemberRemoveOK with default headers values
func NewClusterMemberRemoveOK() *ClusterMemberRemoveOK {
	return &ClusterMemberRemoveOK{}
}

/* ClusterMemberRemoveOK describes a response with status code 200, with default header values.

Base empty response
*/
type ClusterMemberRemoveOK struct {
	Payload *rest_model.Empty
}

func (o *ClusterMemberRemoveOK) Error() string {
	return fmt.Sprintf("[POST /cluster/remove-member][%d] clusterMemberRemoveOK  %+v", 200, o.Payload)
}
func (o *ClusterMemberRemoveOK) GetPayload() *rest_model.Empty {
	return o.Payload
}

func (o *ClusterMemberRemoveOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.Empty)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewClusterMemberRemoveBadRequest creates a ClusterMemberRemoveBadRequest with default headers values
func NewClusterMemberRemoveBadRequest() *ClusterMemberRemoveBadRequest {
	return &ClusterMemberRemoveBadRequest{}
}

/* ClusterMemberRemoveBadRequest describes a response with status code 400, with default header values.

The supplied request contains invalid fields or could not be parsed (json and non-json bodies). The error's code, message, and cause fields can be inspected for further information
*/
type ClusterMemberRemoveBadRequest struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *ClusterMemberRemoveBadRequest) Error() string {
	return fmt.Sprintf("[POST /cluster/remove-member][%d] clusterMemberRemoveBadRequest  %+v", 400, o.Payload)
}
func (o *ClusterMemberRemoveBadRequest) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *ClusterMemberRemoveBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewClusterMemberRemoveUnauthorized creates a ClusterMemberRemoveUnauthorized with default headers values
func NewClusterMemberRemoveUnauthorized() *ClusterMemberRemoveUnauthorized {
	return &ClusterMemberRemoveUnauthorized{}
}

/* ClusterMemberRemoveUnauthorized describes a response with status code 401, with default header values.

The currently supplied session does not have the correct access rights to request this resource
*/
type ClusterMemberRemoveUnauthorized struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *ClusterMemberRemoveUnauthorized) Error() string {
	return fmt.Sprintf("[POST /cluster/remove-member][%d] clusterMemberRemoveUnauthorized  %+v", 401, o.Payload)
}
func (o *ClusterMemberRemoveUnauthorized) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *ClusterMemberRemoveUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewClusterMemberRemoveNotFound creates a ClusterMemberRemoveNotFound with default headers values
func NewClusterMemberRemoveNotFound() *ClusterMemberRemoveNotFound {
	return &ClusterMemberRemoveNotFound{}
}

/* ClusterMemberRemoveNotFound describes a response with status code 404, with default header values.

The requested resource does not exist
*/
type ClusterMemberRemoveNotFound struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *ClusterMemberRemoveNotFound) Error() string {
	return fmt.Sprintf("[POST /cluster/remove-member][%d] clusterMemberRemoveNotFound  %+v", 404, o.Payload)
}
func (o *ClusterMemberRemoveNotFound) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *ClusterMemberRemoveNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewClusterMemberRemoveTooManyRequests creates a ClusterMemberRemoveTooManyRequests with default headers values
func NewClusterMemberRemoveTooManyRequests() *ClusterMemberRemoveTooManyRequests {
	return &ClusterMemberRemoveTooManyRequests{}
}

/* ClusterMemberRemoveTooManyRequests describes a response with status code 429, with default header values.

The resource requested is rate limited and the rate limit has been exceeded
*/
type ClusterMemberRemoveTooManyRequests struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *ClusterMemberRemoveTooManyRequests) Error() string {
	return fmt.Sprintf("[POST /cluster/remove-member][%d] clusterMemberRemoveTooManyRequests  %+v", 429, o.Payload)
}
func (o *ClusterMemberRemoveTooManyRequests) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *ClusterMemberRemoveTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
