# Watcher

## Overview

In this tutorial we will explore using the watcher sdk to create a watcher to listen for events from the EPR server.

This tutorial depends on the [Hello World](../hello_world/README.md) being completed.

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
	return string(record.Value) == "match"
}

func customTaskHandler(record *watcher.Record) error {
	log.Default().Printf("I received a task with value '%s'", record.Value)
	return nil
}

```

Save the file and run `go mod init` then `go mode tidy`.

