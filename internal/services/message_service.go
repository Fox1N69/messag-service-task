package services

import (
	"context"
	"messaggio/internal/repository"
	"messaggio/storage/sqlc/database"
)

type MessageService interface {
	CreateMessage(ctx context.Context, content string, statusID int64) error
	GetStatistics(ctx context.Context) (map[string]int64, error)
	GetMessages(ctx context.Context) ([]database.Message, error)
	GetMessageByID(ctx context.Context, messageID int64) (database.Message, error)
	UpdateMessageStatus(ctx context.Context, statusID, id int64) error
	UpdateMessageProcessingDetails(
		ctx context.Context,
		id int64,
		topic string,
		partition int32,
		offset int64,
	) error
	GetMessageStatistics(ctx context.Context) (map[string]int64, error)
}

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(messageRepository repository.MessageRepository) MessageService {
	return &messageService{
		repo: messageRepository,
	}
}

func (ms *messageService) CreateMessage(ctx context.Context, content string, statusID int64) error {
	return ms.repo.Create(ctx, content, statusID)
}

func (ms *messageService) GetStatistics(ctx context.Context) (map[string]int64, error) {
	// Получаем количество обработанных сообщений
	count, err := ms.repo.GetProcessedMessagesCount(ctx)
	if err != nil {
		return nil, err
	}

	// Возвращаем статистику в виде карты
	stats := map[string]int64{
		"processed_messages": count,
	}
	return stats, nil
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

func (ms *messageService) UpdateMessageProcessingDetails(
	ctx context.Context,
	id int64,
	topic string,
	partition int32,
	offset int64,
) error {
	return ms.repo.UpdateMessageProcessingDetails(ctx, id, topic, partition, offset)
}

func (ms *messageService) GetMessageStatistics(ctx context.Context) (map[string]int64, error) {
	return ms.repo.GetMessageStatistics(ctx)
}
