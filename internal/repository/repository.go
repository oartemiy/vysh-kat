package repository

import (
	"vysh-kat/internal/domain"
	"vysh-kat/internal/novice_policy"
)

type Repository interface {
	AddTransport(t domain.Transport)
	AddThing(th domain.Thing)

	GetThings() []domain.Thing
	GetTransports() []domain.Transport

	GetAllItems() []domain.InventoryItem

	TotalDailyEnergyKWh() float64
	NoviceSuitableTransports() []domain.Transport

	ChangeNovicePolicy(p novice_policy.NovicePolicy)
}
