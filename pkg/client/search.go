// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// Search searches for the given queryFor based on params
func (c *Client) Search(queryName string, queryFor string, params map[string]interface{}, fields []string) (string, error) {
	endpoint, err := c.getGraphQLEndpoint()
	if err != nil {
		return "", err
	}

	gqlBody := NewGraphQLRequest(queryName, queryFor, params, fields)
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

func (c *Client) SearchEvents(params map[string]interface{}, fields []string) ([]storage.Event, error) {
	response, err := c.Search("FindEvents", "", params, fields)
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

	return respObj.Data.Events, nil
}

func (c *Client) SearchEventReceivers(params map[string]interface{}, fields []string) ([]storage.EventReceiver, error) {
	response, err := c.Search("FindEventReceivers", "", params, fields)
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

	return respObj.Data.EventReceivers, nil
}

func (c *Client) SearchEventReceiverGroups(params map[string]interface{}, fields []string) ([]storage.EventReceiverGroup, error) {
	response, err := c.Search("FindEventReceiverGroups", "", params, fields)
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
