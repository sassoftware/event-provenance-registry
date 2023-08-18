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

func (r *MutationResolver) CreateEvent(args struct{ Input EventInput }) (*graphql.ID, error) {
	// TODO: centralize this and make it look better
	eventInput := storage.Event{
		Name:            args.Input.Name,
		Version:         args.Input.Version,
		Release:         args.Input.Release,
		PlatformID:      args.Input.PlatformID,
		Package:         args.Input.Package,
		Description:     args.Input.Description,
		Payload:         args.Input.Payload,
		Success:         args.Input.Success,
		EventReceiverID: args.Input.EventReceiverID,
	}

	eventReciever, err := storage.CreateEvent(r.Connection.Client, eventInput)
	if err != nil {
		logger.Error(err, "error creating event", "input", eventInput)
		return nil, err
	}
	logger.Info("created", "eventReciever", eventReciever)
	return &eventReciever.ID, nil
}

func (r *MutationResolver) CreateEventReceiver(args struct{ Input EventReceiverInput }) (*graphql.ID, error) {
	// TODO: centralize this and make it look better
	eventRecieverInput := storage.EventReceiver{
		Name:        args.Input.Name,
		Type:        args.Input.Type,
		Version:     args.Input.Version,
		Description: args.Input.Description,
		Schema:      args.Input.Schema,
	}

	eventReciever, err := storage.CreateEventReceiver(r.Connection.Client, eventRecieverInput)
	if err != nil {
		logger.Error(err, "error creating event receiver", "input", eventRecieverInput)
		return nil, err
	}
	logger.Info("created", "eventReciever", eventReciever)
	return &eventReciever.ID, nil
}

func (r *MutationResolver) CreateEventReceiverGroup(args struct{ Input EventReceiverGroupInput }) (*graphql.ID, error) {
	// TODO: centralize this and make it look better
	eventRecieverGroupInput := storage.EventReceiverGroup{
		Name:             args.Input.Name,
		Type:             args.Input.Type,
		Version:          args.Input.Version,
		Description:      args.Input.Description,
		Enabled:          true,
		EventReceiverIDs: args.Input.EventReceiverIDs,
	}

	eventRecieverGroup, err := storage.CreateEventReceiverGroup(r.Connection.Client, eventRecieverGroupInput)
	if err != nil {
		logger.Error(err, "error creating event receiver group", "input", eventRecieverGroupInput)
		return nil, err
	}

	logger.Info("created", "eventReceiverGroup", eventRecieverGroup)
	return &eventRecieverGroup.ID, nil
}

func (r *MutationResolver) SetEventReceiverGroupEnabled(args struct{ ID graphql.ID }) (*graphql.ID, error) {
	err := storage.SetEventReceiverGroupEnabled(r.Connection.Client, args.ID, true)
	if err != nil {
		logger.Error(err, "error setting event receiver group enabled", "id", args.ID)
		return nil, err
	}
	logger.Info("updated", "eventRecieverGroupEnabled", args.ID)
	return &args.ID, nil
}

func (r *MutationResolver) SetEventReceiverGroupDisabled(args struct{ ID graphql.ID }) (*graphql.ID, error) {
	err := storage.SetEventReceiverGroupEnabled(r.Connection.Client, args.ID, false)
	if err != nil {
		logger.Error(err, "error setting event receiver group disabled", "id", args.ID)
		return nil, err
	}
	logger.Info("updated", "eventRecieverGroupDisabled", args.ID)
	return &args.ID, nil
}
