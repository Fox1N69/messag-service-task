package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"messaggio/infra"
	"messaggio/infra/kafka"
	"messaggio/internal/repository"
	"messaggio/pkg/util/logger"
)

func main() {
	i := infra.New("config/config.json")
	mode := i.SetMode()
	logger.Init(mode)
	log := i.GetLogger().WithField("op", "consumer.client")

	// Конфигурация подключения
	kafkaConfig := i.Config().Sub("kafka")

	brokers := []string{"kafka:9092"}
	groupID := kafkaConfig.GetString("group_id")
	topic := kafkaConfig.GetString("topic")

	// Подключение к базе данных
	psql, err := i.PSQLClient()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	log.Info("Client connet to database")

	messageRepo := repository.NewMessageRepository(psql.Queries)
	handler := kafka.NewKafkaMessageHandler(messageRepo)

	consumer, err := kafka.NewKafkaConsumer(brokers, groupID, topic, handler)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()
	log.Info("Create Kafka consumer")

	// Канал для завершения работы и канал для сигналов
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Создание контекста для управления завершением работы
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск потребления сообщений
	go func() {
		if err := consumer.ConsumeMessages(); err != nil {
			log.Printf("Error consuming messages: %v", err)
			cancel()
		}
	}()

	// Ожидание сигнала завершения
	select {
	case <-stopChan:
		log.Println("Received shutdown signal")
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	// Добавляем небольшой тайм-аут перед завершением работы
	timeout := time.After(5 * time.Second)
	select {
	case <-timeout:
		log.Println("Shutdown complete")
	case <-ctx.Done():
		log.Println("Shutdown interrupted")
	}
}
