package resolvers

import (
	"log"

	"github.com/graph-gophers/graphql-go"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/graphql/schema/types"
)

type MutationResolver struct{}

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

func (r *MutationResolver) CreateEvent(args struct{ Event EventInput }) *graphql.ID {
	log.Print(args.Event)
	return nil
}

func (r *MutationResolver) CreateEventReceiver(args struct{ EventReceiver EventReceiverInput }) *graphql.ID {
	log.Print(args.EventReceiver)
	return nil
}

func (r *MutationResolver) CreateEventReceiverGroup(args struct{ EventReceiverGroup EventReceiverGroupInput }) *graphql.ID {
	log.Print(args.EventReceiverGroup)
	return nil
}

func (r *MutationResolver) SetEventReceiverGroupEnabled(args struct{ ID graphql.ID }) *graphql.ID {
	log.Print(args.ID)
	return nil
}

func (r *MutationResolver) SetEventReceiverGroupDisabled(args struct{ ID graphql.ID }) *graphql.ID {
	log.Print(args.ID)
	return nil
}
