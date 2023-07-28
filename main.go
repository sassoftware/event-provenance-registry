// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/graph-gophers/graphql-go/relay"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/graphql/schema"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

func main() {
	// set up graphql schema
	schema := schema.New()

	// cmd.Execute()

	connection, err := storage.New("localhost", "postgres", "", "", "postgres")
	if err != nil {
		log.Fatal("unable to connect to db", err)
	}

	err = connection.SyncSchema()
	if err != nil {
		log.Fatal("failed to sync schema", err)
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(page)
	}))

	http.Handle("/query", &relay.Handler{Schema: schema})

	port := ":8080"
	log.Printf("connect to http://localhost%s/ for GraphQL playground", port)

	server := &http.Server{
		Addr:              port,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

var page = []byte(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>GraphiQL</title>
    <style>
      body {
        height: 100%;
        margin: 0;
        width: 100%;
        overflow: hidden;
      }
      #graphiql {
        height: 100vh;
      }
    </style>
    <script src="https://unpkg.com/react@17/umd/react.development.js" integrity="sha512-Vf2xGDzpqUOEIKO+X2rgTLWPY+65++WPwCHkX2nFMu9IcstumPsf/uKKRd5prX3wOu8Q0GBylRpsDB26R6ExOg==" crossorigin="anonymous"></script>
    <script src="https://unpkg.com/react-dom@17/umd/react-dom.development.js" integrity="sha512-Wr9OKCTtq1anK0hq5bY3X/AvDI5EflDSAh0mE9gma+4hl+kXdTJPKZ3TwLMBcrgUeoY0s3dq9JjhCQc7vddtFg==" crossorigin="anonymous"></script>
    <link rel="stylesheet" href="https://unpkg.com/graphiql@2.3.0/graphiql.min.css" />
  </head>
  <body>
    <div id="graphiql">Loading...</div>
    <script src="https://unpkg.com/graphiql@2.3.0/graphiql.min.js" type="application/javascript"></script>
    <script>
      ReactDOM.render(
        React.createElement(GraphiQL, {
          fetcher: GraphiQL.createFetcher({url: '/query'}),
          defaultEditorToolsVisibility: true,
        }),
        document.getElementById('graphiql'),
      );
    </script>
  </body>
</html>
`)
