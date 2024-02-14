package resolvers

import (
	"log/slog"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/epr"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// The MutationResolver type is used to handle mutations in a GraphQL schema and has a connection to a
// database.
// @property Connection - The `Connection` property is a reference to a database connection object. It
// is used to establish a connection to a database and perform various database operations such as
// querying and modifying data.
type MutationResolver struct {
	Connection  *storage.Database
	msgProducer message.TopicProducer
}

func (r *MutationResolver) CreateEvent(args struct{ Event epr.EventInput }) (graphql.ID, error) {
	event, err := epr.CreateEvent(r.msgProducer, r.Connection, args.Event)
	if err != nil {
		return "", err
	}
	return event.ID, nil
}

func (r *MutationResolver) CreateEventReceiver(args struct{ EventReceiver epr.EventReceiverInput }) (graphql.ID, error) {
	eventReceiver, err := epr.CreateEventReceiver(r.msgProducer, r.Connection, args.EventReceiver)
	if err != nil {
		return "", err
	}
	return eventReceiver.ID, nil
}

func (r *MutationResolver) CreateEventReceiverGroup(args struct{ EventReceiverGroup epr.EventReceiverGroupInput }) (graphql.ID, error) {
	eventReceiverGroup, err := epr.CreateEventReceiverGroup(r.msgProducer, r.Connection, args.EventReceiverGroup)
	if err != nil {
		return "", err
	}
	return eventReceiverGroup.ID, nil
}

func (r *MutationResolver) SetEventReceiverGroupEnabled(args struct{ ID graphql.ID }) (graphql.ID, error) {
	err := storage.SetEventReceiverGroupEnabled(r.Connection.Client, args.ID, true)
	if err != nil {
		slog.Error("error setting event receiver group enabled", "error", err, "id", args.ID)
		return "", err
	}
	slog.Info("updated", "eventReceiverGroupEnabled", args.ID)
	return args.ID, nil
}

func (r *MutationResolver) SetEventReceiverGroupDisabled(args struct{ ID graphql.ID }) (graphql.ID, error) {
	err := storage.SetEventReceiverGroupEnabled(r.Connection.Client, args.ID, false)
	if err != nil {
		slog.Error("error setting event receiver group disabled", "error", err, "id", args.ID)
		return "", err
	}
	slog.Info("updated", "eventReceiverGroupDisabled", args.ID)
	return args.ID, nil
}
