// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
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
