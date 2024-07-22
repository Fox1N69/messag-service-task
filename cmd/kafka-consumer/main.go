package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"messaggio/infra"
	"messaggio/infra/kafka"
	"messaggio/internal/repository"
)

func main() {
	i := infra.New("config/config.json")

	// Конфигурация подключения
	brokers := []string{"localhost:9092"}
	groupID := "example-group"
	topic := "messages"

	// Подключение к базе данных
	psql, err := i.PSQLClient()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	messageRepo := repository.NewMessageRepository(psql.Queries)
	handler := kafka.NewKafkaMessageHandler(messageRepo)

	consumer, err := kafka.NewKafkaConsumer(brokers, groupID, topic, handler)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

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
			cancel() // Завершение работы при ошибке
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
