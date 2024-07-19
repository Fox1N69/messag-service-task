package kafka

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.AsyncProducer
	errors   <-chan *sarama.ProducerError
	success  <-chan *sarama.ProducerMessage
}

// NewKafkaProducer создает новый асинхронный Kafka продюсер.
func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	// Инициализируем каналы для обработки сообщений и ошибок.
	errors := producer.Errors()
	success := producer.Successes()

	return &KafkaProducer{
		producer: producer,
		errors:   errors,
		success:  success,
	}, nil
}

// ProduceMessage отправляет сообщение асинхронно.
func (kp *KafkaProducer) ProduceMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	// Отправляем сообщение асинхронно
	kp.producer.Input() <- msg

	// Обработка успешной отправки и ошибок в отдельных горутинах.
	go func() {
		select {
		case <-kp.success:
			log.Printf("Produced message to topic %s: %s", topic, string(message))
		case err := <-kp.errors:
			log.Printf("Failed to produce message: %v", err)
		}
	}()

	return nil
}

// Close закрывает соединение с Kafka и завершает работу асинхронного продюсера.
func (kp *KafkaProducer) Close() error {
	// Закрытие Input канала для остановки асинхронного продюсера.
	kp.producer.AsyncClose()

	// Ожидаем, пока все сообщения будут отправлены и обработаны.
	select {
	case err := <-kp.errors:
		if err != nil {
			return err
		}
	case success := <-kp.success:
		log.Printf("Final success: %v", success)
	case <-time.After(10 * time.Second): // Timeout, если необходимо
		log.Println("Timeout while waiting for all messages to be processed")
	}

	return nil
}
