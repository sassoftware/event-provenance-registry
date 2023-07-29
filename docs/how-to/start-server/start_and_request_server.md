# Start Server

This how-to walks you through starting a server and making your first request
using GraphQL.

## Start database

The database must be started before the server can be started. This can be done
using the command below. This command will start a postgres instance on
port 5432 with no authentication. Be sure your docker daemon is available
prior to running the command.

```bash
docker run -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 postgres
```

## Start server

The server can be started using the default settings. This will make the server
available on localhost:8080.

```bash
go run main.go
```

## Access graphql playground

On successful startup the server will display the message below:

```json
{"level":"info","module":"cmd.root","v":0,"logger":"server","timestamp":"2023-07-29T13:56:22.378783-04:00","message":"connect to http://localhost:8080/api/v1/graphql for GraphQL playground"}
```

The graphql playground will not be accessible at: <http://localhost:8080/api/v1/graphql>

## Making a request

The current schema for all requests is availiable through the UI. A simple mutation
and query command can be found below

## Mutation

```graphql
mutation {
  create_event_receiver(
    event_receiver: {
      name: "grant", 
      version: "1.0.0", 
      type: "some-action", 
      description: "a fake event reciever", 
      schema: "{\"name\": \"value\"}"
    }
  )
}
```

This will return the id of the newly created event receiver.

```json
{
  "data": {
    "create_event_receiver": "01H6HSPWNMR8HJ9WKA5AJWG430"
  }
}
```

This can then be used to create a new event

```graphql
mutation {
  create_event(
    event: {
      name: "grant",
      version: "1.0.0",
      release: "some-action",
      platformID: "platformID",
      package: "package",
      description: "a fake event reciever",
      payload: "{\"name\": \"value\"}"
      event_receiver_id: "01H6HSJGDJ9CH67D3BK30XD2Q5",
      success: true
    }
  )
}
```

## Query

This query is only returning a subset of the available fields. Pass
in the ID of the previously created event

```graphql
{
  event(id: "01H6HSJGDJ9CH67D3BK30XD2Q5") {
    name
    version
    description
  }
}
```
