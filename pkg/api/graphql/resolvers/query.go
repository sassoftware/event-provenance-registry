package resolvers

import (
	"github.com/graph-gophers/graphql-go"
	eprErrors "github.com/sassoftware/event-provenance-registry/pkg/errors"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

type QueryResolver struct {
	Connection *storage.Database
}

func (r *QueryResolver) Events(args struct{ Event FindEventInput }) ([]storage.Event, error) {
	events, err := storage.FindEvent(r.Connection.Client, args.Event.toMap())
	return events, eprErrors.SanitizeError(err)
}

func (r *QueryResolver) EventReceivers(args struct{ EventReceiver FindEventReceiverInput }) ([]storage.EventReceiver, error) {
	receivers, err := storage.FindEventReceiver(r.Connection.Client, args.EventReceiver.toMap())
	return receivers, eprErrors.SanitizeError(err)
}

func (r *QueryResolver) EventReceiverGroups(args struct{ EventReceiverGroup FindEventReceiverGroupInput }) ([]storage.EventReceiverGroup, error) {
	groups, err := storage.FindEventReceiverGroup(r.Connection.Client, args.EventReceiverGroup.toMap())
	return groups, eprErrors.SanitizeError(err)
}

func (r *QueryResolver) EventsByID(args struct{ ID graphql.ID }) ([]storage.Event, error) {
	events, err := storage.FindEventByID(r.Connection.Client, args.ID)
	return events, eprErrors.SanitizeError(err)
}

func (r *QueryResolver) EventReceiversByID(args struct{ ID graphql.ID }) ([]storage.EventReceiver, error) {
	receivers, err := storage.FindEventReceiverByID(r.Connection.Client, args.ID)
	return receivers, eprErrors.SanitizeError(err)
}

func (r *QueryResolver) EventReceiverGroupsByID(args struct{ ID graphql.ID }) ([]storage.EventReceiverGroup, error) {
	groups, err := storage.FindEventReceiverGroupByID(r.Connection.Client, args.ID)
	return groups, eprErrors.SanitizeError(err)
}
