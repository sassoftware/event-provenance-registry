package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

type QueryResolver struct {
	Connection *storage.Database
}

func (r *QueryResolver) Events(args struct{ Event FindEventInput }) ([]storage.Event, error) {
	eventInput := storage.Event{
		ID:              *args.Event.ID,
		Name:            *args.Event.Name,
		Version:         *args.Event.Version,
		Release:         *args.Event.Release,
		PlatformID:      *args.Event.PlatformID,
		Package:         *args.Event.Package,
		Success:         *args.Event.Success,
		EventReceiverID: *args.Event.EventReceiverID,
	}
	return storage.FindEvent(r.Connection.Client, eventInput)
}

func (r *QueryResolver) EventReceivers(args struct{ EventReceiver FindEventReceiverInput }) ([]storage.EventReceiver, error) {
	eventReceiverInput := storage.EventReceiver{
		ID:      *args.EventReceiver.ID,
		Name:    *args.EventReceiver.Name,
		Type:    *args.EventReceiver.Type,
		Version: *args.EventReceiver.Version,
	}
	return storage.FindEventReceiver(r.Connection.Client, eventReceiverInput)
}

func (r *QueryResolver) EventReceiverGroups(args struct{ EventReceiverGroup FindEventReceiverGroupInput }) ([]storage.EventReceiverGroup, error) {
	eventReceiverGroupInput := storage.EventReceiverGroup{
		ID:      *args.EventReceiverGroup.ID,
		Name:    *args.EventReceiverGroup.Name,
		Type:    *args.EventReceiverGroup.Type,
		Version: *args.EventReceiverGroup.Version,
	}
	return storage.FindEventReceiverGroup(r.Connection.Client, eventReceiverGroupInput)
}

func (r *QueryResolver) EventsByID(args struct{ ID graphql.ID }) ([]storage.Event, error) {
	return storage.FindEventByID(r.Connection.Client, args.ID)
}

func (r *QueryResolver) EventReceiversByID(args struct{ ID graphql.ID }) ([]storage.EventReceiver, error) {
	return storage.FindEventReceiverByID(r.Connection.Client, args.ID)
}

func (r *QueryResolver) EventReceiverGroupsByID(args struct{ ID graphql.ID }) ([]storage.EventReceiverGroup, error) {
	return storage.FindEventReceiverGroupByID(r.Connection.Client, args.ID)
}
