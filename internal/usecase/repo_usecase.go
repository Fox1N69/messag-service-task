package usecase

import (
	"messaggio/infra"
	"messaggio/internal/repository"
	"messaggio/pkg/util/logger"
	"sync"
)

type RepoUseCase interface {
	MessageRepository() repository.MessageRepository
}

type repoUseCase struct {
	infra infra.Infra
	log   logger.Logger
}

// NewRepoUseCase creates a new instance of RepoUseCase using the provided infrastructure.
func NewRepoUseCase(infra infra.Infra) RepoUseCase {
	logger := logger.GetLogger()
	return &repoUseCase{infra: infra, log: logger}
}

var (
	messageRepoOnce sync.Once
	messageRepo     repository.MessageRepository
)

func (rm *repoUseCase) MessageRepository() repository.MessageRepository {
	messageRepoOnce.Do(func() {
		psqlClient, err := rm.infra.PSQLClient()
		if err != nil {
			rm.log.Errorf("failed to create PSQL client: %v", err)
			return
		}
		messageRepo = repository.NewMessageRepository(psqlClient.Queries)
	})

	return messageRepo
}

