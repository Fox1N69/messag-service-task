package manager

import (
	"messaggio/infra"
	"messaggio/internal/repository"
	"sync"
)

type RepoManager interface {
	MessageRepository() repository.MessageRepository
}

type repoManager struct {
	infra infra.Infra
}

// NewRepoManager creates a new instance of RepoManager using the provided infrastructure.
func NewRepoManager(infra infra.Infra) RepoManager {
	return &repoManager{infra: infra}
}

var (
	messageRepositoryOnce sync.Once
	messageRepository     repository.MessageRepository
)

func (rm *repoManager) MessageRepository() repository.MessageRepository {
	messageRepositoryOnce.Do(func() {
		messageRepository = repository.NewMessageRepository(rm.infra.PSQLClient().Queries)
	})

	return messageRepository
}
