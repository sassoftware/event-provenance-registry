package e2e

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sassoftware/event-provenance-registry/tests/common"
	"gotest.tools/v3/assert"
)

func TestCreateAndGetEvent(t *testing.T) {
	client := common.NewHTTPClient()

	tests := map[string]struct {
		input         eventInput
		receiverInput eventReceiverInput
	}{
		"successful event to empty schema receiver": {
			input: eventInput{
				Name:        "my-event",
				Version:     "1.0.0",
				Release:     "2024.01.22",
				PlatformID:  "amd64-oci-linux",
				Package:     "docker",
				Description: "my sample event",
				Payload:     "{}",
				Success:     true,
			},
			receiverInput: eventReceiverInput{
				Name:        "receiver without schema",
				Type:        "do.some.work",
				Version:     "3.2.4",
				Description: "Sample receiver for event testing",
				Schema:      "{}",
			},
		},
		"successful event to basic schema receiver": {
			input: eventInput{
				Name:        "commit qwerty123",
				Version:     "2.2.0",
				Release:     "2023.02.05",
				PlatformID:  "amd64-oci-linux",
				Package:     "docker",
				Description: "a commit being made",
				Payload:     `{"authorEmail":"cool.person@company.com"}`,
				Success:     true,
			},
			receiverInput: eventReceiverInput{
				Name:        "accept & merge PRs",
				Type:        "repo.commit.merged",
				Version:     "5.6.7",
				Description: "Sample receiver for event testing",
				Schema:      `{"type":"object","properties":{"authorEmail":{"type":"string"}}}`,
			},
		},
		"unsuccessful event to empty schema receiver": {
			input: eventInput{
				Name:        "my-other-event",
				Version:     "1.0.1",
				Release:     "2024.01.21",
				PlatformID:  "amd64-oci-linux",
				Package:     "docker",
				Description: "my other sample event",
				Payload:     "{}",
				Success:     false,
			},
			receiverInput: eventReceiverInput{
				Name:        "receiver without schema",
				Type:        "do.some.work",
				Version:     "3.2.4",
				Description: "Sample receiver for event testing",
				Schema:      "{}",
			},
		},
		"unsuccessful event to basic schema receiver": {
			input: eventInput{
				Name:        "upload foo-microservice",
				Version:     "1.0.0",
				Release:     "2022.12.03",
				PlatformID:  "arm-oci-linux",
				Package:     "docker",
				Description: "upload build of foo service",
				Payload:     `{"server":"my.artifactory.com","artifactSize":1029354}`,
				Success:     false,
			},
			receiverInput: eventReceiverInput{
				Name:        "publish artifacts to server",
				Type:        "artifactory.publish",
				Version:     "1.10.7",
				Description: "Sample receiver for event testing",
				Schema:      `{"type":"object","properties":{"server":{"type":"string"},"artifactSize":{"type":"number"}}}`,
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			receiverID, err := createReceiver(client, tt.receiverInput)
			assert.NilError(t, err)
			tt.input.EventReceiverID = receiverID

			resp, err := client.Post(eventURI, "application/json", strings.NewReader(tt.input.toPayload()))
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusOK)

			var postBody postEventResponse
			err = json.NewDecoder(resp.Body).Decode(&postBody)
			assert.NilError(t, err)
			assert.Equal(t, len(postBody.Errors), 0, "got error(s) %v", postBody.Errors)
			assert.Check(t, len(postBody.Data) > 0, "expect non-empty id")

			eventID := postBody.Data
			resp, err = client.Get(eventURI + eventID)
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusOK)

			var getBody getEventResponse
			err = json.NewDecoder(resp.Body).Decode(&getBody)
			assert.NilError(t, err)
			assert.Equal(t, len(getBody.Errors), 0, "got error(s) %v", getBody.Errors)
			assert.Equal(t, len(getBody.Data), 1)
			event := getBody.Data[0]
			assert.Equal(t, string(event.ID), eventID)
			assert.Equal(t, event.Name, tt.input.Name)
			assert.Equal(t, event.Version, tt.input.Version)
			assert.Equal(t, event.Release, tt.input.Release)
			assert.Equal(t, event.PlatformID, tt.input.PlatformID)
			assert.Equal(t, event.Package, tt.input.Package)
			assert.Equal(t, event.Description, tt.input.Description)
			assert.Equal(t, event.Payload.String(), tt.input.Payload)
			assert.Equal(t, event.Success, tt.input.Success)
			assert.Equal(t, string(event.EventReceiverID), tt.input.EventReceiverID)
			assert.Check(t, !time.Time(event.CreatedAt.Date).IsZero(), "expect time to be set")
		})
	}
}

func TestCreateEventWithInvalidInput(t *testing.T) {
	client := common.NewHTTPClient()

	receiverInput := eventReceiverInput{
		Name:        "merge PRs",
		Type:        "pr.merge",
		Version:     "1.0.26",
		Description: "accept & merge PRs on github",
		Schema:      `{"type":"object","properties":{"filesChanged":{"type":"number"},"reviewers":{"type":"array","items":{"type":"string"}}}}`,
	}
	receiverID, err := createReceiver(client, receiverInput)
	assert.NilError(t, err)

	tests := map[string]struct {
		input eventInput
	}{
		"incorrect payload field type": {
			input: eventInput{
				Name:            "merge PR for service foobar",
				Version:         "1.0.1",
				Release:         "2024.01.18",
				PlatformID:      "oci-linux",
				Package:         "docker",
				Description:     "merged a PR for foobar",
				Payload:         `{"filesChanged":true,"reviewers":["me"]}`,
				Success:         true,
				EventReceiverID: receiverID,
			},
		},
		"incorrect payload array value type": {
			input: eventInput{
				Name:            "merge PR for service foo",
				Version:         "2.11.1",
				Release:         "2017.04.12",
				PlatformID:      "oci-linux",
				Package:         "docker",
				Description:     "merged a PR for foo",
				Payload:         `{"filesChanged":16,"reviewers":["you", 4]}`,
				Success:         true,
				EventReceiverID: receiverID,
			},
		},
		"empty name": {
			input: eventInput{
				Name:            "",
				Version:         "9.9.9",
				Release:         "2006.01.01",
				PlatformID:      "oci-linux",
				Package:         "docker",
				Description:     "merged a PR for baz",
				Payload:         `{"filesChanged":16,"reviewers":["someone"]}`,
				Success:         true,
				EventReceiverID: receiverID,
			},
		},
		"empty version": {
			input: eventInput{
				Name:            "merge PR to baz repo",
				Version:         "",
				Release:         "2006.01.01",
				PlatformID:      "oci-linux",
				Package:         "docker",
				Description:     "merged a PR for baz",
				Payload:         `{"filesChanged":16,"reviewers":["someone"]}`,
				Success:         true,
				EventReceiverID: receiverID,
			},
		},
		"empty release": {
			input: eventInput{
				Name:            "merge PR to baz repo",
				Version:         "9.9.9",
				Release:         "",
				PlatformID:      "oci-linux",
				Package:         "docker",
				Description:     "merged a PR for baz",
				Payload:         `{"filesChanged":16,"reviewers":["someone"]}`,
				Success:         true,
				EventReceiverID: receiverID,
			},
		},
		"empty platform id": {
			input: eventInput{
				Name:            "merge PR to baz repo",
				Version:         "9.9.9",
				Release:         "2006.01.01",
				PlatformID:      "",
				Package:         "docker",
				Description:     "merged a PR for baz",
				Payload:         `{"filesChanged":16,"reviewers":["someone"]}`,
				Success:         true,
				EventReceiverID: receiverID,
			},
		},
		"empty package": {
			input: eventInput{
				Name:            "merge PR to baz repo",
				Version:         "9.9.9",
				Release:         "2006.01.01",
				PlatformID:      "oci-linux",
				Package:         "",
				Description:     "merged a PR for baz",
				Payload:         `{"filesChanged":16,"reviewers":["someone"]}`,
				Success:         true,
				EventReceiverID: receiverID,
			},
		},
		"empty description": {
			input: eventInput{
				Name:            "merge PR to baz repo",
				Version:         "9.9.9",
				Release:         "2006.01.01",
				PlatformID:      "oci-linux",
				Package:         "docker",
				Description:     "",
				Payload:         `{"filesChanged":16,"reviewers":["someone"]}`,
				Success:         true,
				EventReceiverID: receiverID,
			},
		},
		"empty payload": {
			input: eventInput{
				Name:            "merge PR to baz repo",
				Version:         "9.9.9",
				Release:         "2006.01.01",
				PlatformID:      "oci-linux",
				Package:         "docker",
				Description:     "merged a PR for baz",
				Payload:         "",
				Success:         true,
				EventReceiverID: receiverID,
			},
		},
		"empty event receiver": {
			input: eventInput{
				Name:            "merge PR to baz repo",
				Version:         "9.9.9",
				Release:         "2006.01.01",
				PlatformID:      "oci-linux",
				Package:         "docker",
				Description:     "merged a PR for baz",
				Payload:         `{"filesChanged":16,"reviewers":["someone"]}`,
				Success:         true,
				EventReceiverID: "",
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			resp, err := client.Post(eventURI, "application/json", strings.NewReader(tt.input.toPayload()))
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

			var body postEventResponse
			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NilError(t, err)
			assert.Check(t, len(body.Errors) > 0, "expect a resp error")
			assert.Equal(t, len(body.Data), 0, "shouldn't get an id if event creation fails")
		})
	}
}

func TestGetNonExistentEvent(t *testing.T) {
	client := common.NewHTTPClient()

	resp, err := client.Get(eventURI + "non-existent-event-id")
	assert.NilError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)

	var body getEventResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NilError(t, err)
	assert.Check(t, len(body.Errors) > 0)
	assert.Equal(t, len(body.Data), 0)
}
