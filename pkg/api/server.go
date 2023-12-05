// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"errors"

	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql"
	"github.com/sassoftware/event-provenance-registry/pkg/api/rest"
	"github.com/sassoftware/event-provenance-registry/pkg/config"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

// Server type contains a GraphQL server and a REST server.
// @property GraphQL - The GraphQL property is a server that handles GraphQL requests. GraphQL is a
// query language for APIs and a runtime for executing those queries with your existing data. It allows
// clients to request only the data they need, making it more efficient and flexible than traditional
// REST APIs.
// @property Rest - The `Rest` property is a server that handles RESTful API requests. It is
// responsible for handling HTTP requests and returning appropriate responses based on the requested
// resources and actions.
type Server struct {
	GraphQL *graphql.Server
	Rest    *rest.Server
}

// New function creates a new Server instance with a GraphQL and REST server, using the provided
// database connector, Kafka configuration, and wait group.
func New(ctx context.Context, conn *storage.Database, config *config.KafkaConfig) (*Server, error) {
	if conn == nil {
		return nil, errors.New("database connector cannot be nil")
	}
	return &Server{
		GraphQL: graphql.New(conn, config),
		Rest:    rest.New(ctx, conn, config),
	}, nil
}
