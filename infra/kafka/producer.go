package kafka

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
)

// KafkaProducer represents a producer for sending messages to Kafka.
// It manages the Kafka producer, a channel for messages, and synchronization for message processing.
type KafkaProducer struct {
	producer       sarama.SyncProducer
	messageChannel chan *sarama.ProducerMessage
	waitGroup      sync.WaitGroup
}

// NewKafkaProducer creates a new KafkaProducer instance with the provided
// synchronous producer. It also starts a goroutine for processing messages.
//
// producer - A synchronous Sarama producer for sending messages to Kafka.
//
// Returns:
// *KafkaProducer - A pointer to the newly created KafkaProducer instance.
func NewKafkaProducer(producer sarama.SyncProducer) *KafkaProducer {
	kp := &KafkaProducer{
		producer:       producer,
		messageChannel: make(chan *sarama.ProducerMessage, 100),
	}

	go kp.startProducing()

	return kp
}

// startProducing processes messages from the channel and sends them to Kafka.
// This method runs in a separate goroutine and is initiated when a KafkaProducer
// is created.
func (kp *KafkaProducer) startProducing() {
	for msg := range kp.messageChannel {
		_, _, err := kp.producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message to Kafka: %v", err)
		} else {
			log.Printf("Produced message to topic %s: %s", msg.Topic, string(msg.Value.(sarama.ByteEncoder)))
		}
		kp.waitGroup.Done()
	}
}

// ProduceMessage sends a message to the channel for asynchronous delivery to Kafka.
//
// topic - The name of the Kafka topic to which the message will be sent.
// message - The message data as a byte slice.
func (kp *KafkaProducer) ProduceMessage(topic string, message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	kp.waitGroup.Add(1)
	kp.messageChannel <- msg
}

// Close shuts down the producer and waits for all messages in the channel to be processed.
// Errors during producer shutdown are logged.
func (kp *KafkaProducer) Close() {
	close(kp.messageChannel)
	kp.waitGroup.Wait()
	if err := kp.producer.Close(); err != nil {
		log.Printf("Failed to close Kafka producer: %v", err)
	} else {
		log.Println("Kafka producer closed successfully")
	}
}
