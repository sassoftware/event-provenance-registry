package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

const groupURI = "http://localhost:8042/api/v1/groups/"

type getGroupResponse struct {
	Data   []storage.EventReceiverGroup
	Errors []string
}

type postGroupResponse struct {
	// ID of the created group
	Data   string
	Errors []string
}

type patchGroupResponse struct {
	// ID of the patched group
	Data   string
	Errors []string
}

type eventReceiverGroupInput struct {
	Name        string
	Type        string
	Version     string
	Description string
	Receivers   []string
}

func (g eventReceiverGroupInput) toPayload() string {
	receivers, _ := json.Marshal(g.Receivers)
	return fmt.Sprintf(`{
	"name": "%s",
	"type": "%s",
	"version": "%s",
	"description": "%s",
	"event_receiver_ids": %s
}`, g.Name, g.Type, g.Version, g.Description, string(receivers))
}

// createGroup creates a group with the given input, returning
// its ID or any errors that occurred. It is meant to simplify tests
// needing groups which don't care about details of group creation
func createGroup(client *http.Client, input eventReceiverGroupInput) (string, error) {
	resp, err := client.Post(groupURI, "application/json", strings.NewReader(input.toPayload()))
	if err != nil {
		return "", fmt.Errorf("failed to post group: %w", err)
	}

	var body postGroupResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", fmt.Errorf("failed to decode group resp body: %w", err)
	}

	if len(body.Errors) > 0 {
		return "", fmt.Errorf("group resp body had errors: %v", body.Errors)
	}

	return body.Data, nil
}

func getGroup(client *http.Client, id string) (storage.EventReceiverGroup, error) {
	resp, err := client.Get(groupURI + id)
	if err != nil {
		return storage.EventReceiverGroup{}, fmt.Errorf("failed to get group by id %s: %w", id, err)
	}

	var body getGroupResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return storage.EventReceiverGroup{}, fmt.Errorf("failed to decode group resp body: %w", err)
	}

	if len(body.Errors) > 0 {
		return storage.EventReceiverGroup{}, fmt.Errorf("group resp body had errors: %v", body.Errors)
	}
	if len(body.Data) > 1 {
		return storage.EventReceiverGroup{}, fmt.Errorf("found multiple groups by id %s", id)
	}

	return body.Data[0], nil
}
