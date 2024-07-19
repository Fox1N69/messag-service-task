package handlers

import (
	"context"
	"log"
	"messaggio/infra/kafka"
	"messaggio/internal/dto"
	"messaggio/internal/services"
	"messaggio/pkg/http/response"

	"github.com/gofiber/fiber/v3"
	"github.com/golang/snappy"
)

type MessageHandler interface {
	CreateMessage(c fiber.Ctx) error
	GetStatistics(c fiber.Ctx) error
}

type messageHandler struct {
	messageService services.MessageService
	kafkaProducer  *kafka.KafkaProducer
}

func NewMessageHandler(
	messageService services.MessageService,
	kafkaProducer *kafka.KafkaProducer,
) MessageHandler {
	return &messageHandler{
		messageService: messageService,
		kafkaProducer:  kafkaProducer,
	}
}

func (mh *messageHandler) CreateMessage(c fiber.Ctx) error {
	response := response.New(c)

	var req dto.CreateMessageReq

	if err := c.Bind().Body(&req); err != nil {
		return response.Error(400, err)
	}

	id, err := mh.messageService.CreateMessage(context.Background(), req.Content, req.StatusID)
	if err != nil {
		return response.Error(501, err)
	}

	kafkaMessage := []byte(req.Content)
	compressedMessage := snappy.Encode(nil, kafkaMessage)

	// Отправка сообщения асинхронно
	go func() {
		if err := mh.kafkaProducer.ProduceMessage("messages", compressedMessage); err != nil {
			log.Printf("Failed to produce message: %v", err)
		}
	}()

	return response.WriteMap(201, fiber.Map{
		"message": "Message create success",
		"id":      id,
	})
}

func (mh *messageHandler) GetStatistics(c fiber.Ctx) error {
	response := response.New(c)

	stats, err := mh.messageService.GetStatistics(context.Background())
	if err != nil {
		return response.Error(500, err)
	}

	return response.WriteMap(200, fiber.Map{
		"statistics": stats,
	})
}
