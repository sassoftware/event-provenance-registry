package resolvers

import (
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
)

var logger = utils.MustGetLogger("server", "pkg.api.graphql.resolvers")

type Resolver struct {
	Connection  *storage.Database
	msgProducer message.TopicProducer
}

func New(connection *storage.Database, msgProducer message.TopicProducer) *Resolver {
	return &Resolver{
		Connection:  connection,
		msgProducer: msgProducer,
	}
}

func (r *Resolver) Query() *QueryResolver {
	return &QueryResolver{
		Connection: r.Connection,
	}
}

func (r *Resolver) Mutation() *MutationResolver {
	return &MutationResolver{
		Connection:  r.Connection,
		msgProducer: r.msgProducer,
	}
}
