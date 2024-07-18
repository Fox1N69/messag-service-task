package usecase

import (
	"messaggio/infra"
	"messaggio/internal/repository"
	"sync"
)

type RepoUseCase interface {
	MessageRepository() repository.MessageRepository
}

type repoUseCase struct {
	infra infra.Infra
}

// NewRepoUseCase creates a new instance of RepoUseCase using the provided infrastructure.
func NewRepoUseCase(infra infra.Infra) RepoUseCase {
	return &repoUseCase{infra: infra}
}

var (
	messageRepositoryOnce sync.Once
	messageRepository     repository.MessageRepository
)

func (rm *repoUseCase) MessageRepository() repository.MessageRepository {
	messageRepositoryOnce.Do(func() {
		messageRepository = repository.NewMessageRepository(rm.infra.PSQLClient().Queries)
	})

	return messageRepository
}
