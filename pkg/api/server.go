// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	_ "embed"
	"errors"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/config"
	"sync"

	"gitlab.sas.com/async-event-infrastructure/server/pkg/api/graphql"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/api/rest"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

type Server struct {
	GraphQL *graphql.Server
	Rest    *rest.Server
}

func New(ctx context.Context, conn *storage.Database, config *config.KafkaConfig, wg *sync.WaitGroup) (*Server, error) {
	if conn == nil {
		return nil, errors.New("database connector cannot be nil")
	}
	return &Server{
		GraphQL: graphql.New(conn),
		Rest:    rest.New(ctx, conn, config, wg),
	}, nil
}
