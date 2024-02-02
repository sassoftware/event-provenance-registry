package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

type QueryResolver struct {
	Connection *storage.Database
}

func (r *QueryResolver) Events(args struct{ Event storage.Event }) ([]storage.Event, error) {
	return storage.FindEvent(r.Connection.Client, args.Event)
}

func (r *QueryResolver) EventReceivers(args struct{ EventReceiver storage.EventReceiver }) ([]storage.EventReceiver, error) {
	return storage.FindEventReceiver(r.Connection.Client, args.EventReceiver)
}

func (r *QueryResolver) EventReceiverGroups(args struct{ EventReceiverGroup storage.EventReceiverGroup }) ([]storage.EventReceiverGroup, error) {
	return storage.FindEventReceiverGroup(r.Connection.Client, args.EventReceiverGroup)
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
