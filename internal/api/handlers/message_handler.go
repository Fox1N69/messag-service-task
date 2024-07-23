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

// CreateMessage godoc
// @Summary Create a new message
// @Description Handles the HTTP POST request to create a new message.
//
//	It parses the request body, creates a new message using the message service,
//	and publishes the message content to Kafka asynchronously. Returns the ID of the created message.
//
// @Tags messages
// @Accept json
// @Produce json
// @Param request body dto.CreateMessageReq true "Create Message Request"
// @Success 201 {object} fiber.Map "Message creation success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 501 {object} response.ErrorResponse "Server Error"
// @Router /api/message/ [post]
func (mh *messageHandler) CreateMessage(c fiber.Ctx) error {
	response := response.New(c)

	var req dto.CreateMessageReq

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return response.Error(400, err)
	}

	err := mh.messageService.CreateMessage(context.Background(), req.Content, req.StatusID)
	if err != nil {
		return response.Error(501, err)
	}

	go func() {
		kafkaMessage := []byte(req.Content)
		mh.kafkaProducer.ProduceMessage("messages", kafkaMessage)
	}()

	return response.Write(201, "message create success")
}

// GetStatistics godoc
// @Summary Retrieve message statistics
// @Description Handles the HTTP GET request to retrieve statistics about messages.
//
//	It fetches the statistics from the message service and returns them in the response.
//
// @Tags messages
// @Produce json
// @Success 200 {object} fiber.Map "Message statistics"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/message/stat [get]
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
