# Event Provenance Registry CLI

## Overview

Event Provenance Registry CLI is a command line for interfacing with the Event
Provenance Registry service that stores events and tracks event-receivers and
event-receiver-groups.

## Description

The Event Provenance Registry CLI is a command line for interfacing with the
Event Provenance Registry service that stores events and tracks event-receivers
and event-receiver-groups. EPR CLI provides an command line that lets you create
events, event-receivers, and event-receiver-groups. You can query the EPR using
the EPR CLI to get identifying information about events, event-receivers, and
event-receiver-groups.

[EPR Documentation](../docs/README.md)

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

Check the status of EPR including liveness and readiness

```txt
Usage:
  epr-cli status [flags]

Flags:
      --dry-run      do a dry run of the command
  -h, --help         help for status
      --no-indent    do not indent the JSON output
      --url string   EPR base url (default "http://localhost:8042")
```

Create, Search, and Generate Events for the Event Provenance Registry Server

```text
Usage:
  epr-cli event [command]

Available Commands:
  create      create an event
  generate    Generates a event
  search      Searches for event objects

Flags:
  -h, --help   help for event
```

Create, Search, or Generate Event Receivers for the Event Provenance Registry
Service

```text
Usage:
  epr-cli receiver [command]

Available Commands:
  create      Creates a Event Receiver
  generate    Generates a Event Receiver
  search      Search for Event Receivers

Flags:
  -h, --help   help for receiver
```

Create, Search, or Generate Event Receiver Groups for the Event Provenance
Registry Service

```text
Usage:
  epr-cli group [command]

Available Commands:
  create      creates a Event Receiver Group
  generate    Generates a Event Receiver Group
  search      Searches for Event Receiver Group objects

Flags:
  -h, --help   help for group
```

## Examples

Create Event Receivers

```bash
./bin/epr-cli-darwin-arm64 receiver create --name "foo-cli" --version "1.0.0" --description "foo cli created foo" --type "epr.foo.cli" --schema "{}" --dry-run

./bin/epr-cli-darwin-arm64 receiver create --name "foo-cli" --version "1.0.0" --description "foo cli created foo" --type "epr.foo.cli" --schema "{}"
```

```json
{
  "data": "01HKX0J9KS8AASMRYX61458N41"
}
```

```bash
./bin/epr-cli-darwin-arm64 receiver search --id 01HKX0J9KS8AASMRYX61458N41 --fields all
```

```bash
./bin/epr-cli-darwin-arm64 receiver create --name "bar-cli" --version "1.0.0" --description "bar cli created bar" --type "epr.bar.cli" --schema "{}"  --dry-run

./bin/epr-cli-darwin-arm64 receiver create --name "bar-cli" --version "1.0.0" --description "bar cli created bar" --type "epr.bar.cli" --schema "{}"
```

```json
{
  "data": "01HKX0KY3B31MR3XKJWTDZ4EQ0"
}
```

```bash
./bin/epr-cli-darwin-arm64 receiver search --id 01HKX0KY3B31MR3XKJWTDZ4EQ0 --fields all
```

Create Events

```bash
./bin/epr-cli-darwin-arm64 event create --name foo --version 1.0.0 --release 2024.01 --platform-id x86-64-gnu-linux-9 --package rpm --success true --description "the foo event for foo" --success true --receiver-id 01HKX0J9KS8AASMRYX61458N41 --payload '{"name":"foo"}' --dry-run

./bin/epr-cli-darwin-arm64 event create --name foo --version 1.0.0 --release 2024.01 --platform-id x86-64-gnu-linux-9 --package rpm --success true --description "the foo event for foo" --success true --receiver-id 01HKX0J9KS8AASMRYX61458N41 --payload '{"name":"foo"}'
```

```bash
./bin/epr-cli-darwin-arm64 event search --id 01HKX1TMQZQDS6NC5DG7WNXXCJ --fields all
```

```bash
./bin/epr-cli-darwin-arm64 event create --name bar --version 1.0.0 --release 2024.01 --platform-id x86-64-gnu-linux-9 --package rpm --success true --description "the bar event for bar" --success true --receiver-id 01HKX0KY3B31MR3XKJWTDZ4EQ0 --payload '{"name":"bar"}' --dry-run

./bin/epr-cli-darwin-arm64 event create --name bar --version 1.0.0 --release 2024.01 --platform-id x86-64-gnu-linux-9 --package rpm --success true --description "the bar event for bar" --success true --receiver-id 01HKX0KY3B31MR3XKJWTDZ4EQ0 --payload '{"name":"bar"}'
```

```json
{
  "data": "01HKX7HKVDVH0HRQN0P2VDJ7Y6"
}
```

```bash
./bin/epr-cli-darwin-arm64 event search --id 1HKX7HKVDVH0HRQN0P2VDJ7Y6 --fields all
```

Create Event Receiver Groups

```bash
./bin/epr-cli-darwin-arm64 group create --name "foo-group-cli" --version "1.0.0" --description "foo cli created foo group" --type "epr.foo.group.cli" --receiver-ids "01HKX0J9KS8AASMRYX61458N41 01HKX0KY3B31MR3XKJWTDZ4EQ0"  --dry-run

./bin/epr-cli-darwin-arm64 group create --name "foo-group-cli" --version "1.0.0" --description "foo cli created foo group" --type "epr.foo.group.cli" --receiver-ids "01HKX0J9KS8AASMRYX61458N41 01HKX0KY3B31MR3XKJWTDZ4EQ0"
```

```json
{
  "data": "01HKX90FKWQZ49F6H5V5NQT95Z"
}
```

```bash
./bin/epr-cli-darwin-arm64 group search --id 01HKX90FKWQZ49F6H5V5NQT95Z --fields all
```
