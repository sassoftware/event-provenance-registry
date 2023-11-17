# Hello World

## Overview

In this tutorial we will run the Event Provenance Registry (EPR) Server with
Redpanda and create a event receiver and a few events.

## Requirements

Golang 1.20+ Docker Docker-Compose

## Start dependencies

This project contains a docker-compose file that will start up a postgres
database, a redpanda kafka instance, and a redpanda UI. These three dependencies
can be started by running the following command:

## Start Redpanda

This how-to walks you through starting a server and making your first request
using GraphQL.

```bash
docker compose -f ./docs/tutorials/workshops/compose/docker-compose.yaml up
```

The redpanda admin console will be at `http://localhost:8080/overview`

Create a topic. Only needed for initial setup

Using the admin console to create a topic called "epr.dev.events"

Or use the docker container

```bash
docker exec -it redpanda-1 \
    rpk topic create epr.dev.events --brokers=localhost:9092
```

## Start Event Provenance Registry server

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

## Mutation using GraphQL

Create an event receiver

```graphql
mutation {
  create_event_receiver(
    event_receiver: {
      name: "the_clash"
      version: "1.0.0"
      type: "london.calling"
      description: "The only band that matters"
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
      name: "foo"
      version: "1.0.0"
      release: "20231103"
      platformID: "platformID"
      package: "package"
      description: "The Foo of Brixton"
      payload: "{\"name\": \"value\"}"
      event_receiver_id: "01HEK2ZJHQJZPCTJRY88B04RVT"
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
      name: "foobar"
      version: "1.0.0"
      description: "a fake event receiver group"
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

## Query using GraphQL

This query is only returning a subset of the available fields. Pass in the ID of
the previously created event

```graphql
query {
  event(id: "01H6HSJGDJ9CH67D3BK30XD2Q5") {
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

## Create using the REST API

First thing we need is an event receiver. The event receiver acts as a
classification and gate for events.

Create an event receiver:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/receivers' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "foobar",
  "type": "foo.bar",
  "version": "1.1.2",
  "description": "The event receiver of Brixton",
  "enabled": true,
  "schema": {
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    }
  }
}
}'
```

Create an event using curl.

When you create an event, you must specify an `event_receiver_id` to associate
it with. An event is the record of some action being completed. You cannot
create an event for a non-existent receiver ID. The payload field of the event
must conform to the schema defined on the event receiver that you have given the
ID of.

Create an event:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/events' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "magnificent",
    "version": "7.0.1",
    "release": "2023.11.16",
    "platform_id": "linux",
    "package": "docker",
    "description": "blah",
    "payload": {"name":"joe"},
    "success": true,
    "event_receiver_id": "<PASTE EVENT RECEIVER ID FROM FIRST CURL COMMAND>"
}'
```

Create an event receiver group:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/groups' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "the_clash",
    "type": "foo.bar",
    "version": "3.3.3",
    "description": "The only event receiver group that matters",
    "enabled": true,
    "event_receiver_ids": ["PASTE EVENT RECEIVER ID FROM FIRST CURL COMMAND"]
}'
```
