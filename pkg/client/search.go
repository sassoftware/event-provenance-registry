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
	endpoint, err := c.getGraphQLEndpointQuery()
	if err != nil {
		return "", err
	}

	gqlBody := NewGraphQLRequest(queryName, queryFor, params, fields)
	enc, err := json.Marshal(gqlBody)
	if err != nil {
		return "", err
	}

	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Printf("Query: %s\n", enc)
	fmt.Printf(`curl -X POST -H "content-type:application/json" -d '%s' %s`, enc, endpoint)
	fmt.Println("")

	content, err := c.DoPost(endpoint, enc)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (c *Client) SearchEvents(params map[string]interface{}, fields []string) ([]storage.Event, error) {
	response, err := c.Search("FindEvents", "events", params, fields)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response: %s\n", response)

	respObj, err := DecodeGraphQLRespFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	fmt.Printf("ResponseObj: +%v\n", respObj)

	// Check for presence of errors in respObj from searching eventReceiverGroups
	if respObj.Errors != nil {
		return nil, fmt.Errorf("when searching for Event returned: errors: %s ", respObj.Errors)
	}

	return respObj.Data.Events, nil
}

func (c *Client) SearchEventReceivers(params map[string]interface{}, fields []string) ([]storage.EventReceiver, error) {
	response, err := c.Search("FindEventReceivers", "event_receivers", params, fields)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response: %s\n", response)

	respObj, err := DecodeGraphQLRespFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	fmt.Printf("ResponseObj: +%v\n", respObj)

	// Check for presence of errors in respObj from searching eventReceiver
	if respObj.Errors != nil {
		return nil, fmt.Errorf("when searching for eventReceiver returned: errors: %s ", respObj.Errors)
	}

	return respObj.Data.EventReceivers, nil
}

func (c *Client) SearchEventReceiverGroups(params map[string]interface{}, fields []string) ([]storage.EventReceiverGroup, error) {
	response, err := c.Search("FindEventReceiverGroups", "event_receiver_groups", params, fields)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response: %s\n", response)

	respObj, err := DecodeGraphQLRespFromJSON(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	fmt.Printf("ResponseObj: +%v\n", respObj)

	// Check for presence of errors in respObj from searching eventReceiverGroups
	if respObj.Errors != nil {
		return nil, fmt.Errorf("when searching for eventReceiverGroup returned: errors: %s ", respObj.Errors)
	}

	return respObj.Data.EventReceiverGroups, nil
}
