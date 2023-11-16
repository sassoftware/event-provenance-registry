package resolvers

import (
	"github.com/sassoftware/event-provenance-registry/pkg/config"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
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
