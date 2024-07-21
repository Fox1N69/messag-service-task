package handlers

import (
	"context"
	"messaggio/infra/kafka"
	"messaggio/internal/dto"
	"messaggio/internal/services"
	"messaggio/pkg/http/response"

	"github.com/gofiber/fiber/v3"
	jsoniter "github.com/json-iterator/go"
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

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (mh *messageHandler) CreateMessage(c fiber.Ctx) error {
	response := response.New(c)

	var req dto.CreateMessageReq

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return response.Error(400, err)
	}

	id, err := mh.messageService.CreateMessage(context.Background(), req.Content, req.StatusID)
	if err != nil {
		return response.Error(501, err)
	}

	go func() {
		kafkaMessage := []byte(req.Content)
		mh.kafkaProducer.ProduceMessage("messages", kafkaMessage)
	}()

	return response.WriteMap(201, fiber.Map{
		"message": "Message create success",
		"id":      id,
	})
}

func (mh *messageHandler) GetStatistics(c fiber.Ctx) error {
	response := response.New(c)

	stats, err := mh.messageService.GetMessageStatistics(context.Background())
	if err != nil {
		return response.Error(500, err)
	}

	return response.WriteMap(200, fiber.Map{
		"statistics": stats,
	})
}
