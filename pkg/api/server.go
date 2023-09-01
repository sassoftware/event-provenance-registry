// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	_ "embed"
	"errors"

	"gitlab.sas.com/async-event-infrastructure/server/pkg/api/graphql"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/api/rest"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

type Server struct {
	GraphQL *graphql.Server
	Rest    *rest.Server
}

func New(conn *storage.Database) (*Server, error) {
	if conn == nil {
		return nil, errors.New("database connector cannot be nil")
	}
	return &Server{
		graphql.New(conn),
		rest.New(conn),
	}, nil
}
