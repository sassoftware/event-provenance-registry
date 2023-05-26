// Copyright © 2019, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package message

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/utils"
)

var logger = utils.MustGetLogger("server", "message.consume")

// ConsumerController is an abstraction around consumer groups to make them easier to use. The methods used to build this
// object are exported in the event that this object is insufficient for the task.
type ConsumerController struct {
	consumer      Consumer
	consumerGroup sarama.ConsumerGroup
	topics        []string
}

// NewConsumerController creates a new ConsumerController
func NewConsumerController(kafkaVersion, groupID string, kafkaPeers, topics []string, worker ConsumerWorker) (*ConsumerController, error) {
	consumer := newConsumer(worker)

	group, err := CreateConsumerGroup(kafkaVersion, groupID, kafkaPeers)
	if err != nil {
		return nil, err
	}

	return &ConsumerController{
		consumer:      consumer,
		consumerGroup: group,
		topics:        topics,
	}, nil
}

// NewSecureConsumerController creates a new ConsumerController with security flags invoked.
func NewSecureConsumerController(kafkaVersion, groupID, saslUser, saslPass string, tls bool, kafkaPeers, topics []string, worker ConsumerWorker) (*ConsumerController, error) {
	consumer := newConsumer(worker)

	group, err := CreateSecureConsumerGroup(kafkaVersion, groupID, saslUser, saslPass, tls, kafkaPeers)
	if err != nil {
		return nil, err
	}

	return &ConsumerController{
		consumer:      consumer,
		consumerGroup: group,
		topics:        topics,
	}, nil
}

// NewConsumerControllerEnv creates a new ConsumerController with security flags invoked based on SASL_USERNAME and SASL_PASSWORD
// envirionment variables. If SASL_USERNAME and SASL_PASSWORD are set, a ConsumerController with SASL enabled will be returned.
func NewConsumerControllerEnv(kafkaVersion, groupID string, tls bool, kafkaPeers, topics []string, worker ConsumerWorker) (*ConsumerController, error) {
	consumer := newConsumer(worker)

	group, err := CreateConsumerGroupEnv(kafkaVersion, groupID, tls, kafkaPeers)
	if err != nil {
		return nil, err
	}

	return &ConsumerController{
		consumer:      consumer,
		consumerGroup: group,
		topics:        topics,
	}, nil
}

// BeginConsuming starts consuming messages off of the kafka bus inside of a goroutine.
func (c *ConsumerController) BeginConsuming(ctx context.Context, wg *sync.WaitGroup) {
	ConsumeMessages(ctx, c.topics, c.consumerGroup, c.consumer, wg)
	logger.V(1).Info("waiting for consumer to be ready")
	<-c.consumer.Ready
	logger.V(1).Info("consumer ready")
}

// Close closes the ConsumerGroup
func (c *ConsumerController) Close() error {
	return c.consumerGroup.Close()
}

// Consumer implements the sarama ConsumerGroupHandler interface.
type Consumer struct {
	Ready  chan bool
	Worker ConsumerWorker
}

func newConsumer(worker ConsumerWorker) Consumer {
	return Consumer{
		Ready:  make(chan bool),
		Worker: worker,
	}
}

// ConsumerWorker represents a function that operates on sarama.ConsumerMessage structs. In practice, this means that
// this function should be executed on each messages received from the bus. Refer to ConsumeClaim as an example.
type ConsumerWorker func(*sarama.ConsumerMessage) error

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages(). Each message received is handled by
// the Worker function on the consumer. Messages that process with an error are marked as read and will not be reprocessed.
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		err := c.Worker(message)
		session.MarkMessage(message, "")
		if err != nil {
			logger.Error(err, fmt.Sprintf("message error : %v", err))
			return err
		}
	}

	return nil
}

// ConsumeMessages given a ConsumerGroup and a Consumer, process each message in the consumer group using the consumer
// provided. You will need to close the consumer group when you are done processing.
func ConsumeMessages(ctx context.Context, topics []string, group sarama.ConsumerGroup, cons Consumer, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := group.Consume(ctx, topics, &cons)
				if err != nil {
					logger.Error(err, fmt.Sprintf("error from consumer: %v, will try again", err))
				}
				// check if context was cancelled, signaling that the consumer should stop
				if ctx.Err() != nil {
					return
				}
				cons.Ready = make(chan bool)
			}
		}
	}()
}

// CreateConsumerGroup returns a new ConsumerGroup, ready to consume things. Calls through to CreateSecureConsumerGroup
// without any security pararmeters enabled.
func CreateConsumerGroup(kafkaVersion, groupID string, kafkaPeers []string) (sarama.ConsumerGroup, error) {
	return CreateSecureConsumerGroup(kafkaVersion, groupID, "", "", false, kafkaPeers)
}

// CreateSecureConsumerGroup creates a ConsumerGroup with Kafka security options enabled. Providing saslUser and saslPassword
// allows the consumer to authenticate with a SASL enabled kafka cluster. If left empty, a normal consumer will be created.
// The previous two options imply TLS. TLS can be set separately in cases where TLS is required but auth is not.
func CreateSecureConsumerGroup(kafkaVersion, groupID, saslUser, saslPass string, tls bool, kafkaPeers []string) (sarama.ConsumerGroup, error) {
	saramaCfg, err := NewConfig(kafkaVersion)
	if err != nil {
		return nil, err
	}

	if tls {
		logger.Info("Enabled Kafka TLS")
		saramaCfg.Net.TLS.Enable = true
	}

	if groupID == "" {
		groupID = "re.polaris.watcher"
	}
	saramaCfg.ClientID = groupID + "." + utils.NewULIDAsString()

	if saslUser != "" && saslPass != "" {
		logger.Info("Enabled Kafka SASL")
		saramaCfg, err = NewSCRAMConfig(saslUser, saslPass, kafkaVersion)
		if err != nil {
			return nil, err
		}
	}

	consGroup, err := sarama.NewConsumerGroup(kafkaPeers, groupID, saramaCfg)
	if err != nil {
		return nil, err
	}

	return consGroup, nil
}

// CreateConsumerGroupEnv returns a new ConsumerGroup with Kafka security options enabled or disabled depending on
// configured envirionment variables. If SASL_USERNAME and SASL_PASSWORD are set, a consumer with SASL enabled will be
// returned, else a normal consumer group will be returned. If SASL is enabled SASL_MECHANISM must be set to SCRAM or PLAIN
// Supported SASL_MECHANISMS: SCRAM, PLAIN
// Future    SASL_MECHANISMS: OAUTH2
func CreateConsumerGroupEnv(kafkaVersion, groupID string, tls bool, kafkaPeers []string) (sarama.ConsumerGroup, error) {
	saramaCfg, err := NewConfigEnv(kafkaVersion)
	if err != nil {
		return nil, err
	}

	if tls {
		logger.Info("Enabled Kafka TLS")
		saramaCfg.Net.TLS.Enable = true
	}

	if groupID == "" {
		return nil, fmt.Errorf("Consumer Group groupID cannot be empty")
	}
	saramaCfg.ClientID = groupID + "." + utils.NewULIDAsString()

	consGroup, err := sarama.NewConsumerGroup(kafkaPeers, groupID, saramaCfg)
	if err != nil {
		return nil, err
	}
	return consGroup, nil
}
