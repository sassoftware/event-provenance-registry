//go:build functional

package functional

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"gotest.tools/v3/assert"
)

const receiverURI = "http://localhost:8042/api/v1/receivers/"

type getResponse struct {
	Data   []storage.EventReceiver
	Errors []string
}

type postResponse struct {
	Data   string
	Errors []string
}

func TestGetReceiverValid(t *testing.T) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	tests := map[string]struct {
		receiverId string
		expected   storage.EventReceiver
	}{
		"receiver exists empty schema": {
			receiverId: "receiver-d",
			expected: storage.EventReceiver{
				ID:          "receiver-d",
				Name:        "publish artifact",
				Type:        "artifact.publish",
				Version:     "2.0.0",
				Description: "publish to artifactory",
				Schema:      types.JSON{JSON: []byte("{}")},
				Fingerprint: "vci6fk9",
			},
		},
		"receiver exists with schema": {
			receiverId: "receiver-e",
			expected: storage.EventReceiver{
				ID:          "receiver-e",
				Name:        "manager sign-off",
				Type:        "signoff.complete",
				Version:     "1.0.0",
				Description: "manual sign-off of artifact",
				Schema:      types.JSON{JSON: []byte(`{"type":"object","properties":{"employee_id":{"type":"string"}}}`)},
				Fingerprint: "re5u2al",
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			resp, err := client.Get(receiverURI + tt.receiverId)
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusOK)

			var body getResponse
			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NilError(t, err)

			assert.Equal(t, len(body.Errors), 0, "expect no resp body errors")
			assert.Equal(t, len(body.Data), 1, "expect receiver in resp body")
			assertReceiversEqual(t, body.Data[0], tt.expected)
		})
	}

	t.Run("receiver doesn't exist", func(t *testing.T) {
		resp, err := client.Get(receiverURI + "does-not-exist")
		assert.NilError(t, err)
		assert.Assert(t, isClientErrStatus(resp.StatusCode), "got status %d", resp.StatusCode)

		var body getResponse
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NilError(t, err)

		assert.Equal(t, len(body.Data), 0, "expect no receivers in body")
		assert.Check(t, len(body.Errors) > 0, "expect resp body errors")
	})
}

func TestGetReceiverInvalid(t *testing.T) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	tests := map[string]struct {
		receiverId string
	}{
		"receiver doesn't exist": {
			receiverId: "does-not-exist",
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			resp, err := client.Get(receiverURI + tt.receiverId)
			assert.NilError(t, err)
			assert.Assert(t, isClientErrStatus(resp.StatusCode), "got status %d", resp.StatusCode)

			var body getResponse
			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NilError(t, err)

			assert.Equal(t, len(body.Data), 0, "expect no receivers in body")
			assert.Check(t, len(body.Errors) > 0, "expect resp body errors")
		})
	}
}

func TestCreateReceiverValid(t *testing.T) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	tests := map[string]struct {
		body io.Reader
	}{
		"empty schema": {
			body: strings.NewReader(`{
  "name": "receiver-z",
  "type": "artifact.deploy.eks",
  "version": "0.0.5",
  "description": "deploy to eks",
  "schema": {}
}`),
		},
		"basic schema": {
			body: strings.NewReader(`{
  "name": "receiver-x",
  "type": "artifact.deploy.gcp",
  "version": "0.3.2",
  "description": "deploy to gcp",
  "schema": {
    "type": "object",
    "properties": {
      "region": {
        "type": "string"
      }
    }
  }
}`),
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			resp, err := client.Post(receiverURI, "application/json", tt.body)
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusOK)

			var body postResponse
			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NilError(t, err)
			assert.Check(t, len(body.Data) > 0, "expect a non-empty id")
			assert.Equal(t, len(body.Errors), 0, "expect no errors")
		})
	}
}

func TestCreateReceiverInvalid(t *testing.T) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	tests := map[string]struct {
		body io.Reader
	}{
		"invalid schema": {
			body: strings.NewReader(`{
  "name": "receiver-y",
  "type": "artifact.deploy.openshift",
  "version": "1.0.1",
  "description": "deploy to openshift",
  "schema": {
    "type": "object",
    "properties": {
      "region": {
        "type"": "string"
      }
    }
  }
}`),
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			resp, err := client.Post(receiverURI, "application/json", tt.body)
			assert.NilError(t, err)
			assert.Assert(t, resp.StatusCode >= 400 && resp.StatusCode <= 499, "got status %d", resp.StatusCode)

			var body postResponse
			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NilError(t, err)
			assert.Equal(t, len(body.Data), 0, "expect an empty id, got %s", body.Data)
			assert.Check(t, len(body.Errors) > 0, "expect an error")
		})
	}
}

// assertReceiversEqual checks that two event receivers are functionally identical
func assertReceiversEqual(t *testing.T, got, expected storage.EventReceiver) {
	t.Helper()
	assert.Equal(t, got.ID, expected.ID)
	assert.Equal(t, got.Name, expected.Name)
	assert.Equal(t, got.Type, expected.Type)
	assert.Equal(t, got.Version, expected.Version)
	assert.Equal(t, got.Description, expected.Description)
	assert.Equal(t, got.Fingerprint, expected.Fingerprint)
	assert.Equal(t, got.Schema.String(), expected.Schema.String())

	gotTime := time.Time(got.CreatedAt.Date)
	expectedTime := time.Time(expected.CreatedAt.Date)
	// timestamp typically created upon db insertion & therefore hard to test.
	// If time not set in expected instance, ignore it
	if !expectedTime.IsZero() {
		assert.Check(t, gotTime.Equal(expectedTime), "got time %s, expected time %s", gotTime, expectedTime)
	}
}

func isClientErrStatus(status int) bool {
	return status >= 400 && status <= 499
}
