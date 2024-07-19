package services

import (
	"context"
	"messaggio/internal/repository"
)

type MessageService interface {
	CreateMessage(ctx context.Context, content string, statusID int64) (int64, error)
}

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(messageRepository repository.MessageRepository) MessageService {
	return &messageService{
		repo: messageRepository,
	}
}

func (ms *messageService) CreateMessage(ctx context.Context, content string, statusID int64) (int64, error) {
	return ms.repo.Create(ctx, content, statusID)
}
