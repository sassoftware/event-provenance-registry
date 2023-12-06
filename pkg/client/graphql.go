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

// NewGraphQLRequest returns a new instance of GraphQLRequest
func NewGraphQLRequest(operation string, lookFor string, fields []string, params map[string]interface{}) *GraphQLRequest {
	template := `query %s(%s){%s(%s) {%s}}`
	varDefs := ""
	selSets := ""
	for k, v := range params {
		varDef, selSet := formatValues(k, v)
		varDefs += varDef
		selSets += selSet
	}
	varDefs = strings.Trim(varDefs, ",")
	selSets = strings.Trim(selSets, ",")

	query := fmt.Sprintf(template, operation, varDefs, lookFor, selSets, strings.Join(fields, ","))
	return &GraphQLRequest{
		Query:     query,
		Variables: params,
	}
}

func formatValues(k string, v interface{}) (string, string) {
	switch v.(type) {
	case string:
		return fmt.Sprintf(`$%s: String,`, k), fmt.Sprintf(`%s: $%s,`, k, k)
	case []string:
		return fmt.Sprintf(`$%s: [String],`, k), fmt.Sprintf(`%s: $%s,`, k, k)
	case int:
		return fmt.Sprintf(`$%s: Int,`, k), fmt.Sprintf(`%s: $%s,`, k, k)
	}
	return "", ""
}

// NewGraphQLRequestIds returns a new instance of GraphQLRequest given params
func NewGraphQLRequestIds(params map[string]interface{}) *GraphQLRequest {
	query := `query FindAny($ulids: [String]){any (ulids: $ulids) { events receivers groups }}`
	return &GraphQLRequest{
		Query:     query,
		Variables: params,
	}
}
