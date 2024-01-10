# Start Server

This how-to walks you through starting a server and making your first request
using GraphQL.

## Start dependencies

This project contains a docker-compose file that will start up a postgres
database, a redpanda kafka instance, and a redpanda UI. These three dependencies
can be started by running the following command:

```bash
docker-compose -f ./docker-compose.services.yaml up
```

## Start server

Export the environment variables for the server

```bash
export EPR_TOPIC=epr.dev.events
export EPR_BROKERS=localhost:19092
export EPR_DB=postgres://localhost:5432
```

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
      name: "foo-receiver"
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
    "create_event_receiver": "01HKNDR10NVBA8V7G0V3C15JA6"
  }
}
```

This can then be used to create a new event

```graphql
mutation {
  create_event(
    event: {
      name: "foo-event"
      version: "1.0.0"
      release: "some-action"
      platform_id: "platform-x64"
      package: "package"
      description: "a fake event"
      payload: "{\"name\": \"value\"}"
      event_receiver_id: "01HKNDR10NVBA8V7G0V3C15JA6"
      success: true
    }
  )
}
```

This will return the id of the newly created event.

```json
{
  "data": {
    "create_event": "01HKNDTSFT6ZZ8Q8YNK736TT43"
  }
}
```

We can use the event receiver to create a new event receiver group

```graphql
mutation {
  create_event_receiver_group(
    event_receiver_group: {
      name: "foo-group"
      version: "1.0.0"
      description: "a fake event receiver group"
      event_receiver_ids: ["01HKNDR10NVBA8V7G0V3C15JA6"]
      type: "test"
    }
  )
}
```

This will return the id of the newly created event receiver group.

```json
{
  "data": {
    "create_event_receiver_group": "01HKNE0TJG7GA35GP703D75XTH"
  }
}
```

Event receiver Groups can be updated using the following mutation

```graphql
mutation {
  set_event_receiver_group_enabled(id: "01HKNE0TJG7GA35GP703D75XTH")
}
```

```graphql
mutation {
  set_event_receiver_group_disabled(id: "01HKNE0TJG7GA35GP703D75XTH")
}
```

## Query

This query is only returning a subset of the available fields. Pass in the ID of
the previously created event

```graphql
query {
  events(id: "01HKNDTSFT6ZZ8Q8YNK736TT43") {
    id
    name
    version
    release
    platform_id
    package
    description
    payload
    success
    event_receiver_id
		created_at
  }
}
```

This query is only returning a subset of the available fields. Pass in the ID of
the previously created event_receiver

```graphql
query {
  event_receivers(id: "01HKNDR10NVBA8V7G0V3C15JA6") {
    name
    version
    description
    type
    schema
    fingerprint
    created_at
  }
}
```

This query is only returning a subset of the available fields. Pass in the ID of
the previously created event_receiver_group

```graphql
query {
  event_receiver_groups(id: "01HKNE0TJG7GA35GP703D75XTH") {
    name
    version
    description
    type
    enabled
    event_receiver_ids
    created_at
    updated_at
  }
}
```
