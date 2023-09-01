package graphql

import (
	_ "embed"
	"net/http"

	"github.com/graph-gophers/graphql-go/relay"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/api/graphql/schema"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

type Server struct {
	DBConnector *storage.Database
}

func New(conn *storage.Database) *Server {
	return &Server{
		DBConnector: conn,
	}
}

//go:embed resources/graphql.html
var graphqlHTML []byte

func (s *Server) ServerGraphQLDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(graphqlHTML)
	}
}

func (s *Server) GraphQLHandler(connection *storage.Database) http.HandlerFunc {
	handler := &relay.Handler{Schema: schema.New(connection)}
	return handler.ServeHTTP
}
