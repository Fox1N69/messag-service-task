package kafka

import (
	"context"
	"fmt"
	"messaggio/internal/repository"
	"messaggio/pkg/util/logger"

	"github.com/IBM/sarama"
)

type MessageHandler struct {
	messageRepo repository.MessageRepository
	log         logger.Logger
}

func NewKafkaMessageHandler(
	messageRepo repository.MessageRepository,
) *MessageHandler {
	return &MessageHandler{
		messageRepo: messageRepo,
		log:         logger.GetLogger(),
	}
}

func (h *MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.log.Debugf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		messageID, err := h.handleMessage(msg.Value)
		if err != nil {
			h.log.Errorf("Error handling message: %v", err)
			continue
		}

		err = h.messageRepo.UpdateMessageStatus(context.Background(), 3, messageID) // assuming status ID 3 is "Processed"
		if err != nil {
			h.log.Errorf("Error updating message status: %v", err)
		}

		err = h.messageRepo.UpdateMessageProcessingDetails(
			context.Background(),
			messageID,
			msg.Topic,
			msg.Partition,
			msg.Offset,
		)
		if err != nil {
			h.log.Errorf("Error recording message details: %v", err)
		}

		sess.MarkMessage(msg, "")
	}

	return nil
}

func (h *MessageHandler) handleMessage(message []byte) (int64, error) {
	content := string(message)
	h.log.Debugf("Handling message with content: %s", content)

	// Получаем сообщение по содержимому
	messages, err := h.messageRepo.GetMessages(context.Background())
	if err != nil {
		h.log.Errorf("Failed to get messages from repository: %v", err)
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
		h.log.Errorf("Message not found in repository for content: %s", content)
		return 0, fmt.Errorf("message not found")
	}

	// Обновляем статус сообщения на "Processed"
	err = h.messageRepo.UpdateMessageStatus(context.Background(), 3, messageID) // assuming status ID 3 is "Processed"
	if err != nil {
		h.log.Errorf("Failed to update message status: %v", err)
		return 0, err
	}

	return messageID, nil
}
