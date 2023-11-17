# Watcher

## Overview

In this tutorial we will explore using the watcher sdk to create a watcher to
listen for events from the EPR server.

This tutorial depends on the [Hello World](../hello_world/README.md) being
completed.

## Create a new watcher

Make a new directory for your watcher and create a `main.go` in that directory.

```bash
mkdir foo
cd foo
touch main.go
```

Now open the `main.go` in your favorite editor (Vim).

Add the following code:

```go
package main

import (
	"log"

	"github.com/sassoftware/event-provenance-registry/pkg/watcher"
)

func main() {
	seeds := []string{"localhost:9092"}
	topics := []string{"example.topic"}
	consumerGroup := "my-group-identifier"

	watcher, err := watcher.New(seeds, topics, consumerGroup)
	if err != nil {
		panic(err)
	}
	defer watcher.Client.Close()

	go watcher.StartTaskHandler(customTaskHandler)

	watcher.ConsumeRecords(customMatcher)
}

func customMatcher(record *watcher.Record) bool {
    log.Default().Printf("I see a message with the value '%s'", record.Value)
	return string(record.Value) == "match"
}

func customTaskHandler(record *watcher.Record) error {
	log.Default().Printf("I received a task with value '%s'", record.Value)
	return nil
}

```

Save the file and run `go mod init` then

Add the following line to the `go.mod` file under the module line.

```bash
replace github.com/sassoftware/event-provenance-registry => ../../../../../

```

File should look like this:

```bash
module github.sas.com/sassoftware/event-provenance-registry/docs/tutorials/workshops/watcher/foo

replace github.com/sassoftware/event-provenance-registry => ../../../../../

go 1.21.4
```

Now we can run `go mode tidy` to fill in our dependencies.

## Begin consuming

We can now start up the watcher and start consuming messages.

```bash
go run main.go
```

You should see a log stating that we have begin consuming records.

## Produce message

In a second terminal run the command below:

```bash
docker exec -it redpanda-0 \
    rpk topic produce epr.dev.events --brokers=localhost:9092
```

Type a message you would like to send and enter.

```bash
match

```

When you are finished you can exit the producer with `Ctrl+C`

## Receive message

You should now see a message like the one below.

```bash
2023/09/01 22:11:19 I received a task with value 'match'
```

**Note**: the matcher being run is looking for kafka messages with the value
`match`. All other messages will be ignored.
