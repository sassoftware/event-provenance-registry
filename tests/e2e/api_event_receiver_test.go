package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/sassoftware/event-provenance-registry/tests/common"
	"gotest.tools/v3/assert"
)

const receiverURI = "http://localhost:8042/api/v1/receivers/"

type getReceiverResponse struct {
	Data   []storage.EventReceiver
	Errors []string
}

type postReceiverResponse struct {
	// ID of the created receiver
	Data   string
	Errors []string
}

type eventReceiverInput struct {
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"schema"`
}

func TestCreateAndGetReceiver(t *testing.T) {
	client := common.NewHTTPClient()

	tests := map[string]struct {
		input eventReceiverInput
	}{
		"empty schema": {
			input: eventReceiverInput{
				Name:        "has-no-schema",
				Type:        "artifact.deploy.eks",
				Version:     "0.0.5",
				Description: "deploy to eks",
				Schema:      []byte(`{}`),
			},
		},
		"basic schema": {
			input: eventReceiverInput{
				Name:        "has-basic-schema",
				Type:        "artifact.deploy.gcp",
				Version:     "0.3.2",
				Description: "deploy to gcp",
				Schema:      []byte(`{"type":"object","properties":{"region":{"type":"string"}}}`),
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			// TODO create fixture for resetting db after tests
			input, err := json.Marshal(tt.input)
			assert.NilError(t, err, "failed to marshal body for request")
			resp, err := client.Post(receiverURI, "application/json", bytes.NewBuffer(input))
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusOK)

			var postBody postReceiverResponse
			err = json.NewDecoder(resp.Body).Decode(&postBody)
			assert.NilError(t, err)
			assert.Equal(t, len(postBody.Errors), 0, "got error(s) %v", postBody.Errors)
			assert.Check(t, len(postBody.Data) > 0, "expect non-empty id")

			receiverID := postBody.Data
			resp, err = client.Get(receiverURI + receiverID)
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusOK)

			var getBody getReceiverResponse
			err = json.NewDecoder(resp.Body).Decode(&getBody)
			assert.NilError(t, err)
			assert.Equal(t, len(getBody.Errors), 0)
			assert.Equal(t, len(getBody.Data), 1)
			receiver := getBody.Data[0]
			assert.Equal(t, string(receiver.ID), receiverID)
			assert.Equal(t, receiver.Name, tt.input.Name)
			assert.Equal(t, receiver.Type, tt.input.Type)
			assert.Equal(t, receiver.Version, tt.input.Version)
			assert.Equal(t, receiver.Description, tt.input.Description)
			assert.Equal(t, receiver.Schema.String(), string(tt.input.Schema))
			assert.Check(t, len(receiver.Fingerprint) > 0, "expect non-empty fingerprint")
			assert.Check(t, !time.Time(receiver.CreatedAt.Date).IsZero(), "expect time to be set")
		})
	}
}
