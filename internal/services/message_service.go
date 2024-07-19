package services

import (
	"context"
	"messaggio/internal/repository"
	"messaggio/storage/sqlc/database"
)

type MessageService interface {
	CreateMessage(ctx context.Context, content string, statusID int64) (int64, error)
	GetStatistics(ctx context.Context) (map[string]int64, error)
	GetMessages(ctx context.Context) ([]database.Message, error)
	GetMessageByID(ctx context.Context, messageID int64) (database.Message, error)
	UpdateMessageStatus(ctx context.Context, statusID, id int64) error
	CreateProcessed(
		ctx context.Context,
		messageID int64,
		kafkaTopic string,
		kafkaPartition int32,
		kafkaOffset int64,
	) (int64, error)
	MessageByContent(ctx context.Context, content string) (database.Message, error)
}

type messageService struct {
	repo          repository.MessageRepository
	processedRepo repository.ProcessedMsgRepository
}

func NewMessageService(
	messageRepository repository.MessageRepository,
	processedRepository repository.ProcessedMsgRepository,
) MessageService {
	return &messageService{
		repo:          messageRepository,
		processedRepo: processedRepository,
	}
}

func (ms *messageService) CreateMessage(ctx context.Context, content string, statusID int64) (int64, error) {
	return ms.repo.Create(ctx, content, statusID)
}

func (ms *messageService) GetStatistics(ctx context.Context) (map[string]int64, error) {
	return ms.processedRepo.GetMessageStatistics(ctx)
}

func (ms *messageService) GetMessages(ctx context.Context) ([]database.Message, error) {
	return ms.repo.GetMessages(ctx)
}

func (ms *messageService) GetMessageByID(ctx context.Context, messageID int64) (database.Message, error) {
	return ms.repo.MessageByID(ctx, messageID)
}

func (ms *messageService) UpdateMessageStatus(ctx context.Context, statusID, id int64) error {
	return ms.repo.UpdateMessageStatus(ctx, statusID, id)
}

func (ms *messageService) CreateProcessed(
	ctx context.Context,
	messageID int64,
	kafkaTopic string,
	kafkaPartition int32,
	kafkaOffset int64,
) (int64, error) {
	return ms.processedRepo.Create(ctx, messageID, kafkaTopic, kafkaPartition, kafkaOffset)
}

func (ms *messageService) MessageByContent(ctx context.Context, content string) (database.Message, error) {
	return ms.repo.MessageByContent(ctx, content)
}
