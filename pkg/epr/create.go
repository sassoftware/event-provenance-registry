package epr

import (
	"fmt"
	"log/slog"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	eprErrors "github.com/sassoftware/event-provenance-registry/pkg/errors"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/xeipuuv/gojsonschema"
)

type EventInput struct {
	Name            string     `json:"name"`
	Version         string     `json:"version"`
	Release         string     `json:"release"`
	PlatformID      string     `json:"platform_id"`
	Package         string     `json:"package"`
	Description     string     `json:"description"`
	Payload         types.JSON `json:"payload"`
	Success         bool       `json:"success"`
	EventReceiverID graphql.ID `json:"event_receiver_id"`
}

type EventReceiverInput struct {
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Schema      types.JSON `json:"schema"`
}

type EventReceiverGroupInput struct {
	Name             string       `json:"name"`
	Type             string       `json:"type"`
	Version          string       `json:"version"`
	Description      string       `json:"description"`
	Enabled          bool         `json:"enabled"`
	EventReceiverIDs []graphql.ID `json:"event_receiver_ids"`
}

func CreateEvent(msgProducer message.TopicProducer, db *storage.Database, input EventInput) (*storage.Event, error) {
	partial := storage.Event{
		Name:            input.Name,
		Version:         input.Version,
		Release:         input.Release,
		PlatformID:      input.PlatformID,
		Package:         input.Package,
		Description:     input.Description,
		Payload:         input.Payload,
		Success:         input.Success,
		EventReceiverID: input.EventReceiverID,
	}
	event, err := storage.CreateEvent(db.Client, partial)
	if err != nil {
		slog.Error("error creating event", "error", err, "input", input)
		return nil, err
	}

	msgProducer.Async(message.NewEvent(*event))
	slog.Info("created", "event", event)

	eventReceiverGroups, err := storage.FindTriggeredEventReceiverGroups(db.Client, *event)
	if err != nil {
		slog.Error("error finding triggered event receiver groups", "error", err, "input", input)
		return nil, err
	}

	for _, eventReceiverGroup := range eventReceiverGroups {
		msgProducer.Async(message.NewEventReceiverGroupComplete(*event, eventReceiverGroup))
	}

	return event, nil
}

func CreateEventReceiver(msgProducer message.TopicProducer, db *storage.Database, input EventReceiverInput) (*storage.EventReceiver, error) {
	schema := input.Schema.String()
	if schema == "" {
		return nil, eprErrors.InvalidInputError{Msg: "schema is required"}
	}
	loader := gojsonschema.NewStringLoader(schema)
	_, err := gojsonschema.NewSchema(loader)
	if err != nil {
		err = fmt.Errorf("failed to parse schema: %w", err)
		return nil, eprErrors.InvalidInputError{Msg: err.Error()}
	}

	partial := storage.EventReceiver{
		Name:        input.Name,
		Type:        input.Type,
		Version:     input.Version,
		Description: input.Description,
		Schema:      input.Schema,
	}

	receiver, err := storage.CreateEventReceiver(db.Client, partial)
	if err != nil {
		slog.Error("error creating event receiver", "error", err, "input", input)
		return nil, err
	}

	msgProducer.Async(message.NewEventReceiver(*receiver))
	slog.Info("created", "eventReceiver", receiver)

	return receiver, nil
}

func CreateEventReceiverGroup(msgProducer message.TopicProducer, db *storage.Database, input EventReceiverGroupInput) (*storage.EventReceiverGroup, error) {
	partial := storage.EventReceiverGroup{
		Name:             input.Name,
		Type:             input.Type,
		Version:          input.Version,
		Description:      input.Description,
		Enabled:          input.Enabled,
		EventReceiverIDs: input.EventReceiverIDs,
	}

	group, err := storage.CreateEventReceiverGroup(db.Client, partial)
	if err != nil {
		slog.Error("error creating event receiver group", "error", err, "input", input)
		return nil, err
	}

	msgProducer.Async(message.NewEventReceiverGroupCreated(*group))
	slog.Info("created", "eventReceiverGroup", group)

	return group, nil
}
