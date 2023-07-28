package resolvers

import (
	"log"

	"github.com/graph-gophers/graphql-go"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

type QueryResolver struct{}

func (r *QueryResolver) Event(args struct{ ID graphql.ID }) *storage.Event {
	log.Print(args.ID)
	return nil
}

func (r *QueryResolver) EventReceiver(args struct{ ID graphql.ID }) *storage.EventReceiver {
	log.Print(args.ID)
	return nil
}

func (r *QueryResolver) EventReceiverGroup(args struct{ ID graphql.ID }) *storage.EventReceiverGroup {
	log.Print(args.ID)
	return nil
}
