package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

type QueryResolver struct {
	Connection *storage.Database
}

func (r *QueryResolver) Events(args struct{ ID graphql.ID }) (*[]*storage.Event, error) {
	return storage.FindEvent(r.Connection.Client, args.ID)
}

func (r *QueryResolver) EventReceivers(args struct{ ID graphql.ID }) (*[]*storage.EventReceiver, error) {
	return storage.FindEventReceiver(r.Connection.Client, args.ID)
}

func (r *QueryResolver) EventReceiverGroups(args struct{ ID graphql.ID }) (*[]*storage.EventReceiverGroup, error) {
	return storage.FindEventReceiverGroup(r.Connection.Client, args.ID)
}
