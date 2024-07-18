package services

import "messaggio/internal/repository"

type MessageService interface {
}

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(messageRepository repository.MessageRepository) MessageService {
	return &messageService{
		repo: messageRepository,
	}
}

