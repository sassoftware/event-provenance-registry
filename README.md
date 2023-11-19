# Event Provenance Registry (EPR)

## Overview

Event Provenance Registry (EPR) is a provenance store. It keeps track of events.
With EPR, you can use the API to create, retrieve, and query events, event
receivers, and groups of event receivers.

## Description

The Event Provenance Registry (EPR) is a service that stores events and tracks
event-receivers and event-receiver-groups. EPR provides an API that lets you
create events, event-receivers, and event-receiver-groups. You can query the EPR
using the GraphQL endpoint to get identifying information about events,
event-receivers, and event-receiver-groups. EPR collects events from the supply
chain to record the lifecycle of a unit in the SDLC. EPR validates
event-receivers have events. EPR emits a message when a event is received as
well as when an event-receiver-groups is complete for a unit version.

## Build

```bash
make
```

## Installation

To install in `/usr/local/bin` directory

Linux

```bash
make install
```

Mac OS X

```bash
make install-darwin
```

To install in your go path directory set `PREFIX`

Linux

```bash
make PREFIX=$(go env GOPATH) install
```

Mac OS X

```bash
make PREFIX=$(go env GOPATH) install-darwin
```

## Tests

Run the go unit tests:

```bash
make test
```

## Linter

Run golangci-lint (requires
[golangci-lint](https://golangci-lint.run/usage/install/)):

```bash
make megalint
```

## Usage

### Running the Server

Start up Redpanda. See
[the docs](docs/how-to/redpanda/multi-node/redpanda_deploy.md) for more details.

```bash
docker compose -f docs/how-to/redpanda/multi-node/docker-compose.yaml up -d
```

[Start up Postgres.](docs/how-to/start-server/README.md)

### Interacting with the Server

An event receiver is an object that represents some type of action that would
occur in your pipeline (i.e. a build, a test run, a deployment, etc...).

Create an event receiver:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/receivers' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "foobar",
  "type": "whatever",
  "version": "1.1.2",
  "description": "it does stuff",
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
    "name": "idc",
    "version": "2.2.2",
    "release": "whatever",
    "platform_id": "linux",
    "package": "docker",
    "description": "blah",
    "payload": {"name":"bob"},
    "success": true,
    "event_receiver_id": "01HDS785T0V8KTSTDM9XGT33QQ"
}'
```

Create a group:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/groups' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "group",
    "type": "notyours",
    "version": "3.3.3",
    "description": "foobar",
    "enabled": true,
    "event_receiver_ids": ["01H9GW7FYY4XYE2R930YTFM7FM"]
}'
```

## Contributing

We welcome your contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md)
for details on how to submit contributions to this project.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
