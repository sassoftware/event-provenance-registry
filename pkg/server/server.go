// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package server

import (
	_ "embed"
	"errors"
	"net/http"

	"github.com/graph-gophers/graphql-go/relay"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/graphql/schema"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

type contextKey string

var (
	contextKeySortBy  = contextKey("sortBy")
	contextKeySortDir = contextKey("sortDir")
	contextKeyStart   = contextKey("start")
	contextKeyLimit   = contextKey("limit")
)

type Server struct {
	DBConnector *storage.Database
}

func New(conn *storage.Database) (*Server, error) {
	if conn == nil {
		return nil, errors.New("database connector cannot be nil")
	}
	return &Server{
		DBConnector: conn,
	}, nil
}

// Paginate implements pagination for endpoints
func (s *Server) Paginate(_ http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: figure this out.
		// ctxt, err := setContext(r.Context(), r.URL.Query())
		// if err != nil {
		// 	logger.Error(err.Error())
		// 	render.Status(r, http.StatusBadRequest)
		// 	render.JSON(w, r, resp.NewRespError(err.Error()))
		// }
		// next.ServeHTTP(w, r.WithContext(ctxt))
	}

	return http.HandlerFunc(fn)
}

// func setContext(cont context.Context, query map[string][]string) (context.Context, error) {
// 	options := postgres.NewDefaultOptions()
// 	start := query["start"]
// 	if len(start) == 1 {
// 		strt, err := strconv.Atoi(start[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		options.Start = strt
// 	}
// 	ctxt := context.WithValue(cont, contextKeyStart, options.Start)
// 	limit := query["limit"]
// 	if len(limit) == 1 {
// 		lim, err := strconv.Atoi(limit[0])
// 		if err != nil {
// 			return nil, err
// 		}
// 		options.Limit = lim
// 	}
// 	ctxt = context.WithValue(ctxt, contextKeyLimit, options.Limit)
// 	return ctxt, nil
// }

// Sorting implements sorting for endpoints
func (s *Server) Sorting(_ http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: finish me.
		// ctxt := setSortingContext(r.Context(), r.URL.Query())
		// next.ServeHTTP(w, r.WithContext(ctxt))
	}

	return http.HandlerFunc(fn)
}

// func setSortingContext(cont context.Context, query map[string][]string) context.Context {
// 	options := postgres.NewDefaultOptions()
// 	sortAsc := query["sortAsc"]
// 	if len(sortAsc) == 1 {
// 		options.SortDir = "Ascending"
// 		options.SortBy = sortAsc[0]
// 	}
// 	sortDesc := query["sortDesc"]
// 	if len(sortDesc) == 1 {
// 		options.SortDir = "Descending"
// 		options.SortBy = sortDesc[0]
// 	}
// 	ctxt := context.WithValue(cont, contextKeySortBy, options.SortBy)
// 	ctxt = context.WithValue(ctxt, contextKeySortDir, options.SortDir)
// 	return ctxt
// }

func (s *Server) ServeOpenAPIDoc(_ string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement me
		panic("implement me!")
	}
}

//go:embed graphql.html
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
