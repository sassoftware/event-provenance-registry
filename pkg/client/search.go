// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// getGraphqlEndpoint use for getting endpoints
func (c *Client) getGraphqlEndpoint() string {
	return fmt.Sprintf(`%s%s/graphql`, c.url, c.apiVersion)
}

// SearchEventReceiverGroups searches for the eventReceiverGroups based on params
func (c *Client) SearchEventReceiverGroups(params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, "FindEventReceiverGroup", "eventReceiverGroups", endpoint, fields)
}

// SearchEventReceiverGroupsObj function is used to search for event receiver groups based on the
// provided parameters and retrieve the results as a list of `storage.EventReceiverGroup` objects.
func (c *Client) SearchEventReceiverGroupsObj(params map[string]interface{}, fields []string) ([]storage.EventReceiverGroup, error) {
	response, err := c.SearchEventReceiverGroups(params, fields)
	if err != nil {
		return nil, err
	}

	respObj, err := DecodeGraphQLRespFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	// Check for presence of errors in respObj from searching eventReceiverGroups
	if respObj.Errors != "" {
		return nil, fmt.Errorf("when searching for eventReceiverGroup returned: errors: %s ", respObj.Errors)
	}

	return respObj.Data.EventReceiverGroups, nil
}

// SearchEventReceivers searches for eventReceivers based on params
func (c *Client) SearchEventReceivers(params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, "FindEventReceiver", "eventReceivers", endpoint, fields)
}

// SearchEventReceiversObj function is used to search for event receivers based on the provided
// parameters and retrieve the results as a list of `storage.EventReceiver` objects.
func (c *Client) SearchEventReceiversObj(params map[string]interface{}, fields []string) ([]storage.EventReceiver, error) {
	response, err := c.SearchEventReceivers(params, fields)
	if err != nil {
		return nil, err
	}

	respObj, err := DecodeGraphQLRespFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	// Check for presence of errors in respObj from searching eventReceiver
	if respObj.Errors != "" {
		return nil, fmt.Errorf("when searching for eventReceiver returned: errors: %s ", respObj.Errors)
	}

	return respObj.Data.EventReceivers, nil
}

// SearchEvents searches for events based on params
func (c *Client) SearchEvents(params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, "FindEvent", "events", endpoint, fields)
}

// SearchEventsObj function is used to search for events based on the provided parameters and
// retrieve the results as a list of `storage.Event` objects.
func (c *Client) SearchEventsObj(params map[string]interface{}, fields []string) ([]storage.Event, error) {
	response, err := c.SearchEvents(params, fields)
	if err != nil {
		return nil, err
	}
	respObj, err := DecodeGraphQLRespFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	// Check for presence of errors in responseObj from searching event
	if respObj.Errors != "" {
		return nil, fmt.Errorf("when searching for event returned: errors: %s ", respObj.Errors)
	}

	return respObj.Data.Events, nil
}

// Search searches for the given queryFor based on params
func (c *Client) Search(queryName string, queryFor string, params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, queryName, queryFor, endpoint, fields)
}

// searchQuery implements the searching
func (c *Client) searchQuery(params map[string]interface{}, queryName, queryFor, endpoint string, fields []string) (string, error) {
	gqlBody := NewGraphQLRequest(queryName, queryFor, fields, params)
	enc, err := json.Marshal(gqlBody)
	if err != nil {
		return "", err
	}

	content, err := c.DoPost(endpoint, enc)
	if err != nil {
		return "", err
	}

	return content, nil
}

// IDSearch searches for the given eventReceivers, eventReceiverGroups, events based on params
func (c *Client) IDSearch(params map[string]interface{}) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.idSearchQuery(params, endpoint)
}

// idSearchQuery implements the searching by ids
func (c *Client) idSearchQuery(params map[string]interface{}, endpoint string) (string, error) {
	gqlBody := NewGraphQLRequestIds(params)
	enc, err := json.Marshal(gqlBody)
	if err != nil {
		return "", err
	}

	content, err := c.DoPost(endpoint, enc)
	if err != nil {
		return "", err
	}

	return content, nil
}
