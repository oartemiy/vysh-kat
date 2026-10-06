package mocks

import (
	"vysh-kat/internal/domain"
	"vysh-kat/internal/inspection"
	"vysh-kat/internal/novice_policy"
)

// MockRepository — in-memory mock репозитория, записывает добавленные сущности.
type MockRepository struct {
	Transports []domain.Transport
	Things     []domain.Thing

	TotalEnergy float64
	NoviceList  []domain.Transport
	ExtraItems  []domain.InventoryItem

	ChangePolicyCalls int
	LastPolicy        novice_policy.NovicePolicy
}

func (m *MockRepository) AddTransport(t domain.Transport) {
	m.Transports = append(m.Transports, t)
}

func (m *MockRepository) AddThing(th domain.Thing) {
	m.Things = append(m.Things, th)
}

func (m *MockRepository) HasInventoryNumber(inv string) bool {
	for _, t := range m.Transports {
		if t.InventoryNumber() == inv {
			return true
		}
	}
	for _, th := range m.Things {
		if th.InventoryNumber() == inv {
			return true
		}
	}
	return false
}

func (m *MockRepository) GetThings() []domain.Thing {
	return m.Things
}

func (m *MockRepository) GetTransports() []domain.Transport {
	return m.Transports
}

func (m *MockRepository) GetAllItems() []domain.InventoryItem {
	items := make([]domain.InventoryItem, 0, len(m.Transports)+len(m.Things)+len(m.ExtraItems))
	for _, t := range m.Transports {
		items = append(items, t)
	}
	for _, th := range m.Things {
		items = append(items, th)
	}
	return append(items, m.ExtraItems...)
}

func (m *MockRepository) TotalDailyEnergyKWh() float64 {
	return m.TotalEnergy
}

func (m *MockRepository) NoviceSuitableTransports() []domain.Transport {
	return m.NoviceList
}

func (m *MockRepository) ChangeNovicePolicy(p novice_policy.NovicePolicy) {
	m.ChangePolicyCalls++
	m.LastPolicy = p
}

// MockServiceCenter — сервисный центр с настраиваемым вердиктом, считает вызовы Inspect.
type MockServiceCenter struct {
	Result       inspection.ResultOfInspection
	InspectCalls int
	Inspected    []domain.Transport
}

func (m *MockServiceCenter) Inspect(t domain.Transport) inspection.ResultOfInspection {
	m.InspectCalls++
	m.Inspected = append(m.Inspected, t)
	return m.Result
}

// MockNovicePolicy — политика новичка с произвольным правилом, считает вызовы.
type MockNovicePolicy struct {
	Fn    func(t domain.Transport) bool
	Calls int
}

func (m *MockNovicePolicy) IsNoviceFriendly(t domain.Transport) bool {
	m.Calls++
	if m.Fn == nil {
		return true
	}
	return m.Fn(t)
}
