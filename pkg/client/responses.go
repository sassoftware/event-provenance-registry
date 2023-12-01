// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"io"
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

// // RespGraphQL type represents a response from a GraphQL API that includes data for events, event
// // receivers, and event receiver groups.
// // @property Data - The `Data` property is a struct that contains three arrays: `Events`,
// // `EventReceivers`, and `EventReceiverGroups`. Each of these arrays contains objects that have
// // specific properties related to events, event receivers, and event receiver groups respectively.
// // @property Errors - The `Errors` property is of type `types.JSON` and is used to store any errors
// // that occur during the execution of the GraphQL query. It is an optional property, meaning it may be
// // omitted if there are no errors.
// type RespGraphQL struct {
// 	Data struct {
// 		Events              []ERespGraphQL   `json:"events"`
// 		EventReceivers      []ERRespGraphQL  `json:"event_receivers"`
// 		EventReceiverGroups []ERGRespGraphQL `json:"event_receiver_groups"`
// 	} `json:"data"`
// 	Errors types.JSON `json:"errors,omitempty"`
// }

// DecodeRespFromJSON decodes a JSON input from a reader into a RespGraphQ struct
// in Go.
func DecodeRespFromJSON(reader io.Reader) (*Response, error) {
	r := &Response{}
	err := json.NewDecoder(reader).Decode(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
