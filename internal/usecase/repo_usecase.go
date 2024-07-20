package usecase

import (
	"messaggio/infra"
	"messaggio/internal/repository"
	"messaggio/pkg/util/logger"
	"sync"
)

type RepoUseCase interface {
	MessageRepository() repository.MessageRepository
	ProcessedMsgRepository() repository.ProcessedMsgRepository
}

type repoUseCase struct {
	infra infra.Infra
	log   logger.Logger
}

// NewRepoUseCase creates a new instance of RepoUseCase using the provided infrastructure.
func NewRepoUseCase(infra infra.Infra) RepoUseCase {
	return &repoUseCase{
		infra: infra,
		log:   logger.GetLogger(),
	}
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

var (
	processedMsgRepoOnce sync.Once
	processedMsgRepo     repository.ProcessedMsgRepository
)

func (ruc *repoUseCase) ProcessedMsgRepository() repository.ProcessedMsgRepository {
	processedMsgRepoOnce.Do(func() {
		psqlClient, err := ruc.infra.PSQLClient()
		if err != nil {
			ruc.log.Errorf("failed to create PSQL client: %v", err)
			return
		}
		processedMsgRepo = repository.NewProcessedMsgRepository(psqlClient.Queries)
	})

	return processedMsgRepo
}
