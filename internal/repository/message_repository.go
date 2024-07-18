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
	statusIDPg := pgtype.Int8{Int64: statusID, Valid: true}

	id, err := mr.queries.InsertMessage(ctx, database.InsertMessageParams{
		Content:  content,
		StatusID: statusIDPg,
	})
	if err != nil {
		mr.log.Errorf("failed to insert message: %v", err)
		return 0, fmt.Errorf("faile to insert message: %v", err)
	}

	return id, nil
}

