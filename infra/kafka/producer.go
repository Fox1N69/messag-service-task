package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer}, nil
}

func (kp *KafkaProducer) ProduceMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := kp.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Produced message to topic %s: %s", topic, string(message))
	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}
