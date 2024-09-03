package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/yay14/pulse/ingestion"
)

// KafkaConfig represents the Kafka configuration
type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

// Message represents the data structure of a message consumed from Kafka.
type Message struct {
	SourceID   string            `json:"source_id"`
	SourceType string            `json:"source_type"`
	Metrics    []*ingestion.MetricData `json:"metrics"`
}

// Consumer starts consuming messages from Kafka and processes them using the provided handler.
func Consumer(cfg KafkaConfig, handler func(message Message)) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer group: %v", err)
	}

	ctx := context.Background()

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{cfg.Topic}, &consumer{handler: handler}); err != nil {
				log.Fatalf("Error while consuming messages: %v", err)
			}
		}
	}()
}

// consumer is an implementation of the sarama.ConsumerGroupHandler interface.
type consumer struct {
	handler func(message Message)
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes each message in a consumer claim.
func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var kafkaMessage Message
		if err := json.Unmarshal(message.Value, &kafkaMessage); err != nil {
			log.Printf("Failed to unmarshal Kafka message: %v", err)
			continue
		}

		log.Printf("Consumed message: %s", string(message.Value))
		session.MarkMessage(message, "")

		// Invoke the handler function to process the message
		c.handler(kafkaMessage)
	}

	return nil
}
