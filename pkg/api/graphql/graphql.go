package graphql

import (
	_ "embed"
	"net/http"

	"github.com/graph-gophers/graphql-go/relay"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/api/graphql/schema"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/config"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

type Server struct {
	DBConnector *storage.Database
	kafkaCfg    *config.KafkaConfig
}

func New(conn *storage.Database, cfg *config.KafkaConfig) *Server {
	return &Server{
		DBConnector: conn,
		kafkaCfg:    cfg,
	}
}

//go:embed resources/graphql.html
var graphqlHTML []byte

func (s *Server) ServerGraphQLDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(graphqlHTML)
	}
}

func (s *Server) GraphQLHandler(connection *storage.Database, cfg *config.KafkaConfig) http.HandlerFunc {
	handler := &relay.Handler{Schema: schema.New(connection, cfg)}
	return handler.ServeHTTP
}