package service

import (
	"errors"
	"vysh-kat/internal/domain"
	"vysh-kat/internal/inspection"
	"vysh-kat/internal/repository"
)

type Service struct {
	rep           repository.Repository
	serviceCenter inspection.ServiceCenter
}

func (s *Service) AddTransport(t domain.Transport) error {
	result := s.serviceCenter.Inspect(t)
	if !result.Accepted {
		return errors.New(result.Reason)
	}
	s.rep.AddTransport(t)
	return nil
}

func (s *Service) AddThing(th domain.Thing) {
	s.rep.AddThing(th)
}

func (s *Service) TotalDailyEnergyKWh() float64 {
	return s.rep.TotalDailyEnergyKWh()
}

func (s *Service) NoviceSuitableTransports() []domain.Transport {
	return s.rep.NoviceSuitableTransports()
}
