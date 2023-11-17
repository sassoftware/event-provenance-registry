# Start Server

This how-to walks you through starting a server and making your first request
using GraphQL.

## Start dependencies

This project contains a docker-compose file that will start up a postgres
database, a redpanda kafka instance, and a redpanda UI. These three dependencies
can be started by running the following command:

```bash
docker-compose -f docs/how-to/start-server/docker-compose.yaml up
```

## Start server

The server can be started using the default settings. This will make the server
available on localhost:8042.

```bash
go run main.go
```

## Access graphql playground

On successful startup the server will display the message below:

```json
{
  "level": "info",
  "module": "cmd.root",
  "v": 0,
  "logger": "server",
  "timestamp": "2023-07-29T13:56:22.378783-04:00",
  "message": "connect to http://localhost:8042/api/v1/graphql for GraphQL playground"
}
```

The graphql playground will not be accessible at:
<http://localhost:8042/api/v1/graphql>

## Making a request

The current schema for all requests is available through the UI. A simple
mutation and query command can be found below

## Mutation

Create an event receiver

```graphql
mutation {
  create_event_receiver(
    event_receiver: {
      name: "grant"
      version: "1.0.0"
      type: "some-action"
      description: "a fake event receiver"
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
      name: "grant"
      version: "1.0.0"
      release: "some-action"
      platformID: "platformID"
      package: "package"
      description: "a fake event receiver"
      payload: "{\"name\": \"value\"}"
      event_receiver_id: "01H6HSJGDJ9CH67D3BK30XD2Q5"
      success: true
    }
  )
}
```

This can then be used to create a new event receiver group

```graphql
mutation {
  create_event_receiver_group(
    event_receiver_group: {
      name: "grant"
      version: "1.0.0"
      description: "a fake event receiver"
      event_receiver_ids: ["ID_RETURNED_FROM_PREVIOUS_MUTATION"]
      type: "test"
    }
  )
}
```

This will return the id of the newly created event receiver group.

```json
{
  "data": {
    "create_event_receiver_group": "01H713QGDGCW546NV7QYEK3QJ7"
  }
}
```

Event receiver Groups can be updated using the following mutation

```graphql
mutation {
  set_event_receiver_group_enabled(id: "01H713QGDGCW546NV7QYEK3QJ7")
}
```

```graphql
mutation {
  set_event_receiver_group_disabled(id: "01H713QGDGCW546NV7QYEK3QJ7")
}
```

## Query

This query is only returning a subset of the available fields. Pass in the ID of
the previously created event

```graphql
query {
  event(id: "01HEK4REW0S93V8Y10E8A8M8HC") {
    name
    version
    description
  }
}
```

This query is only returning a subset of the available fields. Pass in the ID of
the previously created event_receiver

```graphql
query {
  event_receiver(id: "01H6HSJGDJ9CH67D3BK30XD2Q5") {
    name
    version
    description
  }
}
```

This query is only returning a subset of the available fields. Pass in the ID of
the previously created event_receiver_group

```graphql
query {
  event_receiver_group(id: "01H6HSJGDJ9CH67D3BK30XD2Q5") {
    name
    version
    description
  }
}
```
