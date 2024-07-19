package services

import (
	"context"
	"messaggio/internal/repository"
)

type MessageService interface {
	CreateMessage(ctx context.Context, content string, statusID int64) (int64, error)
	GetStatistics(ctx context.Context) (map[string]int64, error)
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
