package kafka

import (
	"context"
	"fmt"
	"log"
	"messaggio/internal/repository"

	"github.com/IBM/sarama"
	"github.com/golang/snappy"
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

		// Декодирование сообщения Snappy
		decodedMessage, err := snappy.Decode(nil, msg.Value)
		if err != nil {
			log.Printf("Error decoding Snappy message: %v", err)
			continue
		}

		messageID, err := h.handleMessage(decodedMessage)
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
	content := string(message)
	messageID, err := h.findMessageID(content)
	if err != nil {
		return 0, fmt.Errorf("failed to find message ID: %w", err)
	}

	// Обновляем статус сообщения на "Processed"
	err = h.messageRepo.UpdateMessageStatus(context.Background(), 3, messageID)
	if err != nil {
		return 0, fmt.Errorf("failed to update message status: %w", err)
	}

	return messageID, nil
}

func (h *MessageHandler) findMessageID(content string) (int64, error) {
	message, err := h.messageRepo.MessageByContent(context.Background(), content)
	if err != nil {
		return 0, fmt.Errorf("failed to find message by content: %w", err)
	}

	if message.ID == 0 {
		return 0, fmt.Errorf("message not found")
	}

	return message.ID, nil
}
