package resolvers

import (
	"log"

	"github.com/graph-gophers/graphql-go"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

type QueryResolver struct {
	Connection *storage.Database
}

func (r *QueryResolver) Event(args struct{ ID graphql.ID }) (*storage.Event, error) {
	event, err := storage.FindEvent(r.Connection.Client, args.ID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *QueryResolver) EventReceiver(args struct{ ID graphql.ID }) (*storage.EventReceiver, error) {
	log.Print(args.ID)
	return nil, nil
}

func (r *QueryResolver) EventReceiverGroup(args struct{ ID graphql.ID }) (*storage.EventReceiverGroup, error) {
	log.Print(args.ID)
	return nil, nil
}
