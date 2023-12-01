// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// QueryEventByID for querying EPR for Event by ID. Returns nil if no results.
// fields := []string{"id", "name", "version", "description", "release", "platform_id", "package", "mutipass", "revoked", "timestamp", "eventReceiver_id"}
func (c *Client) QueryEventByID(id string, fields []string) (*storage.Event, error) {
	params := map[string]interface{}{
		"id": id,
	}

	events, err := c.SearchEventsObj(params, fields)
	if err != nil {
		return nil, err
	}

	switch len(events) {
	case 0:
		return nil, nil
	case 1:
		return &events[0], nil
	default:
		// This shouldn't happen, but return first entry if it does
		logger.Error(nil, "two events found with same", "ID:", id)
		return &events[0], nil
	}
}

// QueryEventReceiverByID for querying EPR for EventReceiver by ID. Returns nil if no results.
// fields := []string{"id", "name", "type", "version", "description", "fingerprint", "schema",}
func (c *Client) QueryEventReceiverByID(id string, fields []string) (*storage.EventReceiver, error) {
	params := map[string]interface{}{
		"id": id,
	}

	eventReceivers, err := c.SearchEventReceiversObj(params, fields)
	if err != nil {
		return nil, err
	}
	switch len(eventReceivers) {
	case 0:
		return nil, nil
	case 1:
		return &eventReceivers[0], nil
	default:
		// This shouldn't happen, but return first entry if it does
		logger.Error(nil, "two eventReceivers found with same", "ID:", id)
		return &eventReceivers[0], nil
	}
}

// QueryEventReceiverGroupByID for querying EPR for EventReceiverGroup by ID. Returns nil if no results.
// fields := []string{"id", "name", "type", "version", "description", "fingerprint", "disabled"}
func (c *Client) QueryEventReceiverGroupByID(id string, fields []string) (*storage.EventReceiverGroup, error) {
	params := map[string]interface{}{
		"id": id,
	}

	eventReceiverGroups, err := c.SearchEventReceiverGroupsObj(params, fields)
	if err != nil {
		return nil, err
	}

	switch len(eventReceiverGroups) {
	case 0:
		return nil, nil
	case 1:
		return &eventReceiverGroups[0], nil
	default:
		// This shouldn't happen, but return first entry if it does
		logger.Error(nil, "two eventReceiverGroups found with same", "ID:", id)
		return &eventReceiverGroups[0], nil
	}
}

// getGraphqlEndpoint use for getting endpoints
func (c *Client) getGraphqlEndpoint() string {
	return fmt.Sprintf(`%s%s/graphql`, c.url, c.apiVersion)
}

// SearchEventReceiverGroups searches for the eventReceiverGroups based on params
func (c *Client) SearchEventReceiverGroups(params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, "FindEventReceiverGroup", "eventReceiverGroups", endpoint, fields)
}

// SearchEventReceiverGroupsObj unmarshals the result of SearchEventReceiverGroups into a list of EventReceiverGroups
func (c *Client) SearchEventReceiverGroupsObj(params map[string]interface{}, fields []string) ([]storage.EventReceiverGroup, error) {
	response, err := c.SearchEventReceiverGroups(params, fields)
	if err != nil {
		return nil, err
	}

	respObj, err := DecodeRespFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	// Check for presence of errors in respObj from searching eventReceiverGroups
	if respObj.Errors != nil {
		return nil, fmt.Errorf("when searching for eventReceiverGroup returned: errors: %s ", respObj.Errors)
	}

	eventReceiverGroupList := []storage.EventReceiverGroup{}
	for _, gqlEventReceiverGroup := range respObj.Data.EventReceiverGroups {
		erg, err := gqlEventReceiverGroup.ToEventReceiverGroup()
		if err != nil {
			return eventReceiverGroupList, err
		}

		eventReceiverGroupList = append(eventReceiverGroupList, *erg)
	}

	return eventReceiverGroupList, nil
}

// SearchEventReceivers searches for eventReceivers based on params
func (c *Client) SearchEventReceivers(params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, "FindEventReceiver", "eventReceivers", endpoint, fields)
}

// SearchEventReceiversObj unmarshals the result of SearchEventReceivers into a list of EventReceivers
func (c *Client) SearchEventReceiversObj(params map[string]interface{}, fields []string) ([]storage.EventReceiver, error) {
	response, err := c.SearchEventReceivers(params, fields)
	if err != nil {
		return nil, err
	}

	respObj, err := DecodeRespFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	// Check for presence of errors in respObj from searching eventReceiver
	if respObj.Errors != nil {
		return nil, fmt.Errorf("when searching for eventReceiver returned: errors: %s ", respObj.Errors)
	}

	eventReceiverList := []storage.EventReceiver{}
	for _, gqlEventReceiver := range respObj.Data.EventReceivers {
		er, err := gqlEventReceiver.ToEventReceiver()
		if err != nil {
			return eventReceiverList, err
		}

		eventReceiverList = append(eventReceiverList, *er)
	}
	return eventReceiverList, nil
}

// SearchEvents searches for events based on params
func (c *Client) SearchEvents(params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, "FindEvent", "events", endpoint, fields)
}

// SearchEventsObj unmarshals the result of SearchEvents into a list of Events
func (c *Client) SearchEventsObj(params map[string]interface{}, fields []string) ([]storage.Event, error) {
	response, err := c.SearchEvents(params, fields)
	if err != nil {
		return nil, err
	}
	respObj, err := DecodeRespFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	// Check for presence of errors in responseObj from searching event
	if respObj.Errors != nil {
		return nil, fmt.Errorf("when searching for event returned: errors: %s ", respObj.Errors)
	}

	eventList := []storage.Event{}
	for _, gqlRec := range respObj.Data.Events {
		e, err := gqlRec.ToEvent()
		if err != nil {
			return eventList, err
		}
		eventList = append(eventList, *e)
	}

	return eventList, nil
}

// Search searches for the given queryFor based on params
func (c *Client) Search(queryName string, queryFor string, params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, queryName, queryFor, endpoint, fields)
}

// searchQuery implements the searching for the eventReceiverkeeper
func (c *Client) searchQuery(params map[string]interface{}, queryName, queryFor, endpoint string, fields []string) (string, error) {
	s := 
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

// idSearchQuery implements the searching by ids for the eventReceiverkeeper
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
