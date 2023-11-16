package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema/types"
	"github.com/sassoftware/event-provenance-registry/pkg/config"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// The MutationResolver type is used to handle mutations in a GraphQL schema and has a connection to a
// database.
// @property Connection - The `Connection` property is a reference to a database connection object. It
// is used to establish a connection to a database and perform various database operations such as
// querying and modifying data.
type MutationResolver struct {
	Connection *storage.Database
	kafkaCfg   *config.KafkaConfig
}

type EventInput struct {
	Name            string
	Version         string
	Release         string
	PlatformID      string
	Package         string
	Description     string
	Payload         types.JSON
	Success         bool
	EventReceiverID graphql.ID
}

type EventReceiverInput struct {
	Name        string
	Type        string
	Version     string
	Description string
	Schema      types.JSON
}

type EventReceiverGroupInput struct {
	Name             string
	Type             string
	Version          string
	Description      string
	EventReceiverIDs []graphql.ID
}

func (r *MutationResolver) CreateEvent(args struct{ Event EventInput }) (*graphql.ID, error) {
	// TODO: centralize this and make it look better
	eventInput := storage.Event{
		Name:            args.Event.Name,
		Version:         args.Event.Version,
		Release:         args.Event.Release,
		PlatformID:      args.Event.PlatformID,
		Package:         args.Event.Package,
		Description:     args.Event.Description,
		Payload:         args.Event.Payload,
		Success:         args.Event.Success,
		EventReceiverID: args.Event.EventReceiverID,
	}

	event, err := storage.CreateEvent(r.Connection.Client, eventInput)
	if err != nil {
		logger.Error(err, "error creating event", "input", eventInput)
		return nil, err
	}

	r.kafkaCfg.MsgChannel <- message.NewEvent(event)

	logger.V(1).Info("created", "event", event)
	return &event.ID, nil
}

func (r *MutationResolver) CreateEventReceiver(args struct{ EventReceiver EventReceiverInput }) (*graphql.ID, error) {
	// TODO: centralize this and make it look better
	eventReceiverInput := storage.EventReceiver{
		Name:        args.EventReceiver.Name,
		Type:        args.EventReceiver.Type,
		Version:     args.EventReceiver.Version,
		Description: args.EventReceiver.Description,
		Schema:      args.EventReceiver.Schema,
	}

	eventReceiver, err := storage.CreateEventReceiver(r.Connection.Client, eventReceiverInput)
	if err != nil {
		logger.Error(err, "error creating event receiver", "input", eventReceiverInput)
		return nil, err
	}

	r.kafkaCfg.MsgChannel <- message.NewEventReceiver(eventReceiver)

	logger.V(1).Info("created", "eventReceiver", eventReceiver)
	return &eventReceiver.ID, nil
}

func (r *MutationResolver) CreateEventReceiverGroup(args struct{ EventReceiverGroup EventReceiverGroupInput }) (*graphql.ID, error) {
	// TODO: centralize this and make it look better
	eventReceiverGroupInput := storage.EventReceiverGroup{
		Name:             args.EventReceiverGroup.Name,
		Type:             args.EventReceiverGroup.Type,
		Version:          args.EventReceiverGroup.Version,
		Description:      args.EventReceiverGroup.Description,
		Enabled:          true,
		EventReceiverIDs: args.EventReceiverGroup.EventReceiverIDs,
	}

	eventReceiverGroup, err := storage.CreateEventReceiverGroup(r.Connection.Client, eventReceiverGroupInput)
	if err != nil {
		logger.Error(err, "error creating event receiver group", "input", eventReceiverGroupInput)
		return nil, err
	}

	r.kafkaCfg.MsgChannel <- message.NewEventReceiverGroup(eventReceiverGroup)

	logger.V(1).Info("created", "eventReceiverGroup", eventReceiverGroup)
	return &eventReceiverGroup.ID, nil
}

func (r *MutationResolver) SetEventReceiverGroupEnabled(args struct{ ID graphql.ID }) (*graphql.ID, error) {
	err := storage.SetEventReceiverGroupEnabled(r.Connection.Client, args.ID, true)
	if err != nil {
		logger.Error(err, "error setting event receiver group enabled", "id", args.ID)
		return nil, err
	}
	logger.V(1).Info("updated", "eventReceiverGroupEnabled", args.ID)
	return &args.ID, nil
}

func (r *MutationResolver) SetEventReceiverGroupDisabled(args struct{ ID graphql.ID }) (*graphql.ID, error) {
	err := storage.SetEventReceiverGroupEnabled(r.Connection.Client, args.ID, false)
	if err != nil {
		logger.Error(err, "error setting event receiver group disabled", "id", args.ID)
		return nil, err
	}
	logger.V(1).Info("updated", "eventReceiverGroupDisabled", args.ID)
	return &args.ID, nil
}
