package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

// KafkaConsumer represents a consumer for Kafka messages.
// It manages the connection to Kafka, the topic from which messages are consumed,
// and the handler that processes the messages.
type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	handler  sarama.ConsumerGroupHandler
}

// NewKafkaConsumer creates a new KafkaConsumer instance.
// It initializes a Kafka consumer group with the provided brokers, group ID, and topic.
// It also sets up the message handler for processing consumed messages.
//
// Parameters:
//   - brokers: A slice of Kafka broker addresses.
//   - groupID: The consumer group ID used for this consumer.
//   - topic: The Kafka topic from which to consume messages.
//   - handler: A handler implementing sarama.ConsumerGroupHandler to process messages.
//
// Returns:
//   - A pointer to a KafkaConsumer instance or an error if the consumer cannot be created.
func NewKafkaConsumer(brokers []string, groupID, topic string, handler sarama.ConsumerGroupHandler) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
		handler:  handler,
	}, nil
}

// ConsumeMessages starts consuming messages from the Kafka topic.
// It runs an infinite loop, processing messages using the configured handler.
// The loop will exit if the context is done.
//
// Returns:
//   - An error if there is a failure while consuming messages. If the context is done, it returns nil.
func (kc *KafkaConsumer) ConsumeMessages() error {
	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := kc.consumer.Consume(ctx, []string{kc.topic}, kc.handler)
			if err != nil {
				log.Printf("Error consuming messages: %v", err)
				return err
			}
		}
	}
}

// Close shuts down the Kafka consumer gracefully.
// It closes the connection to Kafka and releases any resources.
//
// Returns:
//   - An error if there is an issue closing the consumer.
func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
