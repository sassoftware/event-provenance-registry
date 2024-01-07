// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

// TestGraphQLRequest test
func TestGraphQLRequest(t *testing.T) {
	t.Run("testNewGraphQLRequest", testNewGraphQLRequest())
	t.Run("testComplexNewGraphQLRequest", testComplexNewGraphQLRequest())
	t.Run("testNewGraphQLRequestWithList", testNewGraphQLRequestWithList())
	t.Run("testNewGraphQLRequestWithInt", testNewGraphQLRequestWithInt())
}

func testNewGraphQLRequest() func(t *testing.T) {
	return func(t *testing.T) {
		queryName := "FindEventReceiver"
		lookFor := "receivers"
		fields := []string{"name", "version", "release"}
		params := map[string]interface{}{
			"name": "A",
		}
		expected := `query FindEventReceiver($name: String){receivers(name: $name) {name,version,release}}`

		req := NewGraphQLRequest(queryName, lookFor, params, fields)
		assert.Equal(t, req.Query, expected, "The generated Query did not match expected")
	}
}

func testNewGraphQLRequestWithList() func(t *testing.T) {
	return func(t *testing.T) {
		queryName := "FindEventReceiverGroup"
		lookFor := "groups"
		fields := []string{"name", "id"}
		params := map[string]interface{}{
			"version": "1.0.0",
		}
		expected := `query FindEventReceiverGroup($version: String){groups(version: $version) {name,id}}`
		req := NewGraphQLRequest(queryName, lookFor, params, fields)
		assert.Equal(t, req.Query, expected, "The generated Query did not match expected")
	}
}

func testNewGraphQLRequestWithInt() func(t *testing.T) {
	return func(t *testing.T) {
		queryName := "FindEvent"
		lookFor := "events"
		fields := []string{"name", "id"}
		params := map[string]interface{}{
			"name":  "my-event",
			"start": 1,
		}
		expected := []string{`query FindEvent($name: String,$start: Int){events(name: $name,start: $start) {name,id}}`,
			`query FindEvent($start: Int,$name: String){events(start: $start,name: $name) {name,id}}`}
		req := NewGraphQLRequest(queryName, lookFor, params, fields)
		assert.Check(t, customStringCompare(req.Query, expected))
	}
}

func testComplexNewGraphQLRequest() func(t *testing.T) {
	return func(t *testing.T) {
		queryName := "foo"
		lookFor := "bar"
		fields := []string{"name", "id", "success"}
		params := map[string]interface{}{
			"id":      "abc123",
			"name":    "wally",
			"success": true,
		}

		expected := []string{`query foo($success: Bool,$id: String,$name: String){bar(success: $success,id: $id,name: $name) {name,id,success}}`,
			`query foo($id: String,$name: String,$success: Bool){bar(id: $id,name: $name,success: $success) {name,id,success}}`}
		req := NewGraphQLRequest(queryName, lookFor, params, fields)
		assert.Check(t, customStringCompare(req.Query, expected))
	}
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
