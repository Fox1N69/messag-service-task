package handlers

import (
	"context"
	"messaggio/internal/domain/models"
	"messaggio/internal/services"
	"messaggio/pkg/http/response"

	"github.com/gofiber/fiber/v2"
)

type MessageHandler interface {
	CreateMessage(c *fiber.Ctx) error
}

type messageHandler struct {
	messageService services.MessageService
}

func NewMessageHandler(messageService services.MessageService) MessageHandler {
	return &messageHandler{
		messageService: messageService,
	}
}

func (mh *messageHandler) CreateMessage(c *fiber.Ctx) error {
	response := response.New(c)

	var req models.CreateMessageReq

	if err := c.BodyParser(&req); err != nil {
		return response.Error(400, err)
	}

	id, err := mh.messageService.CreateMessage(context.Background(), req.Content, req.StatusID)
	if err != nil {
		return response.Error(501, err)
	}

	return response.WriteMap(201, fiber.Map{
		"message": "Message create success",
		"id":      id,
	})
}
