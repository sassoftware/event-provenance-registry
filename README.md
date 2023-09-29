# Async Event Infrastructure Server

## Overview

A server for accepting events, storing events, and producing messages on a
message bus.

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

## Usage

### Running the Server

Start up Redpanda. See [the docs](docs/how-to/redpanda/multi-node/redpanda_deploy.md) for more details.

```bash
docker compose -f docs/how-to/redpanda/multi-node/docker-compose.yaml up -d
```

[Start up Postgres.](docs/how-to/start-server/start_and_request_server.md)

### Interacting with the Server

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
  "schema": {}
}'
```

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
    "payload": "",
    "success": true,
    "event_receiver_id": "01H9GW7FYY4XYE2R930YTFM7FM"
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
