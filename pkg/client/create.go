// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// CreateEvent used to create and Event
func (c *Client) CreateEvent(e *storage.Event) (string, error) {
	endpoint := c.GetEndpoint("/events")

	enc, err := json.Marshal(e)
	if err != nil {
		return "", err
	}

	content, err := c.DoPost(endpoint, enc)
	if err != nil {
		return content, err
	}

	return content, nil
}

// CreateEventReceiver used to create an EventReceiver
func (c *Client) CreateEventReceiver(er *storage.EventReceiver) (string, error) {
	endpoint := c.GetEndpoint("/receivers")

	enc, err := json.Marshal(er)
	if err != nil {
		return "", err
	}

	content, err := c.DoPost(endpoint, enc)
	if err != nil {
		return content, err
	}

	return content, nil
}

// CreateEventReceiverGroup used to create an EventReceiverGroup
func (c *Client) CreateEventReceiverGroup(erg *storage.EventReceiverGroup) (string, error) {
	endpoint := c.GetEndpoint("/groups")

	enc, err := json.Marshal(erg)
	if err != nil {
		return "", nil
	}

	content, err := c.DoPost(endpoint, enc)
	if err != nil {
		return content, err
	}

	return content, nil
}

// ModifyEventReceiverGroup takes a EventReceiverGroup object and updates the "Disabled" field in the EPR based on the EventReceiverGroup ID. This
// function returns a JSON blob with the ID of the EventReceiverGroup it modified.
func (c *Client) ModifyEventReceiverGroup(erg *storage.EventReceiverGroup) (string, error) {
	endpoint := c.GetEndpoint("/groups/" + string(erg.ID))

	enc, err := json.Marshal(erg)
	if err != nil {
		return "", err
	}

	content, err := c.DoPut(endpoint, enc)
	if err != nil {
		return content, err
	}

	return content, nil
}
