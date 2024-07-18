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
}

type messageRepository struct {
	queries *database.Queries
	log     logger.Logger
}

func NewMessageRepository(sqlcQueires *database.Queries) MessageRepository {
	logger := logger.GetLogger()

	return &messageRepository{
		queries: sqlcQueires,
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
		return 0, fmt.Errorf("faile to insert message: %w", err)
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
