package resolvers

import (
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

type QueryResolver struct {
	Connection *storage.Database
}

func (r *QueryResolver) Events(args struct{ Event FindEventInput }) ([]storage.Event, error) {
	return storage.FindEvent(r.Connection.Client, args.Event.toMap())
}

func (r *QueryResolver) EventReceivers(args struct{ EventReceiver FindEventReceiverInput }) ([]storage.EventReceiver, error) {
	return storage.FindEventReceiver(r.Connection.Client, args.EventReceiver.toMap())
}

func (r *QueryResolver) EventReceiverGroups(args struct{ EventReceiverGroup FindEventReceiverGroupInput }) ([]storage.EventReceiverGroup, error) {
	return storage.FindEventReceiverGroup(r.Connection.Client, args.EventReceiverGroup.toMap())
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
