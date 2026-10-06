package app

import (
	"vysh-kat/internal/inspection"
	"vysh-kat/internal/novice_policy"
	"vysh-kat/internal/repository"
	"vysh-kat/internal/service"
)

// Container DI
type Container struct {
	Service *service.Service
}

func NewContainer() *Container {
	// Creating Novice Policy
	policy := &novice_policy.BasicNovicePolicy{}

	// Creating Repository
	rep := repository.NewMemoryRepository(policy)

	// Creating Service Center
	sc := &inspection.RealServiceCenter{}

	// Creating Service
	svc := service.NewService(rep, sc)
	return &Container{Service: svc}
}
