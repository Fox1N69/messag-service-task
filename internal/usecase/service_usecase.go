package usecase

import (
	"messaggio/infra"
)

type ServiceUseCase interface {
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
