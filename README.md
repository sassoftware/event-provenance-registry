# Async Event Infrastructure Server

## Overview

A server for accepting events, storing events, and producing messages on a
message bus.

See these pages for specific details about different components:

[Database](pkg/storage/README.md)
[GraphQL API](pkg/graph/README.md)

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

Create an event receiver:

```bash
curl --location --request POST 'http://localhost:8080/api/v1/receivers' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "foobar",
  "type": "whatever",
  "version": "1.1.2",
  "description": "it does stuff",
  "enabled": true
}'
```

Create an event:

```bash
curl --location --request POST 'http://localhost:8080/api/v1/events' \
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
    "event_receiver_id": "01H5VWNWQ6ETAEW7DT6CMD5EC0"
}'
```

Create a group:

```bash
curl --location --request POST 'http://localhost:8080/api/v1/groups' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "group",
    "type": "notyours",
    "version": "3.3.3",
    "description": "foobar",
    "enabled": true,
    "event_receiver_ids": ["01H5VWNWQ6ETAEW7DT6CMD5EC0"]
}'
```

## Contributing

We welcome your contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md)
for details on how to submit contributions to this project.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
