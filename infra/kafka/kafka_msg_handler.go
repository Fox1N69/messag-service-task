package kafka

import (
	"context"
	"fmt"
	"messaggio/internal/repository"
	"messaggio/pkg/util/logger"

	"github.com/IBM/sarama"
)

// MessageHandler handles the consumption of Kafka messages and their processing.
// It interacts with a message repository to update message statuses and details.
type MessageHandler struct {
	messageRepo repository.MessageRepository
	log         logger.Logger
}

// NewKafkaMessageHandler creates a new instance of MessageHandler.
//
// messageRepo - A repository for performing operations on messages, such as updating their status.
//
// Returns:
// *MessageHandler - A pointer to the newly created MessageHandler instance.
func NewKafkaMessageHandler(
	messageRepo repository.MessageRepository,
) *MessageHandler {
	return &MessageHandler{
		messageRepo: messageRepo,
		log:         logger.GetLogger(),
	}
}

// Setup is a no-op method required by the sarama.ConsumerGroupHandler interface.
// It is called at the beginning of a new session and can be used to initialize resources if needed.
//
// Returns:
// error - Always returns nil in this implementation.
func (h *MessageHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup is a no-op method required by the sarama.ConsumerGroupHandler interface.
// It is called at the end of a session and can be used to clean up resources if needed.
//
// Returns:
// error - Always returns nil in this implementation.
func (h *MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from a Kafka topic claim. It handles each message by:
// 1. Logging its content.
// 2. Processing the message and updating its status in the repository.
// 3. Recording message processing details.
// 4. Marking the message as processed.
//
// sess - The consumer group session associated with the message claim.
// claim - The claim from which messages are consumed.
//
// Returns:
// error - An error if any issues occur during processing; otherwise, returns nil.
func (h *MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	statusUpdater := NewStatusUpdater(h.messageRepo)
	defer statusUpdater.Close()

	for msg := range claim.Messages() {
		h.log.Debugf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		messageID, err := h.handleMessage(msg.Value)
		if err != nil {
			h.log.Errorf("Error handling message: %v", err)
			continue
		}

		// Отправляем обновление статуса в канал
		statusUpdater.UpdateStatus(messageID, 2) // Статус "в процессе"

		// Записываем детали обработки сообщения
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

		// Отмечаем сообщение как обработанное в Kafka
		sess.MarkMessage(msg, "")

		// После успешной обработки, обновляем статус на "обработано"
		statusUpdater.UpdateStatus(messageID, 3) // Статус "обработано"
	}

	return nil
}

// handleMessage processes the message content to find its corresponding ID in the repository.
// It updates the message status to "Processed" if found.
//
// message - The byte slice containing the message content.
//
// Returns:
// int64 - The ID of the message if found; otherwise, returns 0.
// error - An error if the message cannot be found or if there are issues updating the status; otherwise, returns nil.
func (h *MessageHandler) handleMessage(message []byte) (int64, error) {
	content := string(message)
	h.log.Debugf("Handling message with content: %s", content)

	// Getting message by content
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

	// Update messages status on "Processed"
	err = h.messageRepo.UpdateMessageStatus(context.Background(), 3, messageID)
	if err != nil {
		h.log.Errorf("Failed to update message status: %v", err)
		return 0, err
	}

	return messageID, nil
}
