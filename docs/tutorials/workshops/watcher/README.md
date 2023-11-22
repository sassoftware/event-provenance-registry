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
	"encoding/json"
	"log"

	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/watcher"
)

func main() {
	seeds := []string{"localhost:19092"}
	topics := []string{"epr.dev.events"}
	consumerGroup := "watcher-workshop"

	watcher, err := watcher.New(seeds, topics, consumerGroup)
	if err != nil {
		panic(err)
	}
	defer watcher.Client.Close()

	go watcher.StartTaskHandler(customTaskHandler)

	watcher.ConsumeRecords(customMatcher)
}

ffunc customMatcher(msg *message.Message) bool {
	return msg.Type == "foo.bar"
}

func customTaskHandler(msg *message.Message) error {
	log.Default().Printf("I received a task with value '%v'", msg)
	return nil
}

```

Save the file and run `go mod init` then

Now we can run `go mod tidy` to fill in our dependencies.

## Begin consuming

We can now start up the watcher and start consuming messages.

```bash
go run main.go
```

You should see a log stating that we have begin consuming records.

## Create an event receiver

Create an event receiver:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/receivers' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "watcher-workshop",
  "type": "foo.bar",
  "version": "1.0.0",
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

## Produce message

In a second terminal run the command below:

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

## Receive message

You should now see a message like the one below.

```bash
2023/11/17 16:18:30 I received a task with value '{"success":true,"id":"01HFFJCJYZN02RR1JSCE9DDAS4","specversion":"1.0","type":"foo.bar","source":"","api_version":"v1","name":"magnificent","version":"7.0.1","release":"2023.11.16","platform_id":"linux","package":"docker","data":{"events":[{"id":"01HFFJCJYZN02RR1JSCE9DDAS4","name":"magnificent","version":"7.0.1","release":"2023.11.16","platform_id":"linux","package":"docker","description":"blah","payload":{"name":"joe"},"success":true,"created_at":"16:18:30.000879894","event_receiver_id":"01HFFJ69HHJ506SRDYQMFF1H5A","EventReceiver":{"id":"01HFFJ69HHJ506SRDYQMFF1H5A","name":"watcher-workshop","type":"foo.bar","version":"1.0.0","description":"The event receiver of Brixton","schema":{"type":"object","properties":{"name":{"type":"string"}}},"fingerprint":"b183c34c7ba56b17f89dfe0c0b22c0a340889cae88d8e87a3f16bc5bdc8f7acb","created_at":"16:15:04.000626147"}}],"event_receivers":[{"id":"01HFFJ69HHJ506SRDYQMFF1H5A","name":"watcher-workshop","type":"foo.bar","version":"1.0.0","description":"The event receiver of Brixton","schema":{"type":"object","properties":{"name":{"type":"string"}}},"fingerprint":"b183c34c7ba56b17f89dfe0c0b22c0a340889cae88d8e87a3f16bc5bdc8f7acb","created_at":"16:15:04.000626147"}],"event_receiver_groups":null}}
```

**Note**: the matcher being run is looking for kafka messages with the value
`match`. All other messages will be ignored.
