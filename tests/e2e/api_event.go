package e2e

import (
	"fmt"

	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

const eventURI = "http://localhost:8042/api/v1/events/"

type postEventResponse struct {
	// ID of the created event
	Data   string
	Errors []string
}

type getEventResponse struct {
	Data   []storage.Event
	Errors []string
}

type eventInput struct {
	Name            string
	Version         string
	Release         string
	PlatformID      string
	Package         string
	Description     string
	Payload         string
	Success         bool
	EventReceiverID string
}

func (e eventInput) toPayload() string {
	return fmt.Sprintf(`{
	"name": "%s",
	"version": "%s",
	"release": "%s",
	"platform_id": "%s",
	"package": "%s",
	"description": "%s",
	"payload": %s,
	"success": %t,
	"event_receiver_id": "%s"
}`, e.Name, e.Version, e.Release, e.PlatformID, e.Package, e.Description, e.Payload, e.Success, e.EventReceiverID)
}
