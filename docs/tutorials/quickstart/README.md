# Quickstart

## Overview

In this tutorial we will run the Event Provenance Registry (EPR) Server with
Redpanda Message broker and PostgreSQL database. We will create an event receiver, an event receiver group, and a few events. 

## Requirements

- [Golang 1.21+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/engine/install)
- [Docker-Compose](https://docs.docker.com/engine/install)

## Start dependencies

### Introduction

This section guides you through the process of initiating the backend services
essential for working with EPR. Docker Compose will be employed to launch these
services, with Redpanda serving as the message broker for event transmission and
PostgreSQL as the designated database for data storage.

### Deploy Backend Services

Utilize the provided docker-compose file to launch the required dependencies,
including a PostgreSQL database, a Redpanda Kafka instance, and a Redpanda UI.
Execute the following command:

```bash
docker compose -f ./docker-compose.services.yaml up
```

The Redpanda admin console will be at `http://localhost:8080/overview`

Create a topic. Only needed for initial setup

Using the admin console to create a topic called "epr.dev.events"

Or use the docker container

```bash
docker exec -it redpanda \
    rpk topic create epr.dev.events --brokers=localhost:19092
```

## Start Event Provenance Registry server

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

## REST API

### Create using the REST API

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

The results should look like this:

```json
{ "data": "01HPW0DY340VMM3DNMX8JCQDGN" }
```

We need the ULID of the event receiver in the next step.

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

The results of the command should look like this:

```json
{ "data": "01HPW0GV9PY8HT2Q0XW1QMRBY9" }
```

Event Receiver Groups are a way to group together several event receivers. When
all the event receivers in a group have successful events for a given unit the
event receiver group will produce a message on the topic.

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

### Query using the REST API

We can query the event information using a GET on the events endpoint as
follows:

```bash
curl --header 'Content-Type: application/json' --location \
  --request GET 'http://localhost:8042/api/v1/events/01HPW0GV9PY8HT2Q0XW1QMRBY9'
```

Query the information for an event receiver:

```bash
curl --header 'Content-Type: application/json' --location \
  --request GET 'http://localhost:8042/api/v1/receivers/01HPW0DY340VMM3DNMX8JCQDGN'
```

And query the information for an event receiver group:

```bash
curl --header 'Content-Type: application/json' --location \
  --request GET 'http://localhost:8042/api/v1/groups/01HPW0JXG82Q0FBEC9M8P2Q6J8'
```

## GraphQL

### Access graphql playground

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

The graphql playground will now be accessible at:
<http://localhost:8042/api/v1/graphql>

### Making a request

The current schema for all requests is available through the UI. A simple
mutation and query command can be found below

### Mutation using GraphQL

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
    "create_event_receiver": "01HPVZY1V8SVXGQY03ZG90CA3S"
  }
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
      type: "test.test.test"
    }
  )
}
```

This will return the id of the newly created event receiver group.

```json
{
  "data": {
    "create_event_receiver_group": "01HPW02R3G3QP3EJB036M41J9J"
  }
}
```

Event receiver Groups can be updated using the following mutation

```graphql
mutation {
  set_event_receiver_group_enabled(id: "01HPW02R3G3QP3EJB036M41J9J")
}
```

```graphql
mutation {
  set_event_receiver_group_disabled(id: "01HPW02R3G3QP3EJB036M41J9J")
}
```

Now can create a new event for the event receiver ID in the previous step

```graphql
mutation {
  create_event(
    event: {
      name: "foo"
      version: "1.0.0"
      release: "20231103"
      platform_id: "platformID"
      package: "package"
      description: "The Foo of Brixton"
      payload: "{\"name\": \"value\"}"
      event_receiver_id: "ID_RETURNED_FROM_PREVIOUS_MUTATION"
      success: true
    }
  )
}
```

This will return the id of the newly created event.

```json
{
  "data": {
    "create_event": "01HPW06R5QXK0C2GZM8H442Q9F"
  }
}
```

### Query using GraphQL

This query is only returning a subset of the available fields. Pass in the ID of
the previously created event

```graphql
query {
  events_by_id(id: "01HPW06R5QXK0C2GZM8H442Q9F") {
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
  event_receivers_by_id(id: "01HPVZY1V8SVXGQY03ZG90CA3S") {
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
  event_receiver_groups_by_id(id: "01HPW02R3G3QP3EJB036M41J9J") {
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

We can use graphql to search for events by name and version.

Create another event for the same event receiver

```graphql
mutation {
  create_event(
    event: {
      name: "foo"
      version: "1.0.0"
      release: "20240217"
      platform_id: "x86-64-gnu-linux-7"
      package: "rpm"
      description: "The RPM Foo of Brixton"
      payload: "{\"name\": \"foo\"}"
      event_receiver_id: "01HPVZY1V8SVXGQY03ZG90CA3S"
      success: true
    }
  )
}
```

In the graphql window create a query with the following:

```graphql
query{
  events(event: {name: "foo", version: "1.0.0"}) {
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

This query will return all events with the name foo and version 1.0.0

As follows:

```json
{
  "data": {
    "events": [
      {
        "id": "01HPW06R5QXK0C2GZM8H442Q9F",
        "name": "foo",
        "version": "1.0.0",
        "release": "20231103",
        "platform_id": "platformID",
        "package": "package",
        "description": "The Foo of Brixton",
        "payload": {
          "name": "value"
        },
        "success": true,
        "event_receiver_id": "01HPVZY1V8SVXGQY03ZG90CA3S",
        "created_at": "2024-02-17T12:00:45.62347-05:00"
      },
      {
        "id": "01HPW1WMPK0ZDHQYD6T40MZBGW",
        "name": "foo",
        "version": "1.0.0",
        "release": "20240217",
        "platform_id": "x86-64-gnu-linux-7",
        "package": "rpm",
        "description": "The RPM Foo of Brixton",
        "payload": {
          "name": "foo"
        },
        "success": true,
        "event_receiver_id": "01HPVZY1V8SVXGQY03ZG90CA3S",
        "created_at": "2024-02-17T12:30:11.539643-05:00"
      }
    ]
  }
}
```

We can use graphql to search for event receivers by name and version.

Create another event receiver

```graphql
mutation {
  create_event_receiver(
    event_receiver: {
      name: "the_clash"
      version: "1.0.0"
      type: "london.calling.from.the.far.away.town"
      description: "The only band that matters"
      schema: "{\"name\": \"joe\"}"
    }
  )
}
```

In the graphql window create a query with the following:

```graphql
query {
  event_receivers(event_receiver: {name: "the_clash", version: "1.0.0"})
    {
    id
    name
    version
    type
    description
    created_at
  }
}
```

This query will return all events with the name the_clash and version 1.0.0

As follows:

```graphql
{
  "data": {
    "event_receivers": [
      {
        "id": "01HPVZY1V8SVXGQY03ZG90CA3S",
        "name": "the_clash",
        "version": "1.0.0",
        "type": "london.calling",
        "description": "The only band that matters",
        "created_at": "2024-02-17T11:56:00.616222-05:00"
      },
      {
        "id": "01HPW2BBY8HGPS6JG10DCXE6EH",
        "name": "the_clash",
        "version": "1.0.0",
        "type": "london.calling.from.the.far.away.town",
        "description": "The only band that matters",
        "created_at": "2024-02-17T12:38:14.088635-05:00"
      }
    ]
  }
}
```

### Query using the GraphQL with Curl

We need to craft a GraphQL query. First thing we need is an event receiver. The event receiver acts as a classification and gate for events.

We can find and event reciever by id using the following graphql query:

```json
{
  "query": "query ($er: FindEventReceiverInput!){event_receivers(event_receiver: $er) {id,name,type,version,description}}",
  "variables": {
    "er": {
      "id": "01HPW652DSJBHR5K4KCZQ97GJP"
    }
  }
}
```

We can query the event reciever information using a POST on the graphql endpoint as follows:

```bash
curl -X POST -H "content-type:application/json" -d '{"query":"query ($er: FindEventReceiverInput!){event_receivers(event_receiver: $er) {id,name,type,version,description}}","variables":{"er":{"id":"01HPW652DSJBHR5K4KCZQ97GJP"}}}' http://localhost:8042/api/v1/graphql/query
```

We can query for an event by name and version using the following graphql query:

```json
{
  "query": "query ($e: FindEventInput!){events(event: $e) {id,name,version,release,platform_id,package,description,success,event_receiver_id}}",
  "variables": {
    "e": {
      "name": "foo",
      "version": "1.0.0"
    }
  }
}
```

We can query the event reciever information using a POST on the graphql endpoint as follows:

```bash
curl -X POST -H "content-type:application/json" -d '{"query":"query ($e : FindEventInput!){events(event: $e) {id,name,version,release,platform_id,package,description,success,event_receiver_id}}","variables":{"e": {"name":"foo","version":"1.0.0"}}}' http://localhost:8042/api/v1/graphql/query
```

```bash
curl -X POST -H "content-type:application/json" -d '{"query":"query {events(event: {name: \"foo\", version: \"1.0.0\"}) {id,name,version,release,platform_id,package,description,success,event_receiver_id}}}' http://localhost:8042/api/v1/graphql/query
```

We can query for an event receiver group by name and version using the following graphql query:

```json
{
  "query": "query ($erg: FindEventReceiverGroupInput!){event_receiver_groups(event_receiver_group: $erg) {id,name,type,version,description}}",
  "variables": {
    "erg": {
      "name": "foobar",
      "version": "1.0.0"
    }
  }
}
```

We can query the event reciever information using a POST on the graphql endpoint as follows:

```bash
curl -X POST -H "content-type:application/json" -d '{"query":"query ($erg: FindEventReceiverGroupInput!){event_receiver_groups(event_receiver_group: $erg) {id,name,type,version,description}}","variables":{"erg": {"name":"foobar","version":"1.0.0"}}}' http://localhost:8042/api/v1/graphql/query
```


