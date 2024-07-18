package repository

import (
	"context"
	"fmt"
	"messaggio/pkg/util/logger"
	"messaggio/storage/sqlc/database"

	"github.com/jackc/pgx/v5/pgtype"
)

type ProcessedMsgRepository interface {
}

type processedMsgRepository struct {
	queries *database.Queries
	log     logger.Logger
}

func NewProcessedMsgRepository(sqlcQueries *database.Queries) ProcessedMsgRepository {
	logger := logger.GetLogger()

	return &processedMsgRepository{
		queries: sqlcQueries,
		log:     logger,
	}
}

func (pmr *processedMsgRepository) Create(
	ctx context.Context,
	messageID int64,
	kafkaTopic string,
	kafkaPartition int32,
	kafkaOffset int64,
) (int64, error) {
	id, err := pmr.queries.InsertProcessedMessage(ctx, database.InsertProcessedMessageParams{
		MessageID:      pgtype.Int8{Int64: messageID, Valid: true},
		KafkaTopic:     pgtype.Text{String: kafkaTopic, Valid: true},
		KafkaPartition: pgtype.Int4{Int32: kafkaPartition, Valid: true},
		KafkaOffset:    pgtype.Int8{Int64: kafkaOffset, Valid: true},
	})
	if err != nil {
		pmr.log.Errorf("failed to insert processed message: %v", err)
		return 0, fmt.Errorf("failed to insert processed message: %w", err)
	}

	return id, nil
}

func (pmr *processedMsgRepository) GetProcessedMessage(ctx context.Context) ([]database.GetProcessedMessagesRow, error) {
	processed, err := pmr.queries.GetProcessedMessages(ctx)
	if err != nil {
		pmr.log.Errorf("failed to get processed message: %v", err)
		return nil, fmt.Errorf("failed to get processed message: %w", err)
	}

	return processed, nil
}
