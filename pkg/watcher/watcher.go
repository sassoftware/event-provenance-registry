package watcher

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Record struct {
	*kgo.Record
}

type Watcher struct {
	// Ensure Client is closed to preserve proper state in partitions
	//
	// defer watcher.Client.Close()
	Client *kgo.Client

	taskChan chan *Record
}

/**
	A sample watcher can be found below using the functions in this SDK
	
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
**/
func New(brokers, topics []string, consumerGroup string) (*Watcher, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(topics...),
		kgo.ConsumerGroup(consumerGroup),
	)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		Client:   client,
		taskChan: make(chan *Record, 100),
	}, nil
}

func (w *Watcher) ConsumeRecords(matches func(record *Record) bool) {
	log.Default().Println("consuming records...")
	ctx := context.Background()
	for {
		fetches := w.Client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Fatal(fmt.Sprint(errs))
		}

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			p.EachRecord(func(r *kgo.Record) {
				record := &Record{Record: r}
				if match := matches(record); match {
					w.taskChan <- record
				}
			})
		})
	}
}

func (w *Watcher) StartTaskHandler(taskHandler func(*Record) error) {
	for {
		task := <-w.taskChan
		if task == nil {
			log.Default().Println("task is null, leaving task handler")
			return
		}

		err := taskHandler(task)
		if err != nil {
			log.Default().Println(err)
		}
	}
}
