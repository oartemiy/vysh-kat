package repository

import (
	"vysh-kat/internal/domain"
	"vysh-kat/internal/novice_policy"
)

type MemoryRepository struct {
	transports []domain.Transport
	things     []domain.Thing
	policy     novice_policy.NovicePolicy
}

func NewMemoryRepository(policy novice_policy.NovicePolicy) *MemoryRepository {
	return &MemoryRepository{policy: policy}
}

func (r *MemoryRepository) ChangeNovicePolicy(policy novice_policy.NovicePolicy) {
	r.policy = policy
}

func (r *MemoryRepository) AddTransport(t domain.Transport) {
	r.transports = append(r.transports, t)
}

func (r *MemoryRepository) AddThing(th domain.Thing) {
	r.things = append(r.things, th)
}

func (r *MemoryRepository) GetItems() []domain.Thing {
	return r.things
}

func (r *MemoryRepository) GetTransports() []domain.Transport {
	return r.transports
}

func (r *MemoryRepository) TotalDailyEnergyKWh() float64 {
	total := 0.0
	for _, t := range r.transports {
		total += t.DailyEnergyKWh()
	}

	for _, th := range r.things {
		if ec, ok := th.(domain.EnergyConsumer); ok {
			total += ec.DailyEnergyKWh()
		}
	}
	return total
}

func (r *MemoryRepository) NoviceSuitableTransports() []domain.Transport {
	var result []domain.Transport
	for _, t := range r.transports {
		if t.GetSimplicity() >= 6 {
			result = append(result, t)
		}
	}
	return result
}
