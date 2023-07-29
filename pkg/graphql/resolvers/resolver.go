package resolvers

import "gitlab.sas.com/async-event-infrastructure/server/pkg/storage"

type Resolver struct {
	Connection *storage.Database
}

func New(connection *storage.Database) *Resolver {
	return &Resolver{
		Connection: connection,
	}
}

func (r *Resolver) Query() *QueryResolver {
	return &QueryResolver{
		Connection: r.Connection,
	}
}

func (r *Resolver) Mutation() *MutationResolver {
	return &MutationResolver{
		Connection: r.Connection,
	}
}
