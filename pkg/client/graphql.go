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
func NewGraphQLRequest(operation string, lookFor string, params map[string]interface{}, fields []string) *GraphQLRequest {
	template := `query %s(%s){%s(%s) {%s}}`
	varDefs := ""
	selSets := ""
	for k, v := range params {
		varDef, selSet := schemaValues(k, v)
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
	case bool:
		return fmt.Sprintf(`$%s: Bool,`, k), fmt.Sprintf(`%s: $%s,`, k, k)
	}
	return "", ""
}

func schemaValues(k string, v interface{}) (string, string) {
	switch k {
	case `id`:
		return fmt.Sprintf(`$%s: ID!,`, k), fmt.Sprintf(`%s: $%s,`, k, k)
	default:
		return formatValues(k, v)
	}
}
