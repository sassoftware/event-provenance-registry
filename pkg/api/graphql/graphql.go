package graphql

import (
	_ "embed"
	"net/http"

	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

type Server struct {
	DBConnector *storage.Database
	msgProducer message.TopicProducer
}

func New(conn *storage.Database, msgProducer message.TopicProducer) *Server {
	return &Server{
		DBConnector: conn,
		msgProducer: msgProducer,
	}
}

//go:embed resources/graphql.html
var graphqlHTML []byte

func (s *Server) ServerGraphQLDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(graphqlHTML)
	}
}

func (s *Server) GraphQLHandler() http.HandlerFunc {
	handler := &relay.Handler{Schema: schema.New(s.DBConnector, s.msgProducer)}
	return handler.ServeHTTP
}
