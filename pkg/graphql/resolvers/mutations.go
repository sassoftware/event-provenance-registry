package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/graphql/schema/types"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/utils"
)

var logger = utils.MustGetLogger("server", "pkg.graphql.resolvers.mutations")

type MutationResolver struct {
	Connection *storage.Database
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

	eventReciever, err := storage.CreateEvent(r.Connection.Client, eventInput)
	if err != nil {
		logger.Error(err, "error creating event", "input", eventInput)
		return nil, err
	}
	logger.Info("created", "eventReciever", eventReciever)
	return &eventReciever.ID, nil
}

func (r *MutationResolver) CreateEventReceiver(args struct{ EventReceiver EventReceiverInput }) (*graphql.ID, error) {
	// TODO: centralize this and make it look better
	eventRecieverInput := storage.EventReceiver{
		Name:        args.EventReceiver.Name,
		Type:        args.EventReceiver.Type,
		Version:     args.EventReceiver.Version,
		Description: args.EventReceiver.Description,
		Schema:      args.EventReceiver.Schema,
	}

	eventReciever, err := storage.CreateEventReceiver(r.Connection.Client, eventRecieverInput)
	if err != nil {
		logger.Error(err, "error creating event receiver", "input", eventRecieverInput)
		return nil, err
	}
	logger.Info("created", "eventReciever", eventReciever)
	return &eventReciever.ID, nil
}

func (r *MutationResolver) CreateEventReceiverGroup(args struct{ EventReceiverGroup EventReceiverGroupInput }) (*graphql.ID, error) {
	logger.Info("created", "eventReceiverGroup", args.EventReceiverGroup)
	return nil, nil
}

func (r *MutationResolver) SetEventReceiverGroupEnabled(args struct{ ID graphql.ID }) (*graphql.ID, error) {
	logger.Info("updated", "eventRecieverGroupEnabled", args.ID)
	return nil, nil
}

func (r *MutationResolver) SetEventReceiverGroupDisabled(args struct{ ID graphql.ID }) (*graphql.ID, error) {
	logger.Info("updated", "eventRecieverGroupDisabled", args.ID)
	return nil, nil
}
