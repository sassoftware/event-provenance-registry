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

[EPR Documentation](./docs/README.md)

## Develop

### Build

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

```txt
Usage:
  epr-server [flags]

Flags:
      --brokers string   broker uris separated by commas (default "localhost:9092")
      --config string    config file (default is $XDG_CONFIG_HOME/epr/epr.yaml)
      --db string        database connection string (default "postgres://localhost:5432")
      --debug            Enable debugging statements
  -h, --help             help for epr-server
      --host string      host to listen on (default "localhost")
      --port string      port to listen on (default "8042")
      --topic string     topic to produce events on (default "epr.dev.events")
```

### Running the Server

See [the docs](docs/how-to/start-server/README.md) for how to start EPR
and its dependencies.

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

## GraphQL

Curl commands for GraphQL Search

```bash
curl -X POST -H "content-type:application/json" -d '{"query":"query FindEvents($id: ID!){events(id: $id) {id,name,version,release,platform_id,package,description,payload,success,created_at,event_receiver_id}}","variables":{"id":"01HKX1TMQZQDS6NC5DG7WNXXCJ"}}' http://localhost:8042/api/v1/graphql/query
```

```bash
curl -X POST -H "content-type:application/json" -d '{"query":"query FindEventReceivers($id: ID!){event_receivers(id: $id) {id,name,version,description,type,schema,created_at}}","variables":{"id":"01HKX0KY3B31MR3XKJWTDZ4EQ0"}}' http://localhost:8042/api/v1/graphql/query
```

```bash
curl -X POST -H "content-type:application/json" -d '{"query":"query FindEventReceiverGroups($id: ID!){event_receiver_groups(id: $id) {id,name,type,version,description,enabled,event_receiver_ids,created_at,updated_at}}","variables":{"id":"01HKX90FKWQZ49F6H5V5NQT95Z"}}' http://localhost:8042/api/v1/graphql/query
```

## Keycloak

Follow this guide to get started https://www.keycloak.org/getting-started/getting-started-docker.

That'll get you a new user and client. Unlike the example, set `client authentication` to `on`, otherwise, you can't
access the client credentials tab in the admin console http://localhost:8083/admin/master/console/#/<realm>/clients/<ulid>/credentials.

To start up Keycloak. I changed the port number because the default 8080 is in use by EPR itself.

```bash
docker run -p 8083:8083 -e KC_HTTP_PORT=8083 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin \
  quay.io/keycloak/keycloak:24.0.1 start-dev
```

http://localhost:8083/admin/master/console/#/test/clients
https://github.com/coreos/go-oidc
https://github.com/coreos/go-oidc/blob/v3/example/userinfo/app.go

Need to grab a token from Keycloak first.  You can get oidc config info from here http://localhost:8083/realms/test/.well-known/openid-configuration
Log into your new realm with your account from here: http://localhost:8083/realms/<realmName>/account.

Once you've got that, 

```bash
export access_token=$(\
  curl -X POST http://localhost:8083/realms/test/protocol/openid-connect/token \
  --user epr-client-id:rGMO0kRpvUj3XD9It678AoTlgMtGxItJ \
  -d 'username=testbob&password=abc123&grant_type=password' | jq --raw-output '.access_token'\
)
```

The `--user` flag is the credentials for the client connecting to Keycloak. In our case, that'll be EPR. The payload
content depends on the grant type. For our purposes, password is easiest. The username and password will be those of
the user you created in the tutorial.

I think the flow is roughly:

1. User hits a logon endpoint and gets redirected to Keycloak.
2. EPR passes up client credentials plus user info.
3. Keycloak returns a token
4. EPR passes token to user.

Although, I think we could shorten it to just get a JWT straight from Keycloak without passing it back through EPR...
I'll have to play with it.

## Contributing

We welcome your contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md)
for details on how to submit contributions to this project.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
