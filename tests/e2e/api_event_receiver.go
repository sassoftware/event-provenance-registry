package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
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

func (r eventReceiverInput) toPayload() string {
	return fmt.Sprintf(`{
	"name": "%s",
	"type": "%s",
	"version": "%s",
	"description": "%s",
	"schema": %s
}`, r.Name, r.Type, r.Version, r.Description, r.Schema)
}

// createReceiver creates a receiver with the given input, returning
// its ID or any errors that occurred. It is meant to simplify tests
// needing receivers which don't care about details of receiver creation
func createReceiver(client *http.Client, input eventReceiverInput) (string, error) {
	resp, err := client.Post(receiverURI, "application/json", strings.NewReader(input.toPayload()))
	if err != nil {
		return "", fmt.Errorf("failed to post receiver: %w", err)
	}

	var body postReceiverResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", fmt.Errorf("failed to decode receiver resp body: %w", err)
	}

	if len(body.Errors) > 0 {
		return "", fmt.Errorf("receiver resp body had errors: %v", body.Errors)
	}

	return body.Data, nil
}
