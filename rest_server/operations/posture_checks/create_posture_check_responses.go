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

package posture_checks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/openziti/edge/rest_model"
)

// CreatePostureCheckOKCode is the HTTP code returned for type CreatePostureCheckOK
const CreatePostureCheckOKCode int = 200

/*CreatePostureCheckOK The create request was successful and the resource has been added at the following location

swagger:response createPostureCheckOK
*/
type CreatePostureCheckOK struct {

	/*
	  In: Body
	*/
	Payload *rest_model.CreateEnvelope `json:"body,omitempty"`
}

// NewCreatePostureCheckOK creates CreatePostureCheckOK with default headers values
func NewCreatePostureCheckOK() *CreatePostureCheckOK {

	return &CreatePostureCheckOK{}
}

// WithPayload adds the payload to the create posture check o k response
func (o *CreatePostureCheckOK) WithPayload(payload *rest_model.CreateEnvelope) *CreatePostureCheckOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create posture check o k response
func (o *CreatePostureCheckOK) SetPayload(payload *rest_model.CreateEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreatePostureCheckOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreatePostureCheckBadRequestCode is the HTTP code returned for type CreatePostureCheckBadRequest
const CreatePostureCheckBadRequestCode int = 400

/*CreatePostureCheckBadRequest The supplied request contains invalid fields or could not be parsed (json and non-json bodies). The error's code, message, and cause fields can be inspected for further information

swagger:response createPostureCheckBadRequest
*/
type CreatePostureCheckBadRequest struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewCreatePostureCheckBadRequest creates CreatePostureCheckBadRequest with default headers values
func NewCreatePostureCheckBadRequest() *CreatePostureCheckBadRequest {

	return &CreatePostureCheckBadRequest{}
}

// WithPayload adds the payload to the create posture check bad request response
func (o *CreatePostureCheckBadRequest) WithPayload(payload *rest_model.APIErrorEnvelope) *CreatePostureCheckBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create posture check bad request response
func (o *CreatePostureCheckBadRequest) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreatePostureCheckBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreatePostureCheckUnauthorizedCode is the HTTP code returned for type CreatePostureCheckUnauthorized
const CreatePostureCheckUnauthorizedCode int = 401

/*CreatePostureCheckUnauthorized The currently supplied session does not have the correct access rights to request this resource

swagger:response createPostureCheckUnauthorized
*/
type CreatePostureCheckUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewCreatePostureCheckUnauthorized creates CreatePostureCheckUnauthorized with default headers values
func NewCreatePostureCheckUnauthorized() *CreatePostureCheckUnauthorized {

	return &CreatePostureCheckUnauthorized{}
}

// WithPayload adds the payload to the create posture check unauthorized response
func (o *CreatePostureCheckUnauthorized) WithPayload(payload *rest_model.APIErrorEnvelope) *CreatePostureCheckUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create posture check unauthorized response
func (o *CreatePostureCheckUnauthorized) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreatePostureCheckUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
