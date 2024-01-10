//go:build functional

package functional

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"gotest.tools/v3/assert"
)

const groupURI = "http://localhost:8042/api/v1/groups/"

type getGroupResponse struct {
	Data   []storage.EventReceiverGroup
	Errors []string
}

type postGroupResponse struct {
	Data   string
	Errors []string
}

func TestGetGroup(t *testing.T) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	t.Run("group exists", func(t *testing.T) {
		resp, err := client.Get(groupURI + "group-a")
		assert.NilError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		var body getGroupResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NilError(t, err)

		assert.Equal(t, len(body.Errors), 0, "expect no resp body errors")
		assert.Equal(t, len(body.Data), 1, "expect receiver group in resp body")
		assertGroupsEqual(t, body.Data[0], storage.EventReceiverGroup{
			ID:               "group-a",
			Name:             "deploy to prod",
			Type:             "deploy.prod",
			Version:          "1.1.1",
			Description:      "deploy to production when artifact vetted",
			Enabled:          true,
			EventReceiverIDs: []graphql.ID{"receiver-b", "receiver-c", "receiver-e"},
		})
	})

	t.Run("group doesn't exist", func(t *testing.T) {
		resp, err := client.Get(groupURI + "does-not-exist")
		assert.NilError(t, err)
		assert.Assert(t, isClientErrStatus(resp.StatusCode), "got status %d", resp.StatusCode)

		var body getGroupResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NilError(t, err)

		assert.Equal(t, len(body.Data), 0, "expect no receiver group")
		assert.Check(t, len(body.Errors) > 0, "expect resp body errors")
	})
}

func TestCreateGroupValid(t *testing.T) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	tests := map[string]struct {
		body io.Reader
	}{
		"enabled group": {
			body: strings.NewReader(`{
  "name": "group-x",
  "type": "security.scan",
  "version": "1.0.8",
  "description": "security scans of source code",
  "enabled": true,
  "event_receiver_ids": ["receiver-b", "receiver-c"]
}`),
		},
		"disabled group": {
			body: strings.NewReader(`{
  "name": "group-y",
  "type": "security.scan",
  "version": "5.4.3",
  "description": "security scans of source code",
  "enabled": false,
  "event_receiver_ids": ["receiver-b", "receiver-c"]
}`),
		},
		"single receiver": {
			body: strings.NewReader(`{
  "name": "group-z",
  "type": "security.scan",
  "version": "2.4.0",
  "description": "security scans of source code",
  "enabled": true,
  "event_receiver_ids": ["receiver-b"]
}`),
		},
		"multiple receivers": {
			body: strings.NewReader(`{
  "name": "group-v",
  "type": "lots.of.receivers",
  "version": "4.0.0",
  "description": "a whole lot of receivers",
  "enabled": true,
  "event_receiver_ids": ["receiver-a", "receiver-b", "receiver-c", "receiver-d"]
}`),
		},
		// TODO this should probably fail
		"no receivers": {
			body: strings.NewReader(`{
  "name": "group-q",
  "type": "absolutely.no.receivers",
  "version": "9.8.7",
  "description": "no receivers whatsoever",
  "enabled": true,
  "event_receiver_ids": []
}`),
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			resp, err := client.Post(groupURI, "application/json", tt.body)
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusOK)

			var body postGroupResponse
			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NilError(t, err)

			assert.Equal(t, len(body.Errors), 0, "expect no resp body errors")
			assert.Check(t, len(body.Data) > 0, "expect non-empty id")
		})
	}
}

func assertGroupsEqual(t *testing.T, got, expected storage.EventReceiverGroup) {
	t.Helper()
	assert.Equal(t, got.ID, expected.ID)
	assert.Equal(t, got.Name, expected.Name)
	assert.Equal(t, got.Type, expected.Type)
	assert.Equal(t, got.Version, expected.Version)
	assert.Equal(t, got.Description, expected.Description)
	assert.Equal(t, got.Enabled, expected.Enabled)

	gotCreatedTime := time.Time(got.CreatedAt.Date)
	expectedCreatedTime := time.Time(expected.CreatedAt.Date)
	// timestamp typically created upon db insertion & therefore hard to test.
	// If time not set in expected instance, ignore it
	if !expectedCreatedTime.IsZero() {
		assert.Check(t, gotCreatedTime.Equal(expectedCreatedTime), "got time %s, expected time %s", gotCreatedTime, expectedCreatedTime)
	}

	gotUpdatedTime := time.Time(got.UpdatedAt.Date)
	expectedUpdatedTime := time.Time(expected.UpdatedAt.Date)
	if !expectedUpdatedTime.IsZero() {
		assert.Check(t, gotUpdatedTime.Equal(expectedUpdatedTime), "got time %s, expected time %s", gotUpdatedTime, expectedUpdatedTime)
	}
}
