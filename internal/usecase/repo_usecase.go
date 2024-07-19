package usecase

import (
	"messaggio/infra"
	"messaggio/internal/repository"
	"sync"
)

type RepoUseCase interface {
	MessageRepository() repository.MessageRepository
	ProcessedMsgRepository() repository.ProcessedMsgRepository
}

type repoUseCase struct {
	infra infra.Infra
}

// NewRepoUseCase creates a new instance of RepoUseCase using the provided infrastructure.
func NewRepoUseCase(infra infra.Infra) RepoUseCase {
	return &repoUseCase{infra: infra}
}

var (
	messageRepoOnce sync.Once
	messageRepo     repository.MessageRepository
)

func (rm *repoUseCase) MessageRepository() repository.MessageRepository {
	messageRepoOnce.Do(func() {
		messageRepo = repository.NewMessageRepository(rm.infra.PSQLClient().Queries)
	})

	return messageRepo
}

var (
	processedMsgRepoOnce sync.Once
	processedMsgRepo     repository.ProcessedMsgRepository
)

func (ruc *repoUseCase) ProcessedMsgRepository() repository.ProcessedMsgRepository {
	processedMsgRepoOnce.Do(func() {
		processedMsgRepo = repository.NewProcessedMsgRepository(ruc.infra.PSQLClient().Queries)
	})

	return processedMsgRepo
}
