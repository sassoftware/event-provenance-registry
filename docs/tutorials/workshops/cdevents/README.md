# CDEvents

## Overview

[CDEvents](https://cdevents.dev/) is a common specification for Continuous
Delivery events, enabling interoperability in the complete software production
ecosystem.

In this tutorial we will learn how we can use CDEvents with the Event Provenance
Registry (EPR).

## Requirements

The [Hello World](../hello_world/README.md) has been completed and the EPR
server is running.

## Create the Event Receiver

First we will create the event receiver and apply the cdevents schema for
artifact packaged.

Create the event receiver:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/receivers' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "artifact-packaged",
  "type": "dev.cdevents.artifact.packaged.0.1.1",
  "version": "1.0.0",
  "description": "CDEvents Artifact Packaged",
  "enabled": true,
  "schema": {
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://cdevents.dev/0.4.0-draft/schema/artifact-packaged-event",
  "properties": {
    "context": {
      "properties": {
        "version": {
          "type": "string",
          "minLength": 1
        },
        "id": {
          "type": "string",
          "minLength": 1
        },
        "source": {
          "type": "string",
          "minLength": 1,
          "format": "uri-reference"
        },
        "type": {
          "type": "string",
          "enum": [
            "dev.cdevents.artifact.packaged.0.1.1"
          ],
          "default": "dev.cdevents.artifact.packaged.0.1.1"
        },
        "timestamp": {
          "type": "string",
          "format": "date-time"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": [
        "version",
        "id",
        "source",
        "type",
        "timestamp"
      ]
    },
    "subject": {
      "properties": {
        "id": {
          "type": "string",
          "minLength": 1
        },
        "source": {
          "type": "string",
          "minLength": 1,
          "format": "uri-reference"
        },
        "type": {
          "type": "string",
          "minLength": 1,
          "enum": [
            "artifact"
          ],
          "default": "artifact"
        },
        "content": {
          "properties": {
            "change": {
              "properties": {
                "id": {
                  "type": "string",
                  "minLength": 1
                },
                "source": {
                  "type": "string",
                  "minLength": 1,
                  "format": "uri-reference"
                }
              },
              "additionalProperties": false,
              "type": "object",
              "required": [
                "id"
              ]
            }
          },
          "additionalProperties": false,
          "type": "object",
          "required": [
            "change"
          ]
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": [
        "id",
        "type",
        "content"
      ]
    },
    "customData": {
      "oneOf": [
        {
          "type": "object"
        },
        {
          "type": "string",
          "contentEncoding": "base64"
        }
      ]
    },
    "customDataContentType": {
      "type": "string"
    }
  },
  "additionalProperties": false,
  "type": "object",
  "required": [
    "context",
    "subject"
  ]
}
}'
```

The results of the command should look like this:

```json
{ "data": "01HFW5HXQR28951NR8NH3WJBN6" }
```

Next we will POST and event to the event receiver. The event payload will be in
the form of an artifact published event.

Create an event:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/events' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "foo",
    "version": "1.0.1",
    "release": "2023.11.16",
    "platform_id": "aarch64-gnu-linux-7",
    "package": "oci",
    "description": "packaged oci image foo",
    "payload": {
  "context": {
    "version": "0.4.0-draft",
    "id": "271069a8-fc18-44f1-b38f-9d70a1695819",
    "source": "/event/source/123",
    "type": "dev.cdevents.artifact.packaged.0.1.1",
    "timestamp": "2023-03-20T14:27:05.315384Z"
  },
  "subject": {
    "id": "pkg:golang/mygit.com/myorg/myapp@234fd47e07d1004f0aed9c",
    "source": "/event/source/123",
    "type": "artifact",
    "content": {
      "change": {
        "id": "myChange123",
        "source": "my-git.example/an-org/a-repo"
      }
    }
  }
}
    ,
    "success": true,
    "event_receiver_id": "01HFFDS17FA20PZRWR23KHPK9Y"
}'
```

The results of the command should look like this:

```json
{ "data": "01HFW5MZARPAQME9M9VKC3Z2ZD" }
```

Now we send an event with a payload that doesn't match the schema and it should
error out.

```bash
curl --location --request POST 'http://localhost:8042/api/v1/events' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "foo",
    "version": "1.0.1",
    "release": "2023.11.16",
    "platform_id": "aarch64-gnu-linux-7",
    "package": "oci",
    "description": "packaged oci image foo",
    "payload": { "name" : "foo" },
    "success": true,
    "event_receiver_id": "01HFW5HXQR28951NR8NH3WJBN6"
}'
```

Error Message

```json
{
  "data": "",
  "errors": [
    "event payload did not match event receiver schema\n(root): context is required\n(root): subject is required\n(root): Additional property name is not allowed"
  ]
}
```

## Create a watcher to match CDEvent

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

// customMatcher matches a cdevent type
func customMatcher(msg *message.Message) bool {
	return msg.Type == "dev.cdevents.artifact.packaged.0.1.1"
}

func customTaskHandler(msg *message.Message) error {
	log.Default().Printf("I received a task with value '%v'", msg)
	return nil
}

```

We can now start up the watcher and start consuming messages.

```bash
go run main.go
```

You should see a log stating that we have begin consuming records.

Now we create a new event with a CDEvents payload:

Create an event:

```bash
curl --location --request POST 'http://localhost:8042/api/v1/events' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "foo",
    "version": "1.0.2",
    "release": "2023.11.20",
    "platform_id": "aarch64-gnu-linux-7",
    "package": "oci",
    "description": "packaged oci image foo",
    "payload": {
  "context": {
    "version": "0.4.0-draft",
    "id": "271069a8-fc18-44f1-b38f-9d70a1695819",
    "source": "/event/source/123",
    "type": "dev.cdevents.artifact.packaged.0.1.1",
    "timestamp": "2023-03-20T14:27:05.315384Z"
  },
  "subject": {
    "id": "pkg:golang/mygit.com/myorg/myapp@234fd47e07d1004f0aed9c",
    "source": "/event/source/123",
    "type": "artifact",
    "content": {
      "change": {
        "id": "myChange123",
        "source": "my-git.example/an-org/a-repo"
      }
    }
  }
}
    ,
    "success": true,
    "event_receiver_id": "01HFW5HXQR28951NR8NH3WJBN6"
}'
```
