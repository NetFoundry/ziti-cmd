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

package service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/openziti/ziti/controller/rest_model"
)

// UpdateServiceOKCode is the HTTP code returned for type UpdateServiceOK
const UpdateServiceOKCode int = 200

/*UpdateServiceOK The update request was successful and the resource has been altered

swagger:response updateServiceOK
*/
type UpdateServiceOK struct {

	/*
	  In: Body
	*/
	Payload *rest_model.Empty `json:"body,omitempty"`
}

// NewUpdateServiceOK creates UpdateServiceOK with default headers values
func NewUpdateServiceOK() *UpdateServiceOK {

	return &UpdateServiceOK{}
}

// WithPayload adds the payload to the update service o k response
func (o *UpdateServiceOK) WithPayload(payload *rest_model.Empty) *UpdateServiceOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update service o k response
func (o *UpdateServiceOK) SetPayload(payload *rest_model.Empty) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateServiceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateServiceBadRequestCode is the HTTP code returned for type UpdateServiceBadRequest
const UpdateServiceBadRequestCode int = 400

/*UpdateServiceBadRequest The supplied request contains invalid fields or could not be parsed (json and non-json bodies). The error's code, message, and cause fields can be inspected for further information

swagger:response updateServiceBadRequest
*/
type UpdateServiceBadRequest struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewUpdateServiceBadRequest creates UpdateServiceBadRequest with default headers values
func NewUpdateServiceBadRequest() *UpdateServiceBadRequest {

	return &UpdateServiceBadRequest{}
}

// WithPayload adds the payload to the update service bad request response
func (o *UpdateServiceBadRequest) WithPayload(payload *rest_model.APIErrorEnvelope) *UpdateServiceBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update service bad request response
func (o *UpdateServiceBadRequest) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateServiceBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateServiceUnauthorizedCode is the HTTP code returned for type UpdateServiceUnauthorized
const UpdateServiceUnauthorizedCode int = 401

/*UpdateServiceUnauthorized The currently supplied session does not have the correct access rights to request this resource

swagger:response updateServiceUnauthorized
*/
type UpdateServiceUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewUpdateServiceUnauthorized creates UpdateServiceUnauthorized with default headers values
func NewUpdateServiceUnauthorized() *UpdateServiceUnauthorized {

	return &UpdateServiceUnauthorized{}
}

// WithPayload adds the payload to the update service unauthorized response
func (o *UpdateServiceUnauthorized) WithPayload(payload *rest_model.APIErrorEnvelope) *UpdateServiceUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update service unauthorized response
func (o *UpdateServiceUnauthorized) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateServiceUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateServiceNotFoundCode is the HTTP code returned for type UpdateServiceNotFound
const UpdateServiceNotFoundCode int = 404

/*UpdateServiceNotFound The requested resource does not exist

swagger:response updateServiceNotFound
*/
type UpdateServiceNotFound struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewUpdateServiceNotFound creates UpdateServiceNotFound with default headers values
func NewUpdateServiceNotFound() *UpdateServiceNotFound {

	return &UpdateServiceNotFound{}
}

// WithPayload adds the payload to the update service not found response
func (o *UpdateServiceNotFound) WithPayload(payload *rest_model.APIErrorEnvelope) *UpdateServiceNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update service not found response
func (o *UpdateServiceNotFound) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateServiceNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateServiceTooManyRequestsCode is the HTTP code returned for type UpdateServiceTooManyRequests
const UpdateServiceTooManyRequestsCode int = 429

/*UpdateServiceTooManyRequests The resource requested is rate limited and the rate limit has been exceeded

swagger:response updateServiceTooManyRequests
*/
type UpdateServiceTooManyRequests struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewUpdateServiceTooManyRequests creates UpdateServiceTooManyRequests with default headers values
func NewUpdateServiceTooManyRequests() *UpdateServiceTooManyRequests {

	return &UpdateServiceTooManyRequests{}
}

// WithPayload adds the payload to the update service too many requests response
func (o *UpdateServiceTooManyRequests) WithPayload(payload *rest_model.APIErrorEnvelope) *UpdateServiceTooManyRequests {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update service too many requests response
func (o *UpdateServiceTooManyRequests) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateServiceTooManyRequests) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(429)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateServiceServiceUnavailableCode is the HTTP code returned for type UpdateServiceServiceUnavailable
const UpdateServiceServiceUnavailableCode int = 503

/*UpdateServiceServiceUnavailable The request could not be completed due to the server being busy or in a temporarily bad state

swagger:response updateServiceServiceUnavailable
*/
type UpdateServiceServiceUnavailable struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewUpdateServiceServiceUnavailable creates UpdateServiceServiceUnavailable with default headers values
func NewUpdateServiceServiceUnavailable() *UpdateServiceServiceUnavailable {

	return &UpdateServiceServiceUnavailable{}
}

// WithPayload adds the payload to the update service service unavailable response
func (o *UpdateServiceServiceUnavailable) WithPayload(payload *rest_model.APIErrorEnvelope) *UpdateServiceServiceUnavailable {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update service service unavailable response
func (o *UpdateServiceServiceUnavailable) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateServiceServiceUnavailable) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(503)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
