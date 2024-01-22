package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Name        string
	Type        string
	Version     string
	Description string
	Schema      string
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
				Schema:      `{}`,
			},
		},
		"basic schema": {
			input: eventReceiverInput{
				Name:        "has-basic-schema",
				Type:        "artifact.deploy.gcp",
				Version:     "0.3.2",
				Description: "deploy to gcp",
				Schema:      `{"type":"object","properties":{"region":{"type":"string"}}}`,
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			resp, err := client.Post(receiverURI, "application/json", strings.NewReader(tt.input.toPayload()))
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
			assert.Equal(t, receiver.Schema.String(), tt.input.Schema)
			assert.Check(t, len(receiver.Fingerprint) > 0, "expect non-empty fingerprint")
			assert.Check(t, !time.Time(receiver.CreatedAt.Date).IsZero(), "expect time to be set")
		})
	}
}

func TestCreateInvalidReceiver(t *testing.T) {
	client := common.NewHTTPClient()

	receiver := eventReceiverInput{
		Name:        "upload artifact",
		Type:        "artifact.publish",
		Version:     "1.0.9",
		Description: "upload an artifact somewhere",
		Schema:      `{abcd}`,
	}

	resp, err := client.Post(receiverURI, "application/json", strings.NewReader(receiver.toPayload()))
	assert.NilError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	var body postReceiverResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NilError(t, err)
	assert.Check(t, len(body.Errors) > 0)
	assert.Equal(t, len(body.Data), 0, "shouldn't get id if creating receiver fails")
}

func TestGetNonExistentReceiver(t *testing.T) {
	client := common.NewHTTPClient()

	resp, err := client.Get(receiverURI + "non-existent-receiver-id")
	assert.NilError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)

	var body getReceiverResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NilError(t, err)
	assert.Check(t, len(body.Errors) > 0)
	assert.Equal(t, len(body.Data), 0)
}

func (r *eventReceiverInput) toPayload() string {
	return fmt.Sprintf(`{
	"name": "%s",
	"type": "%s",
	"version": "%s",
	"description": "%s",
	"schema": %s
}`, r.Name, r.Type, r.Version, r.Description, r.Schema)
}
