// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"io"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// Response type is a struct that represents a JSON response with a data field and an optional
// errors field.
// @property Data - The `Data` property is of type `interface{}` which means it can hold any type of
// data. It is used to store the actual data that needs to be returned in the response.
// @property {string} Errors - The "Errors" property is a string that is used to store any error
// messages or error information related to the response. It is marked as "omitempty" in the JSON tag,
// which means that if the value of the "Errors" property is empty or zero, it will be omitted from the
// JSON
type Response struct {
	Data   interface{} `json:"data"`
	Errors string      `json:"errors,omitempty"`
}

// DecodeRespFromJSON decodes a JSON input from a reader into a Response struct
// in Go.
func DecodeRespFromJSON(reader io.Reader) (*Response, error) {
	r := &Response{}
	err := json.NewDecoder(reader).Decode(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// RespGraphQL type is a struct that represents the response data from a GraphQL query, including
// event, event receiver, and event receiver group information.
// @property Data - The `Data` property is a struct that contains the following properties:
// @property {string} Errors - The `Errors` property is a string that is used to store any error
// messages that may occur during the execution of the GraphQL query. It is optional and will only be
// present if there are any errors.
type RespGraphQL struct {
	Data struct {
		Events                   []storage.Event              `json:"events,omitempty"`
		EventReceivers           []storage.EventReceiver      `json:"event_receivers,omitempty"`
		EventReceiverGroups      []storage.EventReceiverGroup `json:"event_receiver_groups,omitempty"`
		CreateEvent              graphql.ID                   `json:"create_event,omitempty"`
		CreateEventReceiver      graphql.ID                   `json:"create_event_receiver,omitempty"`
		CreateEventReceiverGroup graphql.ID                   `json:"create_event_receiver_group,omitempty"`
	} `json:"data"`
	Errors string `json:"errors,omitempty"`
}

// The function `DecodeGraphQLRespFromJSON` decodes a JSON response into a `RespGraphQL` struct.
func DecodeGraphQLRespFromJSON(reader io.Reader) (*RespGraphQL, error) {
	r := &RespGraphQL{}
	err := json.NewDecoder(reader).Decode(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
