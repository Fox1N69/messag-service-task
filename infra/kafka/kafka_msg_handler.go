package kafka

import (
	"context"
	"fmt"
	"log"
	"messaggio/internal/repository"

	"github.com/IBM/sarama"
)

type MessageHandler struct {
	processedRepo repository.ProcessedMsgRepository
	messageRepo   repository.MessageRepository
}

func NewKafkaMessageHandler(
	processedRepo repository.ProcessedMsgRepository,
	messageRepo repository.MessageRepository,
) *MessageHandler {
	return &MessageHandler{
		processedRepo: processedRepo,
		messageRepo:   messageRepo,
	}
}

func (h *MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		messageID, err := h.handleMessage(msg.Value)
		if err != nil {
			log.Printf("Error handling message: %v", err)
			continue
		}

		_, err = h.processedRepo.Create(
			context.Background(),
			messageID,
			msg.Topic,
			msg.Partition,
			msg.Offset,
		)
		if err != nil {
			log.Printf("Error marking message as processed: %v", err)
		}

		sess.MarkMessage(msg, "")
	}

	return nil
}

func (h *MessageHandler) handleMessage(message []byte) (int64, error) {
	// Декодируем сообщение и извлекаем ID
	// Это может зависеть от формата сообщений
	content := string(message)
	// Получаем сообщение по содержимому
	messages, err := h.messageRepo.GetMessages(context.Background())
	if err != nil {
		return 0, err
	}

	var messageID int64
	for _, m := range messages {
		if m.Content == content {
			messageID = m.ID
			break
		}
	}

	if messageID == 0 {
		return 0, fmt.Errorf("message not found")
	}

	// Обновляем статус сообщения на "Processed"
	err = h.messageRepo.UpdateMessageStatus(context.Background(), 3, messageID) // assuming status ID 3 is "Processed"
	if err != nil {
		return 0, err
	}

	return messageID, nil
}
