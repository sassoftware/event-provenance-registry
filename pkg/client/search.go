// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// QueryEventReceiverByID for querying EPR for EventReceiver by ID. Returns nil if no results.
// fields := []string{"id", "name", "type", "version", "description", "schemaVersion", "maintainer", "fingerprint", "tags", "encoding", "schema", "timestamp", "metadata"}
func (c *Client) QueryEventReceiverByID(id string, fields []string) (*storage.EventReceiver, error) {
	params := map[string]interface{}{
		"id": id,
	}

	gates, err := c.SearchEventReceiversObj(params, fields)
	if err != nil {
		return nil, err
	}
	switch len(gates) {
	case 0:
		return nil, nil
	case 1:
		return &gates[0], nil
	default:
		// This shouldn't happen, but return first entry if it does
		logger.Errorf("two gates found with same ID: '%s'", id)
		return &gates[0], nil
	}
}

// QueryEventByID for querying EPR for Event by ID. Returns nil if no results.
// fields := []string{"id", "name", "version", "description", "release", "platform_id", "package", "mutipass", "revoked", "timestamp", "gate_id"}
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
		logger.Errorf("two events found with same ID: '%s'", id)
		return &events[0], nil
	}
}

// QueryEventReceiverGroupByID for querying EPR for EventReceiverGroup by ID. Returns nil if no results.
// fields := []string{"id", "name", "type", "version", "description", "fingerprint", "tags", "gates", "timestamp", "metadata", "disabled"}
func (c *Client) QueryEventReceiverGroupByID(id string, fields []string) (*storage.EventReceiverGroup, error) {
	params := map[string]interface{}{
		"id": id,
	}

	stages, err := c.SearchEventReceiverGroupsObj(params, fields)
	if err != nil {
		return nil, err
	}

	switch len(stages) {
	case 0:
		return nil, nil
	case 1:
		return &stages[0], nil
	default:
		// This shouldn't happen, but return first entry if it does
		logger.Errorf("two stages found with same ID: '%s'", id)
		return &stages[0], nil
	}
}

// getGraphqlEndpoint use for getting endpoints
func (c *Client) getGraphqlEndpoint() string {
	return fmt.Sprintf(`%s%s/graphql`, c.url, c.apiVersion)
}

// SearchEventReceiverGroups searches for the stages based on params
func (c *Client) SearchEventReceiverGroups(params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, "FindEventReceiverGroup", "stages", endpoint, fields)
}

// SearchEventReceiverGroupsObj unmarshals the result of SearchEventReceiverGroups into a list of EventReceiverGroups
func (c *Client) SearchEventReceiverGroupsObj(params map[string]interface{}, fields []string) ([]storage.EventReceiverGroup, error) {
	response, err := c.SearchEventReceiverGroups(params, fields)
	if err != nil {
		return nil, err
	}

	respObj, err := responses.DecodeRespGraphQLEventReceiverGroupsFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	// Check for presence of errors in respObj from searching stages
	if respObj.Errors != nil {
		return nil, fmt.Errorf("when searching for stage returned: errors: %s ", respObj.Errors)
	}

	stageList := []storage.EventReceiverGroup{}
	for _, gqlEventReceiverGroup := range respObj.Data.EventReceiverGroups {
		erg, err := gqlEventReceiverGroup.ToEventReceiverGroup()
		if err != nil {
			return stageList, err
		}

		stageList = append(stageList, *erg)
	}

	return stageList, nil
}

// SearchEventReceivers searches for gates based on params
func (c *Client) SearchEventReceivers(params map[string]interface{}, fields []string) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.searchQuery(params, "FindEventReceiver", "gates", endpoint, fields)
}

// SearchEventReceiversObj unmarshals the result of SearchEventReceivers into a list of EventReceivers
func (c *Client) SearchEventReceiversObj(params map[string]interface{}, fields []string) ([]storage.EventReceiver, error) {
	response, err := c.SearchEventReceivers(params, fields)
	if err != nil {
		return nil, err
	}

	respObj, err := responses.DecodeRespGraphQLEventReceiversFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	// Check for presence of errors in respObj from searching gate
	if respObj.Errors != nil {
		return nil, fmt.Errorf("when searching for gate returned: errors: %s ", respObj.Errors)
	}

	gateList := []storage.EventReceiver{}
	for _, gqlEventReceiver := range respObj.Data.EventReceivers {
		er, err := gqlEventReceiver.ToEventReceiver()
		if err != nil {
			return gateList, err
		}

		gateList = append(gateList, *er)
	}
	return gateList, nil
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
	respObj, err := responses.DecodeRespGraphQLEventsFromJSON(strings.NewReader(response))
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

// searchQuery implements the searching for the gatekeeper
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

// IDSearch searches for the given gates, stages, events based on params
func (c *Client) IDSearch(params map[string]interface{}) (string, error) {
	endpoint := c.getGraphqlEndpoint()
	return c.idSearchQuery(params, endpoint)
}

// idSearchQuery implements the searching by ids for the gatekeeper
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
