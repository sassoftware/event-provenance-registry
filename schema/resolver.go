package schema

import "gitlab.sas.com/async-event-infrastructure/server/pkg/storage"

type Resolver struct{}

func (*Resolver) Query() *QueryResolver {
	return &QueryResolver{}
}

type QueryResolver struct{}

func (r *QueryResolver) Event(args struct{ ID string }) *storage.Event {
	return nil
}

func (r *QueryResolver) EventReceiver(args struct{ ID string }) *storage.EventReceiver {
	return nil
}

func (r *QueryResolver) EventReceiverGroup(args struct{ ID string }) *storage.EventReceiverGroup {
	return nil
}

func (*Resolver) Mutation() *MutationResolver {
	return &MutationResolver{}
}

type MutationResolver struct{}

