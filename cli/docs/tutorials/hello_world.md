# Hello World

## Overview

In this how-to we will cover getting started with the EPR CLI.

## Install the CLI

In this step we will build the EPR CLI and install it into our `GOPATH`.

### Build

```bash
make
```

### Installation

To install in your go path directory set `PREFIX`

Linux

```bash
make PREFIX=$(go env GOPATH) install
```

Mac OS X

```bash
make PREFIX=$(go env GOPATH) install-darwin
```

## Hello World

## Requirements

This tutorial requires that the
[Hello World](../../../docs/tutorials/hello_world/README.md) has been completed
and the EPR server is running.

### Create Event Receivers

First thing we create are a few event receivers.

The commands to create an event receiver are as follows:

```bash
epr-cli receiver create --name "foo-cli" --version "1.0.0" --description "foo cli created foo" --type "epr.foo.cli" --schema "{}" --dry-run
```

The first command will return the structure of the event receiver. We can
validate that the receiver has the data we expect before creating it.

```bash
epr-cli receiver create --name "foo-cli" --version "1.0.0" --description "foo cli created foo" --type "epr.foo.cli" --schema "{}"
```

The second commands will return the id of the newly created `foo-cli` event
receiver.

```json
{
  "data": "01HKX0J9KS8AASMRYX61458N41"
}
```

We can validate the new event receiver exists with the search command using the
`ID` from the create reciever command.

```bash
epr-cli receiver search --id 01HKX0J9KS8AASMRYX61458N41 --fields all
```

Now we will create a second event receiver as follows.

```bash
epr-cli receiver create --name "bar-cli" --version "1.0.0" --description "bar cli created bar" --type "epr.bar.cli" --schema "{}"  --dry-run

epr-cli receiver create --name "bar-cli" --version "1.0.0" --description "bar cli created bar" --type "epr.bar.cli" --schema "{}"
```

```json
{
  "data": "01HKX0KY3B31MR3XKJWTDZ4EQ0"
}
```

And validate we have created the `bar-cli` receiver.

```bash
epr-cli receiver search --id 01HKX0KY3B31MR3XKJWTDZ4EQ0 --fields all
```

We can also search for all the event receivers that have a type of `epr.foo.cli`

```bash
epr-cli receiver search --type epr.foo.cli --fields all
```

And search for event receivers by name and version

```bash
epr-cli receiver search --name foo-cli --version 1.0.0 --fields all
```

### Create Event Receiver Groups

Next we Next we will create an event receiver group.

The commands to create an event receiver group are as follows:

```bash
epr-cli group create --name "foo-group-cli" --version "1.0.0" --description "foo cli created foo group" --type "epr.foo.group.cli" --receiver-ids "01HKX0J9KS8AASMRYX61458N41 01HKX0KY3B31MR3XKJWTDZ4EQ0"  --dry-run

epr-cli group create --name "foo-group-cli" --version "1.0.0" --description "foo cli created foo group" --type "epr.foo.group.cli" --receiver-ids "01HKX0J9KS8AASMRYX61458N41 01HKX0KY3B31MR3XKJWTDZ4EQ0"
```

```json
{
  "data": "01HKX90FKWQZ49F6H5V5NQT95Z"
}
```

And validate we have created the `foo-group-cli` receiver group.

```bash
epr-cli group search --id 01HKX90FKWQZ49F6H5V5NQT95Z --fields all
```

### Create Events

Now that we have created the receivers and receiver groups we can create events.

The commands to create a `foo` event are as follows:

```bash
epr-cli event create --name foo --version 1.0.0 --release 2024.01 --platform-id x86-64-gnu-linux-9 --package rpm --success true --description "the foo event for foo" --success true --receiver-id 01HKX0J9KS8AASMRYX61458N41 --payload '{"name":"foo"}' --dry-run

epr-cli event create --name foo --version 1.0.0 --release 2024.01 --platform-id x86-64-gnu-linux-9 --package rpm --success true --description "the foo event for foo" --success true --receiver-id 01HKX0J9KS8AASMRYX61458N41 --payload '{"name":"foo"}'
```

We can validate we have created the `foo` event with the following command:

```bash
epr-cli event search --id 01HKX1TMQZQDS6NC5DG7WNXXCJ --fields all
```

or search by name and version

```bash
epr-cli event search --name foo --version 1.0.0 --fields all
```

Next we can create a `bar` event and post it to the other receiver in the
receiver group.

The commands to create a `bar` event are as follows:

```bash
epr-cli event create --name bar --version 1.0.0 --release 2024.01 --platform-id x86-64-gnu-linux-9 --package rpm --success true --description "the bar event for bar" --success true --receiver-id 01HKX0KY3B31MR3XKJWTDZ4EQ0 --payload '{"name":"bar"}' --dry-run

epr-cli event create --name bar --version 1.0.0 --release 2024.01 --platform-id x86-64-gnu-linux-9 --package rpm --success true --description "the bar event for bar" --success true --receiver-id 01HKX0KY3B31MR3XKJWTDZ4EQ0 --payload '{"name":"bar"}'
```

```json
{
  "data": "01HKX7HKVDVH0HRQN0P2VDJ7Y6"
}
```

We can validate we have created the `bar` event with the following command:

```bash
epr-cli event search --id 1HKX7HKVDVH0HRQN0P2VDJ7Y6 --fields all
```

Or we can search by name and version and success

```bash
epr-cli event search --name bar --version 1.0.0 --success true --fields all
```

### Triggering the Event Receiver Group Message

By posting the `foo` event to the first receiver and `bar` to the second
receiver we trigger the receiver group message to fire. We should be able to see
the message in the redpanda UI and in the event queue.
