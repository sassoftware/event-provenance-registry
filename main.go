// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/graph"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

func main() {
	// cmd.Execute()

	port := "8080"

	connection, err := storage.New("localhost", "postgres", "", "", "postgres")
	if err != nil {
		log.Fatal("unable to connect to db", err)
	}

	err = connection.SyncSchema()
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Database: connection,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
