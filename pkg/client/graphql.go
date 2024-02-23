// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"strings"
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
	case eventsQuery:
		query = fmt.Sprintf(template, findEventInputQuery, operation, eventQuery, strings.Join(fields, ","))
	case eventReceiversQuery:
		query = fmt.Sprintf(template, findEventReceiverInputQuery, operation, eventReceiverQuery, strings.Join(fields, ","))
	case eventReceiverGroupsQuery:
		query = fmt.Sprintf(template, findEventReceiverGroupInputQuery, operation, eventReceiverGroupQuery, strings.Join(fields, ","))
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
		query = fmt.Sprintf(template, createEventInputQuery, createEventQuery, eventQuery)
	case `create_event_receiver`:
		query = fmt.Sprintf(template, createEventReceiverInputQuery, createEventReceiverQuery, eventReceiverQuery)
	case `create_event_receiver_group`:
		query = fmt.Sprintf(template, createEventReceiverGroupInputQuery, createEventReceiverGroupQuery, eventReceiverGroupQuery)
	}
	variables := map[string]interface{}{
		"obj": params,
	}
	return &GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
}
