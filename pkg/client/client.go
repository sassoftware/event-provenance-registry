// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
)

var logger = utils.MustGetLogger("client", "client")

// ensure it implements the interface
var _ Contract = &Client{}

// Contract handles communications with the EPR service
type Contract interface {
	CreateEvent(e *storage.Event) (string, error)
	CreateEventReceiver(er *storage.EventReceiver) (string, error)
	CreateEventReceiverGroup(erg *storage.EventReceiverGroup) (string, error)
	ModifyEventReceiverGroup(erg *storage.EventReceiverGroup) (string, error)
	Search(queryName string, queryFor string, params map[string]interface{}, fields []string) (string, error)
	SearchEvents(params map[string]interface{}, fields []string) ([]storage.Event, error)
	SearchEventReceivers(params map[string]interface{}, fields []string) ([]storage.EventReceiver, error)
	SearchEventReceiverGroups(params map[string]interface{}, fields []string) ([]storage.EventReceiverGroup, error)
	CheckReadiness() (bool, error)
	CheckLiveness() (bool, error)
	CheckStatus() (string, error)
	GetEndpoint(end string) (string, error)
}

// Client is a struct for EPR Client configuration
type Client struct {
	url        string
	apiVersion string
	health     string
	client     *http.Client
}

// New returns a new instance of Client struct. Requires a URL to an
// instance of EPR. Use options functions for setting specific parameters.
func New(url string) (*Client, error) {
	client := &http.Client{}
	c := &Client{
		url:        url,
		apiVersion: "/api/v1",
		health:     "/healthz",
		client:     client,
	}

	return c, nil
}

// DoGet makes a GET request to the endpoint and returns the response
func (c *Client) DoGet(endpoint string) (string, error) {
	return c.doReq(http.MethodGet, endpoint, nil)
}

// DoPost makes a POST request to the endpoint and returns the response
func (c *Client) DoPost(endpoint string, payload []byte) (string, error) {
	return c.doReq(http.MethodPost, endpoint, payload)
}

// DoDelete makes a DELETE request to the endpoint and returns the response
func (c *Client) DoDelete(endpoint string, payload []byte) (string, error) {
	return c.doReq(http.MethodDelete, endpoint, payload)
}

// DoPatch makes a PATCH request to the endpoint and returns the response
func (c *Client) DoPatch(endpoint string, payload []byte) (string, error) {
	return c.doReq(http.MethodPatch, endpoint, payload)
}

// DoPut makes a PUT request to the endpoint and returns the response
func (c *Client) DoPut(endpoint string, payload []byte) (string, error) {
	return c.doReq(http.MethodPut, endpoint, payload)
}

// doReq makes a given request type to the endpoint and returns a response
func (c *Client) doReq(reqType string, endpoint string, payload []byte) (string, error) {
	var req *http.Request
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	switch reqType {
	case http.MethodGet:
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	case http.MethodPost:
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	case http.MethodDelete:
		req, err = http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, bytes.NewReader(payload))
	case http.MethodPatch:
		req, err = http.NewRequestWithContext(ctx, http.MethodPatch, endpoint, bytes.NewReader(payload))
	case http.MethodPut:
		req, err = http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(payload))
	default:
		return "", fmt.Errorf("request type %s not supported", reqType)
	}
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= http.StatusBadRequest || resp.StatusCode < http.StatusOK {
		r := &Response{}
		err := json.Unmarshal(content, r)
		if err != nil {
			return string(content), fmt.Errorf("request returned status code %d", resp.StatusCode)
		}
		return string(content), fmt.Errorf("request returned status code %d (%s)", resp.StatusCode, r.Errors)
	}

	return string(content), nil
}

// GetEndpoint formats endoint url
func (c *Client) GetEndpoint(end string) (string, error) {
	s, err := url.JoinPath(c.url, c.apiVersion, end)
	if err != nil {
		return "", err
	}
	return s, nil
}

// getHealthEndpoint formats health endpoints
func (c *Client) getHealthEndpoint(end string) (string, error) {
	s, err := url.JoinPath(c.url, c.health, end)
	if err != nil {
		return "", err
	}
	return s, nil
}

// getGraphQLEndpoint formats graphql endpoints
func (c *Client) getGraphQLEndpoint() (string, error) {
	s, err := url.JoinPath(c.url, c.apiVersion, `graphql`)
	if err != nil {
		return "", err
	}
	return s, nil
}

// getGraphQLEndpointQuery formats graphql endpoints
func (c *Client) getGraphQLEndpointQuery() (string, error) {
	s, err := url.JoinPath(c.url, c.apiVersion, `graphql`, `query`)
	if err != nil {
		return "", err
	}
	return s, nil
}
