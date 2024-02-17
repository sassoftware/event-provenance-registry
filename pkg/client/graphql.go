// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"strings"
)

const (
	constEvents              = `events`
	constEventReceivers      = `event_receivers`
	constEventReceiverGroups = `event_receiver_groups`
	constEvent               = `event`
	constEventReceiver       = `event_receiver`
	constEventReceiverGroup  = `event_receiver_group`
)

// GraphQLRequest struct for graphql request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operationName,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// NewGraphQLSearchRequest creates a new GraphQLRequest
// operation can be events or event_receivers or event_receiver_groups
func NewGraphQLSearchRequest(operation string, params map[string]interface{}, fields []string) *GraphQLRequest {
	// {"query":"query ($erg: FindEventReceiverGroupInput!){event_receiver_groups(event_receiver_group: $erg) {id,name,type,version,description}}","variables":{"erg": {"name":"foobar","version":"1.0.0"}}}
	template := `query ($obj: %s){%s(%s: $obj) {%s}}`
	var query string
	switch operation {
	case constEvents:
		query = fmt.Sprintf(template, `FindEventInput!`, operation, constEvent, strings.Join(fields, ","))
	case constEventReceivers:
		query = fmt.Sprintf(template, `FindEventReceiverInput!`, operation, constEventReceiver, strings.Join(fields, ","))
	case constEventReceiverGroups:
		query = fmt.Sprintf(template, `FindEventReceiverGroupInput!`, operation, constEventReceiverGroup, strings.Join(fields, ","))
	}
	variables := map[string]interface{}{
		"obj": params,
	}
	return &GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
}

func NewGraphQLMutationRequest(operation string, params map[string]interface{}) *GraphQLRequest {
	// {"query":"mutation ($er: CreateEventReceiverInput!){create_event_receiver(event_receiver: $er)}","variables":{"er": {"name":"foobar","version":"1.0.0","description":"foobar is the description","type": "foobar.test", "schema" : "{}"}}}' http://localhost:8042/api/v1/graphql/query
	template := `mutation ($obj: %s){%s(%s: $obj)}`
	var query string
	switch operation {
	case `create_event`:
		query = fmt.Sprintf(template, `CreateEventInput!`, `create_event`, constEvent)
	case `create_event_receiver`:
		query = fmt.Sprintf(template, `CreateEventReceiverInput!`, `create_event_receiver`, constEventReceiver)
	case `create_event_receiver_group`:
		query = fmt.Sprintf(template, `CreateEventReceiverGroupInput!`, `create_event_receiver_group`, constEventReceiverGroup)
	}
	variables := map[string]interface{}{
		"obj": params,
	}
	return &GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
}
