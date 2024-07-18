package manager

import (
	"messaggio/infra"
)

type ServiceManager interface {
}

type serviceManager struct {
	infra infra.Infra
	repo  RepoManager
}

// NewServiceManager ...
func NewServiceManager(infra infra.Infra) ServiceManager {
	return &serviceManager{
		infra: infra,
		repo:  NewRepoManager(infra),
	}
}
