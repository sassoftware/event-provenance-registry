// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestNewGraphQLRequest(t *testing.T) {
	queryName := "FindEventReceivers"
	lookFor := "event_receivers"
	fields := []string{"id", "name", "version", "type"}
	params := map[string]interface{}{
		"id": "01HKMQM136XW7JYP2293N4EBR4",
	}
	expected := `query FindEventReceivers($id: ID!){event_receivers(id: $id) {id,name,version,type}}`
	variables := `01HKMQM136XW7JYP2293N4EBR4`

	req := NewGraphQLRequest(queryName, lookFor, params, fields)
	assert.Equal(t, expected, req.Query, "The generated Query did not match expected")
	assert.Equal(t, variables, req.Variables["id"], "The generated Variables did not match expected")
}

func TestNewGraphQLRequestWithList(t *testing.T) {
	queryName := "FindEventReceiverGroups"
	lookFor := "groups"
	fields := []string{"id", "name", "version", "type"}
	params := map[string]interface{}{
		"id":      "01HKMQM136XW7JYP2293N4EBR4",
		"version": "1.0.0",
	}
	expected := `query FindEventReceiverGroups($id: ID!,$version: String){groups(id: $id,version: $version) {id,name,version,type}}`
	req := NewGraphQLRequest(queryName, lookFor, params, fields)
	assert.Equal(t, req.Query, expected, "The generated Query did not match expected")
}

func TestNewGraphQLRequestWithInt(t *testing.T) {
	queryName := "FindEvents"
	lookFor := "events"
	fields := []string{"name", "id"}
	params := map[string]interface{}{
		"name":  "my-event",
		"start": 1,
	}
	expected := []string{`query FindEvents($name: String,$start: Int){events(name: $name,start: $start) {name,id}}`,
		`query FindEvent($start: Int,$name: String){events(start: $start,name: $name) {name,id}}`}
	req := NewGraphQLRequest(queryName, lookFor, params, fields)
	assert.Check(t, customStringCompare(req.Query, expected))
}

func TestComplexNewGraphQLRequest(t *testing.T) {
	queryName := "foo"
	lookFor := "bar"
	fields := []string{"name", "id", "success"}
	params := map[string]interface{}{
		"id":      "abc123",
		"name":    "wally",
		"success": true,
	}

	expected := []string{`query foo($success: Bool,$id: ID!,$name: String){bar(success: $success,id: $id,name: $name) {name,id,success}}`,
		`query foo($id: ID!,$name: String,$success: Bool){bar(id: $id,name: $name,success: $success) {name,id,success}}`,
		`query foo($name: String,$success: Bool,$id: ID!){bar(name: $name,success: $success,id: $id) {name,id,success}}`}
	req := NewGraphQLRequest(queryName, lookFor, params, fields)
	assert.Check(t, customStringCompare(req.Query, expected))
}

func customStringCompare(query string, options []string) cmp.Comparison {
	return func() cmp.Result {
		for _, possible := range options {
			if possible == query {
				return cmp.ResultSuccess
			}
		}
		return cmp.ResultFailure(fmt.Sprintf("%q did not match available options %v", query, options))
	}
}
