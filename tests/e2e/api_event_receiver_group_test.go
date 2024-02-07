package e2e

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/tests/common"
	"gotest.tools/v3/assert"
)

func TestCreateAndGetGroup(t *testing.T) {
	client := common.NewHTTPClient()

	tests := map[string]struct {
		input          eventReceiverGroupInput
		receiverInputs []eventReceiverInput
	}{
		"group with one receiver": {
			input: eventReceiverGroupInput{
				Name:        "one-receiver",
				Type:        "one.receiver",
				Version:     "4.5.6",
				Description: "has a single receiver",
			},
			receiverInputs: []eventReceiverInput{
				{
					Name:        "simple receiver",
					Type:        "does.something",
					Version:     "5.0.0",
					Description: "tracks that something happens",
					Schema:      `{}`,
				},
			},
		},
		"group with multiple receivers": {
			input: eventReceiverGroupInput{
				Name:        "multiple-receivers",
				Type:        "multiple.receivers",
				Version:     "7.8.9",
				Description: "has multiple receivers",
			},
			receiverInputs: []eventReceiverInput{
				{
					Name:        "receiver A",
					Type:        "does.something",
					Version:     "16.0.7",
					Description: "tracks that something happens",
					Schema:      `{}`,
				},
				{
					Name:        "receiver B",
					Type:        "does.something.else",
					Version:     "4.4.4",
					Description: "tracks that something else happens",
					Schema:      `{"type":"object","properties":{"someProp":{"type":"string"}}}`,
				},
				{
					Name:        "receiver C",
					Type:        "does.something.different",
					Version:     "1.21.4",
					Description: "tracks that something different happens",
					Schema:      `{}`,
				},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			var receiverIDs []string
			for _, receiverInput := range tt.receiverInputs {
				receiverID, err := createReceiver(client, receiverInput)
				assert.NilError(t, err)
				receiverIDs = append(receiverIDs, receiverID)
			}
			tt.input.Receivers = receiverIDs

			resp, err := client.Post(groupURI, "application/json", strings.NewReader(tt.input.toPayload()))
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusOK)

			var postBody postGroupResponse
			err = json.NewDecoder(resp.Body).Decode(&postBody)
			assert.NilError(t, err)
			assert.Equal(t, len(postBody.Errors), 0, "get error(s) %v", postBody.Errors)
			assert.Check(t, len(postBody.Data) > 0, "expect non-empty id")

			groupID := postBody.Data
			resp, err = client.Get(groupURI + groupID)
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusOK)

			var getBody getGroupResponse
			err = json.NewDecoder(resp.Body).Decode(&getBody)
			assert.NilError(t, err)
			assert.Equal(t, len(getBody.Errors), 0, "got error(s) %v", getBody.Errors)
			assert.Equal(t, len(getBody.Data), 1)
			group := getBody.Data[0]
			assert.Equal(t, string(group.ID), groupID)
			assert.Equal(t, group.Name, tt.input.Name)
			assert.Equal(t, group.Type, tt.input.Type)
			assert.Equal(t, group.Version, tt.input.Version)
			assert.Equal(t, group.Description, tt.input.Description)
			assert.Equal(t, group.Enabled, true)
			assert.DeepEqual(t, graphIDsToStrings(group.EventReceiverIDs), tt.input.Receivers)
			createdAt := time.Time(group.CreatedAt.Date)
			updatedAt := time.Time(group.UpdatedAt.Date)
			assert.Check(t, !createdAt.IsZero(), "expect time to be set")
			assert.Check(t, !updatedAt.IsZero(), "expect time to be set")
			assert.Check(t, createdAt.Equal(updatedAt), "on init times should be the same")
		})
	}
}

func TestCreateInvalidGroup(t *testing.T) {
	client := common.NewHTTPClient()

	// common valid receiver used for mixing in with invalid inputs
	receiverID, err := createReceiver(client, eventReceiverInput{
		Name:        "common-group-receiver",
		Type:        "my.receiver.type",
		Version:     "0.0.3",
		Description: "common receiver for invalid group testing",
		Schema:      `{}`,
	})
	assert.NilError(t, err)

	tests := map[string]struct {
		input eventReceiverGroupInput
	}{
		"empty receiver id": {
			input: eventReceiverGroupInput{
				Name:        "my group",
				Type:        "my.group",
				Version:     "1.9.7",
				Description: "my group",
				Receivers:   []string{"", receiverID},
			},
		},
		"non-existent receiver": {
			input: eventReceiverGroupInput{
				Name:        "your group",
				Type:        "your.group",
				Version:     "5.3.2",
				Description: "your group",
				Receivers:   []string{receiverID, "should-not-exist"},
			},
		},
		"group with no receivers": {
			input: eventReceiverGroupInput{
				Name:        "no-receivers",
				Type:        "no.receivers",
				Version:     "1.2.3",
				Description: "doesn't have any receivers",
				Receivers:   []string{},
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			resp, err := client.Post(groupURI, "application/json", strings.NewReader(tt.input.toPayload()))
			assert.NilError(t, err)
			assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

			var body postGroupResponse
			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NilError(t, err)
			assert.Check(t, len(body.Errors) > 0, "expect an error for invalid input")
			assert.Equal(t, len(body.Data), 0, "expect no id generated")
		})
	}
}

func TestToggleGroup(t *testing.T) {
	client := common.NewHTTPClient()

	groupID, err := createGroup(client, eventReceiverGroupInput{
		Name:        "group to be toggled",
		Type:        "toggled.group",
		Version:     "2.0.11",
		Description: "group that will be toggled",
		Receivers:   []string{},
	})
	assert.NilError(t, err)

	originalGroup, err := getGroup(client, groupID)
	assert.NilError(t, err)
	originalUpdatedAt := time.Time(originalGroup.UpdatedAt.Date)
	assert.Equal(t, originalGroup.Enabled, true, "group should be enabled by default")

	err = toggleGroup(client, groupID, false)
	assert.NilError(t, err)
	group, err := getGroup(client, groupID)
	assert.NilError(t, err)
	assert.Equal(t, group.Enabled, false, "group should be disabled")
	afterDisableUpdatedAt := time.Time(group.UpdatedAt.Date)
	assert.Check(t, !afterDisableUpdatedAt.Equal(originalUpdatedAt), "updated time should be updated after disable")

	err = toggleGroup(client, groupID, true)
	assert.NilError(t, err)
	group, err = getGroup(client, groupID)
	assert.NilError(t, err)
	assert.Equal(t, group.Enabled, true, "group should be enabled")
	afterEnableUpdatedAt := time.Time(group.UpdatedAt.Date)
	assert.Check(t, !afterEnableUpdatedAt.Equal(afterDisableUpdatedAt), "updated time should be updated after enable")
}

func graphIDsToStrings(ids []graphql.ID) []string {
	var s []string
	for _, id := range ids {
		s = append(s, string(id))
	}
	return s
}
