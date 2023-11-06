package resolvers

import (
	"gitlab.sas.com/async-event-infrastructure/server/pkg/config"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/utils"
)

var logger = utils.MustGetLogger("server", "pkg.api.graphql.resolvers")

type Resolver struct {
	Connection *storage.Database
	kafkaCfg   *config.KafkaConfig
}

func New(connection *storage.Database, cfg *config.KafkaConfig) *Resolver {
	return &Resolver{
		Connection: connection,
		kafkaCfg:   cfg,
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
		kafkaCfg:   r.kafkaCfg,
	}
}
