package main

import (
    "log"
    "gitlab.sas.com/async-event-infrastructure/server/pkg/watcher"
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
