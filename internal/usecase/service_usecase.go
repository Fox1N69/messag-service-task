package usecase

import (
	"messaggio/infra"
	"messaggio/internal/services"
	"sync"
)

type ServiceUseCase interface {
	MessageService() services.MessageService
}

type serviceUseCase struct {
	infra infra.Infra
	repo  RepoUseCase
}

// NewServiceUseCase ...
func NewServiceUseCase(infra infra.Infra) ServiceUseCase {
	return &serviceUseCase{
		infra: infra,
		repo:  NewRepoUseCase(infra),
	}
}

var (
	messageServiceOnce sync.Once
	messageService     services.MessageService
)

func (suc *serviceUseCase) MessageService() services.MessageService {
	messageServiceOnce.Do(func() {
		repo := suc.repo.MessageRepository()
		messageService = services.NewMessageService(repo)
	})

	return messageService
}
