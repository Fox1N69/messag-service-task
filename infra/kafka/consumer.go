package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	handler  sarama.ConsumerGroupHandler
}

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


func (kc *KafkaConsumer) ConsumeMessages() error {
	ctx := context.Background()
	for {
		if err := kc.consumer.Consume(ctx, []string{kc.topic}, kc.handler); err != nil {
			log.Printf("Error consuming messages: %v", err)
			return err
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
