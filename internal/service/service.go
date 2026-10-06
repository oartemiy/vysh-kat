package service

import (
	"fmt"
	"strings"
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

func validateIdentity(name, inv string) error {
	if strings.TrimSpace(name) == "" {
		return domain.ErrEmptyName
	}
	if strings.TrimSpace(inv) == "" {
		return domain.ErrEmptyInventoryNumber
	}
	return nil
}

func (s *Service) AddTransport(t domain.Transport) error {
	if err := validateIdentity(t.Name(), t.InventoryNumber()); err != nil {
		return err
	}
	if s.repository.HasInventoryNumber(t.InventoryNumber()) {
		return fmt.Errorf("%w: %q", domain.ErrDuplicateInventory, t.InventoryNumber())
	}
	result := s.serviceCenter.Inspect(t)
	if !result.Accepted {
		return fmt.Errorf("%w: %s", domain.ErrInspectionRejected, result.Reason)
	}
	s.repository.AddTransport(t)
	return nil
}

func (s *Service) AddThing(th domain.Thing) error {
	if err := validateIdentity(th.Name(), th.InventoryNumber()); err != nil {
		return err
	}
	if s.repository.HasInventoryNumber(th.InventoryNumber()) {
		return fmt.Errorf("%w: %q", domain.ErrDuplicateInventory, th.InventoryNumber())
	}
	s.repository.AddThing(th)
	return nil
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
