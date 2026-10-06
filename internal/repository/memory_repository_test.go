package repository_test

import (
	"testing"

	"vysh-kat/internal/domain"
	"vysh-kat/internal/mocks"
	"vysh-kat/internal/repository"
)

type energyBox struct {
	name string
	inv  string
	kwh  float64
}

func (e energyBox) Name() string            { return e.name }
func (e energyBox) InventoryNumber() string { return e.inv }
func (e energyBox) DailyEnergyKWh() float64 { return e.kwh }

func TestNoviceSuitableTransports_UsesProvidedPolicy(t *testing.T) {
	policy := &mocks.MockNovicePolicy{Fn: func(tr domain.Transport) bool {
		return tr.GetSimplicity() >= 3
	}}
	rep := repository.NewMemoryRepository(policy)

	rep.AddTransport(domain.NewBicycle("Сложный", "INV-1", 1, 2))
	rep.AddTransport(domain.NewBicycle("Простой", "INV-2", 1, 4))

	got := rep.NoviceSuitableTransports()
	if len(got) != 1 || got[0].Name() != "Простой" {
		t.Fatalf("политика должна отбирать по своему правилу, получено %d: %v", len(got), got)
	}
	if policy.Calls != 2 {
		t.Fatalf("политика должна вызываться для каждого транспорта, вызовов: %d", policy.Calls)
	}
}

func TestNewMemoryRepository_NilPolicyUsesDefault(t *testing.T) {
	rep := repository.NewMemoryRepository(nil)

	rep.AddTransport(domain.NewBicycle("Граница", "INV-1", 1, 5))
	rep.AddTransport(domain.NewBicycle("Подходит", "INV-2", 1, 6))

	got := rep.NoviceSuitableTransports()
	if len(got) != 1 || got[0].Name() != "Подходит" {
		t.Fatalf("ожидался дефолтный порог 6, получено %d: %v", len(got), got)
	}
}

func TestChangeNovicePolicy(t *testing.T) {
	rep := repository.NewMemoryRepository(nil)
	rep.AddTransport(domain.NewBicycle("Самокат", "INV-1", 1, 8))

	if got := rep.NoviceSuitableTransports(); len(got) != 1 {
		t.Fatalf("до смены политики ожидался 1 транспорт, получено %d", len(got))
	}

	strict := &mocks.MockNovicePolicy{Fn: func(domain.Transport) bool { return false }}
	rep.ChangeNovicePolicy(strict)

	if got := rep.NoviceSuitableTransports(); len(got) != 0 {
		t.Fatalf("после смены политики ожидался пустой список, получено %d", len(got))
	}
	if strict.Calls != 1 {
		t.Fatalf("новая политика должна использоваться, вызовов: %d", strict.Calls)
	}

	rep.ChangeNovicePolicy(nil)
	if strict.Calls != 1 {
		t.Fatalf("nil-политика не должна применяться")
	}
}

func TestHasInventoryNumber(t *testing.T) {
	rep := repository.NewMemoryRepository(nil)

	if rep.HasInventoryNumber("INV-1") {
		t.Fatalf("пустой репозиторий не должен содержать номер")
	}

	rep.AddTransport(domain.NewScooter("Самокат", "INV-1", 1, 5))
	if !rep.HasInventoryNumber("INV-1") {
		t.Fatalf("номер транспорта должен найтись")
	}

	rep.AddThing(domain.NewHelmet("Шлем", "INV-2"))
	if !rep.HasInventoryNumber("INV-2") {
		t.Fatalf("номер вещи должен найтись")
	}
	if rep.HasInventoryNumber("INV-3") {
		t.Fatalf("чужой номер не должен найтись")
	}
}

func TestTotalDailyEnergyKWh(t *testing.T) {
	rep := repository.NewMemoryRepository(nil)
	rep.AddTransport(domain.NewScooter("Самокат", "INV-1", 1.5, 5))
	rep.AddTransport(domain.NewBicycle("Велосипед", "INV-2", 2.0, 5))
	rep.AddThing(domain.NewHelmet("Шлем", "INV-3"))
	rep.AddThing(energyBox{name: "Аккумулятор", inv: "INV-4", kwh: 0.5})

	got := rep.TotalDailyEnergyKWh()
	if got != 4.0 {
		t.Fatalf("ожидалось 4.0 кВт·ч, получено %f", got)
	}
}

func TestGetAllItems(t *testing.T) {
	rep := repository.NewMemoryRepository(nil)
	rep.AddTransport(domain.NewScooter("Самокат", "INV-1", 1, 5))
	rep.AddThing(domain.NewHelmet("Шлем", "INV-2"))

	items := rep.GetAllItems()
	if len(items) != 2 {
		t.Fatalf("ожидалось 2 элемента, получено %d", len(items))
	}
	if _, ok := items[0].(domain.Transport); !ok {
		t.Fatalf("первым должен идти транспорт, получено %T", items[0])
	}
	if _, ok := items[1].(domain.Thing); !ok {
		t.Fatalf("вторым должна идти вещь, получено %T", items[1])
	}
}
