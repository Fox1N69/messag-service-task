package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
}

func NewKafkaConsumer(brokers, groupID, topic string) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: c}, nil
}

func (kc *KafkaConsumer) ConsumerMessage(processMessage func([]byte) error) {
	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}
		if err := processMessage(msg.Value); err != nil {
			log.Printf("Failed to process message: %v\n", err)
		}
	}
}

func (kc *KafkaConsumer) Close() {
	kc.consumer.Close()
}
