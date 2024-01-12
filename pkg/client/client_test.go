// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetGraphQLEndpoint(t *testing.T) {
	c := &Client{
		url:        "http://example.com",
		apiVersion: "v1",
	}

	expected := "http://example.com/v1/graphql"

	result, err := c.getGraphQLEndpoint()
	assert.NilError(t, err, "Unexpected error: %v", err)
	assert.Equal(t, expected, result, "Expected %s, but got %s", expected, result)
}

func TestGetGraphQLEndpointInvalidURL(t *testing.T) {
	c := &Client{
		url:        "http://example.com",
		apiVersion: "v1",
	}

	// Set an invalid URL to trigger an error
	c.url = ":invalid_url"

	_, err := c.getGraphQLEndpoint()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func TestGetGraphQLEndpointQuery(t *testing.T) {
	c := &Client{
		url:        "http://example.com",
		apiVersion: "v1",
	}

	expected := "http://example.com/v1/graphql/query"

	result, err := c.getGraphQLEndpointQuery()
	assert.NilError(t, err, "Unexpected error: %v", err)
	assert.Equal(t, expected, result, "Expected %s, but got %s", expected, result)
}

func TestGetGraphQLEndpointQueryInvalidURL(t *testing.T) {
	c := &Client{
		url:        "http://example.com",
		apiVersion: "v1",
	}

	// Set an invalid URL to trigger an error
	c.url = ":invalid_url"

	_, err := c.getGraphQLEndpointQuery()
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}
