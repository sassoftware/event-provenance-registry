package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

type QueryResolver struct {
	Connection *storage.Database
}

func (r *QueryResolver) Event(args struct{ ID graphql.ID }) (*storage.Event, error) {
	return storage.FindEvent(r.Connection.Client, args.ID)
}

func (r *QueryResolver) EventReceiver(args struct{ ID graphql.ID }) (*storage.EventReceiver, error) {
	return storage.FindEventReceiver(r.Connection.Client, args.ID)
}

func (r *QueryResolver) EventReceiverGroup(args struct{ ID graphql.ID }) (*storage.EventReceiverGroup, error) {
	return storage.FindEventReceiverGroup(r.Connection.Client, args.ID)
}
