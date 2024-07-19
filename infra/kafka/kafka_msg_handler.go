package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"messaggio/internal/repository"
	"messaggio/storage/sqlc/database"

	"github.com/IBM/sarama"
)

type MessageHandler struct {
	processedRepo repository.ProcessedMsgRepository
}

func NewKafkaMessageHandler(processedRepo repository.ProcessedMsgRepository) *MessageHandler {
	return &MessageHandler{processedRepo: processedRepo}
}

func (h *MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		messageID, err := extractMessageID(msg.Value)
		if err != nil {
			log.Printf("Error extracting message ID: %v", err)
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

func extractMessageID(message []byte) (int64, error) {
	var msg database.Message

	if err := json.Unmarshal(message, &msg); err != nil {
		return 0, errors.New("failed to unmarshal message")
	}

	return msg.ID, nil
}
