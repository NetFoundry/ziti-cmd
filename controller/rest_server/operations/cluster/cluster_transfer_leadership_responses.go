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
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/openziti/ziti/controller/rest_model"
)

// ClusterTransferLeadershipOKCode is the HTTP code returned for type ClusterTransferLeadershipOK
const ClusterTransferLeadershipOKCode int = 200

/*ClusterTransferLeadershipOK Base empty response

swagger:response clusterTransferLeadershipOK
*/
type ClusterTransferLeadershipOK struct {

	/*
	  In: Body
	*/
	Payload *rest_model.Empty `json:"body,omitempty"`
}

// NewClusterTransferLeadershipOK creates ClusterTransferLeadershipOK with default headers values
func NewClusterTransferLeadershipOK() *ClusterTransferLeadershipOK {

	return &ClusterTransferLeadershipOK{}
}

// WithPayload adds the payload to the cluster transfer leadership o k response
func (o *ClusterTransferLeadershipOK) WithPayload(payload *rest_model.Empty) *ClusterTransferLeadershipOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the cluster transfer leadership o k response
func (o *ClusterTransferLeadershipOK) SetPayload(payload *rest_model.Empty) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ClusterTransferLeadershipOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ClusterTransferLeadershipUnauthorizedCode is the HTTP code returned for type ClusterTransferLeadershipUnauthorized
const ClusterTransferLeadershipUnauthorizedCode int = 401

/*ClusterTransferLeadershipUnauthorized The currently supplied session does not have the correct access rights to request this resource

swagger:response clusterTransferLeadershipUnauthorized
*/
type ClusterTransferLeadershipUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewClusterTransferLeadershipUnauthorized creates ClusterTransferLeadershipUnauthorized with default headers values
func NewClusterTransferLeadershipUnauthorized() *ClusterTransferLeadershipUnauthorized {

	return &ClusterTransferLeadershipUnauthorized{}
}

// WithPayload adds the payload to the cluster transfer leadership unauthorized response
func (o *ClusterTransferLeadershipUnauthorized) WithPayload(payload *rest_model.APIErrorEnvelope) *ClusterTransferLeadershipUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the cluster transfer leadership unauthorized response
func (o *ClusterTransferLeadershipUnauthorized) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ClusterTransferLeadershipUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ClusterTransferLeadershipNotFoundCode is the HTTP code returned for type ClusterTransferLeadershipNotFound
const ClusterTransferLeadershipNotFoundCode int = 404

/*ClusterTransferLeadershipNotFound The requested resource does not exist

swagger:response clusterTransferLeadershipNotFound
*/
type ClusterTransferLeadershipNotFound struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewClusterTransferLeadershipNotFound creates ClusterTransferLeadershipNotFound with default headers values
func NewClusterTransferLeadershipNotFound() *ClusterTransferLeadershipNotFound {

	return &ClusterTransferLeadershipNotFound{}
}

// WithPayload adds the payload to the cluster transfer leadership not found response
func (o *ClusterTransferLeadershipNotFound) WithPayload(payload *rest_model.APIErrorEnvelope) *ClusterTransferLeadershipNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the cluster transfer leadership not found response
func (o *ClusterTransferLeadershipNotFound) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ClusterTransferLeadershipNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ClusterTransferLeadershipTooManyRequestsCode is the HTTP code returned for type ClusterTransferLeadershipTooManyRequests
const ClusterTransferLeadershipTooManyRequestsCode int = 429

/*ClusterTransferLeadershipTooManyRequests The resource requested is rate limited and the rate limit has been exceeded

swagger:response clusterTransferLeadershipTooManyRequests
*/
type ClusterTransferLeadershipTooManyRequests struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewClusterTransferLeadershipTooManyRequests creates ClusterTransferLeadershipTooManyRequests with default headers values
func NewClusterTransferLeadershipTooManyRequests() *ClusterTransferLeadershipTooManyRequests {

	return &ClusterTransferLeadershipTooManyRequests{}
}

// WithPayload adds the payload to the cluster transfer leadership too many requests response
func (o *ClusterTransferLeadershipTooManyRequests) WithPayload(payload *rest_model.APIErrorEnvelope) *ClusterTransferLeadershipTooManyRequests {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the cluster transfer leadership too many requests response
func (o *ClusterTransferLeadershipTooManyRequests) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ClusterTransferLeadershipTooManyRequests) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(429)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ClusterTransferLeadershipInternalServerErrorCode is the HTTP code returned for type ClusterTransferLeadershipInternalServerError
const ClusterTransferLeadershipInternalServerErrorCode int = 500

/*ClusterTransferLeadershipInternalServerError The supplied request contains invalid fields or could not be parsed (json and non-json bodies). The error's code, message, and cause fields can be inspected for further information

swagger:response clusterTransferLeadershipInternalServerError
*/
type ClusterTransferLeadershipInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *rest_model.APIErrorEnvelope `json:"body,omitempty"`
}

// NewClusterTransferLeadershipInternalServerError creates ClusterTransferLeadershipInternalServerError with default headers values
func NewClusterTransferLeadershipInternalServerError() *ClusterTransferLeadershipInternalServerError {

	return &ClusterTransferLeadershipInternalServerError{}
}

// WithPayload adds the payload to the cluster transfer leadership internal server error response
func (o *ClusterTransferLeadershipInternalServerError) WithPayload(payload *rest_model.APIErrorEnvelope) *ClusterTransferLeadershipInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the cluster transfer leadership internal server error response
func (o *ClusterTransferLeadershipInternalServerError) SetPayload(payload *rest_model.APIErrorEnvelope) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ClusterTransferLeadershipInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
