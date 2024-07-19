package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type MessageHandler struct{}

func (h *MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
		sess.MarkMessage(msg, "")
	}
	return nil
}
