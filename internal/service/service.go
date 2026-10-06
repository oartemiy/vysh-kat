package service

import (
	"errors"
	"vysh-kat/internal/domain"
	"vysh-kat/internal/inspection"
	"vysh-kat/internal/repository"
)

type Service struct {
	repository    repository.Repository
	serviceCenter inspection.ServiceCenter
}

func NewService(rep repository.Repository, servCenter inspection.ServiceCenter) *Service {
	return &Service{repository: rep, serviceCenter: servCenter}
}

func (s *Service) AddTransport(t domain.Transport) error {
	result := s.serviceCenter.Inspect(t)
	if !result.Accepted {
		return errors.New(result.Reason)
	}
	s.repository.AddTransport(t)
	return nil
}

func (s *Service) AddThing(th domain.Thing) {
	s.repository.AddThing(th)
}

func (s *Service) TotalDailyEnergyKWh() float64 {
	return s.repository.TotalDailyEnergyKWh()
}

func (s *Service) NoviceSuitableTransports() []domain.Transport {
	return s.repository.NoviceSuitableTransports()
}

func (s *Service) AllItems() []domain.InventoryItem {
	return s.repository.GetAllItems()
}
