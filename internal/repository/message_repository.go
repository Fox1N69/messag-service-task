package repository

import (
	"context"
	"fmt"
	"messaggio/pkg/util/logger"
	"messaggio/storage/sqlc/database"

	"github.com/jackc/pgx/v5/pgtype"
)

type MessageRepository interface {
	Create(ctx context.Context, content string, statusID int64) (int64, error)
	GetMessages(ctx context.Context) ([]database.Message, error)
	MessageByID(ctx context.Context, messageID int64) (database.Message, error)
	UpdateMessageStatus(ctx context.Context, statusID, id int64) error
	GetProcessedMessages(ctx context.Context) ([]database.Message, error)
	GetProcessedMessagesCount(ctx context.Context) (int64, error)
	UpdateMessageProcessingDetails(
		ctx context.Context,
		id int64,
		topic string,
		partition int32,
		offset int64,
	) error
	GetMessageStatistics(ctx context.Context) (map[string]int64, error)
}

type messageRepository struct {
	queries *database.Queries
	log     logger.Logger
}

func NewMessageRepository(sqlcQueries *database.Queries) MessageRepository {
	logger := logger.GetLogger()

	return &messageRepository{
		queries: sqlcQueries,
		log:     logger,
	}
}

func (mr *messageRepository) Create(ctx context.Context, content string, statusID int64) (int64, error) {
	id, err := mr.queries.InsertMessage(ctx, database.InsertMessageParams{
		Content:  content,
		StatusID: pgtype.Int8{Int64: statusID, Valid: true},
	})
	if err != nil {
		mr.log.Errorf("failed to insert message: %v", err)
		return 0, fmt.Errorf("failed to insert message: %w", err)
	}

	return id, nil
}

func (mr *messageRepository) GetMessages(ctx context.Context) ([]database.Message, error) {
	messages, err := mr.queries.GetMessages(ctx)
	if err != nil {
		mr.log.Errorf("failed to get messages: %v", err)
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}

func (mr *messageRepository) MessageByID(ctx context.Context, messageID int64) (database.Message, error) {
	message, err := mr.queries.GetMessageByID(ctx, messageID)
	if err != nil {
		mr.log.Errorf("failed to get message by id: %v", err)
		return database.Message{}, fmt.Errorf("failed to get message by id: %w", err)
	}

	return message, nil
}

func (mr *messageRepository) UpdateMessageStatus(ctx context.Context, statusID, id int64) error {
	if err := mr.queries.UpdateMessageStatus(ctx, database.UpdateMessageStatusParams{
		ID:       id,
		StatusID: pgtype.Int8{Int64: statusID, Valid: true},
	}); err != nil {
		mr.log.Errorf("failed to update message status: %v", err)
		return fmt.Errorf("failed to update message status: %w", err)
	}

	return nil
}

func (mr *messageRepository) GetProcessedMessages(ctx context.Context) ([]database.Message, error) {
	messages, err := mr.queries.GetProcessedMessages(ctx)
	if err != nil {
		mr.log.Errorf("failed to get processed messages: %v", err)
		return nil, fmt.Errorf("failed to get processed messages: %w", err)
	}

	return messages, nil
}

func (mr *messageRepository) GetMessageStatistics(ctx context.Context) (map[string]int64, error) {
	stats, err := mr.queries.GetMessageStatistics(ctx)
	if err != nil {
		mr.log.Errorf("failed to get message statistics: %v", err)
		return nil, fmt.Errorf("failed to get message statistics: %w", err)
	}

	return map[string]int64{
		"processed":  stats.ProcessedCount,
		"received":   stats.ReceivedCount,
		"processing": stats.ProcessingCount,
	}, nil
}

func (mr *messageRepository) GetProcessedMessagesCount(ctx context.Context) (int64, error) {
	count, err := mr.queries.GetProcessedMessagesCount(ctx)
	if err != nil {
		mr.log.Errorf("failed to get processed messages count: %v", err)
		return 0, fmt.Errorf("failed to get processed messages count: %w", err)
	}

	return count, nil
}

func (mr *messageRepository) UpdateMessageProcessingDetails(
	ctx context.Context,
	id int64,
	topic string,
	partition int32,
	offset int64,
) error {
	if err := mr.queries.UpdateMessageProcessingDetails(ctx, database.UpdateMessageProcessingDetailsParams{
		ID:             id,
		KafkaTopic:     pgtype.Text{String: topic, Valid: true},
		KafkaPartition: pgtype.Int4{Int32: partition, Valid: true},
		KafkaOffset:    pgtype.Int8{Int64: offset, Valid: true},
	}); err != nil {
		mr.log.Errorf("failed to update message processing details: %v", err)
		return fmt.Errorf("failed to update message processing details: %w", err)
	}

	return nil
}
