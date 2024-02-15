package epr

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

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

func (e EventInput) Validate() error {
	var err error

	if strings.TrimSpace(e.Name) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "name cannot be blank"})
	}
	if strings.TrimSpace(e.Version) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "version cannot be blank"})
	}
	if strings.TrimSpace(e.Release) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "release cannot be blank"})
	}
	if strings.TrimSpace(e.PlatformID) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "platform id cannot be blank"})
	}
	if strings.TrimSpace(e.Package) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "package cannot be blank"})
	}
	if strings.TrimSpace(e.Description) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "description cannot be blank"})
	}
	if strings.TrimSpace(string(e.EventReceiverID)) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "event receiver id cannot be blank"})
	}

	return err
}

type EventReceiverInput struct {
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Schema      types.JSON `json:"schema"`
}

func (r EventReceiverInput) Validate() error {
	var err error

	if strings.TrimSpace(r.Name) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "name cannot be blank"})
	}
	if strings.TrimSpace(r.Type) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "type cannot be blank"})
	}
	if strings.TrimSpace(r.Version) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "version cannot be blank"})
	}
	if strings.TrimSpace(r.Description) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "description cannot be blank"})
	}
	schemaErr := func() error {
		schema := r.Schema.String()
		if schema == "" {
			return eprErrors.InvalidInputError{Msg: "schema is required"}
		}
		loader := gojsonschema.NewStringLoader(schema)
		_, err := gojsonschema.NewSchema(loader)
		if err != nil {
			err = fmt.Errorf("failed to parse schema: %w", err)
			return eprErrors.InvalidInputError{Msg: err.Error()}
		}
		return nil
	}()
	err = errors.Join(err, schemaErr)

	return err
}

type EventReceiverGroupInput struct {
	Name             string       `json:"name"`
	Type             string       `json:"type"`
	Version          string       `json:"version"`
	Description      string       `json:"description"`
	Enabled          bool         `json:"enabled"`
	EventReceiverIDs []graphql.ID `json:"event_receiver_ids"`
}

func (g EventReceiverGroupInput) Validate() error {
	var err error

	if strings.TrimSpace(g.Name) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "name cannot be blank"})
	}
	if strings.TrimSpace(g.Type) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "type cannot be blank"})
	}
	if strings.TrimSpace(g.Version) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "version cannot be blank"})
	}
	if strings.TrimSpace(g.Description) == "" {
		err = errors.Join(err, eprErrors.InvalidInputError{Msg: "description cannot be blank"})
	}
	receiverIDErr := func() error {
		if len(g.EventReceiverIDs) == 0 {
			return eprErrors.InvalidInputError{Msg: "need at least one event receiver id"}
		}
		for _, receiverID := range g.EventReceiverIDs {
			if strings.TrimSpace(string(receiverID)) == "" {
				return eprErrors.InvalidInputError{Msg: "event receiver ids cannot be blank"}
			}
		}
		return nil
	}()
	err = errors.Join(err, receiverIDErr)

	return err
}

func CreateEvent(msgProducer message.TopicProducer, db *storage.Database, input EventInput) (*storage.Event, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

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
	err := input.Validate()
	if err != nil {
		return nil, err
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
	err := input.Validate()
	if err != nil {
		return nil, err
	}

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
