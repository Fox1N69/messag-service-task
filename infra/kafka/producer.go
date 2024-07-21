package kafka

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer       sarama.SyncProducer
	messageChannel chan *sarama.ProducerMessage
	waitGroup      sync.WaitGroup
}

// NewKafkaProducer создает новый KafkaProducer с заданным синхронным продюсером.
func NewKafkaProducer(producer sarama.SyncProducer) *KafkaProducer {
	kp := &KafkaProducer{
		producer:       producer,
		messageChannel: make(chan *sarama.ProducerMessage, 100), // Буферизующий канал с размером 100
	}

	// Запускаем горутину для обработки сообщений
	go kp.startProducing()

	return kp
}

// startProducing обрабатывает сообщения из канала и отправляет их в Kafka.
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

// ProduceMessage отправляет сообщение в канал для последующей асинхронной отправки в Kafka.
func (kp *KafkaProducer) ProduceMessage(topic string, message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	kp.waitGroup.Add(1)
	kp.messageChannel <- msg
}

// Close закрывает продюсер и ожидает завершения обработки всех сообщений в канале.
func (kp *KafkaProducer) Close() {
	close(kp.messageChannel)
	kp.waitGroup.Wait()
	if err := kp.producer.Close(); err != nil {
		log.Printf("Failed to close Kafka producer: %v", err)
	} else {
		log.Println("Kafka producer closed successfully")
	}
}
