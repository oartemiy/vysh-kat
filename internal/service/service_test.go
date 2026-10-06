package service_test

import (
	"errors"
	"strings"
	"testing"

	"vysh-kat/internal/domain"
	"vysh-kat/internal/inspection"
	"vysh-kat/internal/mocks"
	"vysh-kat/internal/service"
)

func newTestService() (*service.Service, *mocks.MockRepository, *mocks.MockServiceCenter) {
	rep := &mocks.MockRepository{}
	sc := &mocks.MockServiceCenter{Result: inspection.ResultOfInspection{Accepted: true, Reason: "ok"}}
	return service.NewService(rep, sc), rep, sc
}

func TestAddTransport_Success(t *testing.T) {
	svc, rep, sc := newTestService()

	err := svc.AddTransport(domain.NewScooter("Красный самокат", "INV-1", 1.5, 8))
	if err != nil {
		t.Fatalf("ожидался успех, получена ошибка: %v", err)
	}
	if len(rep.Transports) != 1 {
		t.Fatalf("ожидался 1 транспорт в репозитории, получено %d", len(rep.Transports))
	}
	if sc.InspectCalls != 1 {
		t.Fatalf("ожидался 1 вызов Inspect, получено %d", sc.InspectCalls)
	}
}

func TestAddTransport_InspectionRejected(t *testing.T) {
	rep := &mocks.MockRepository{}
	reason := "суточный расход энергии не может быть отрицательным"
	sc := &mocks.MockServiceCenter{Result: inspection.ResultOfInspection{Accepted: false, Reason: reason}}
	svc := service.NewService(rep, sc)

	err := svc.AddTransport(domain.NewScooter("Самокат", "INV-1", -1, 8))
	if !errors.Is(err, domain.ErrInspectionRejected) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrInspectionRejected, err)
	}
	if !strings.Contains(err.Error(), reason) {
		t.Fatalf("в сообщении нет причины отказа: %v", err)
	}
	if len(rep.Transports) != 0 {
		t.Fatalf("отклонённый транспорт не должен попадать в репозиторий")
	}
	if sc.InspectCalls != 1 {
		t.Fatalf("ожидался 1 вызов Inspect, получено %d", sc.InspectCalls)
	}
}

func TestAddTransport_EmptyName(t *testing.T) {
	svc, rep, sc := newTestService()

	err := svc.AddTransport(domain.NewScooter("   ", "INV-1", 1, 8))
	if !errors.Is(err, domain.ErrEmptyName) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrEmptyName, err)
	}
	if sc.InspectCalls != 0 {
		t.Fatalf("валидация должна выполняться до осмотра, вызовов: %d", sc.InspectCalls)
	}
	if len(rep.Transports) != 0 {
		t.Fatalf("репозиторий должен оставаться пустым")
	}
}

func TestAddTransport_EmptyInventoryNumber(t *testing.T) {
	svc, _, sc := newTestService()

	err := svc.AddTransport(domain.NewScooter("Самокат", "  ", 1, 8))
	if !errors.Is(err, domain.ErrEmptyInventoryNumber) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrEmptyInventoryNumber, err)
	}
	if sc.InspectCalls != 0 {
		t.Fatalf("валидация должна выполняться до осмотра, вызовов: %d", sc.InspectCalls)
	}
}

func TestAddTransport_DuplicateInventoryNumber(t *testing.T) {
	svc, rep, sc := newTestService()
	rep.AddTransport(domain.NewBicycle("Старый велосипед", "INV-1", 1, 5))

	err := svc.AddTransport(domain.NewScooter("Новый самокат", "INV-1", 1, 8))
	if !errors.Is(err, domain.ErrDuplicateInventory) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrDuplicateInventory, err)
	}
	if !strings.Contains(err.Error(), "INV-1") {
		t.Fatalf("в сообщении должен быть номер: %v", err)
	}
	if sc.InspectCalls != 0 {
		t.Fatalf("дубликат должен отсекаться до осмотра, вызовов: %d", sc.InspectCalls)
	}
	if len(rep.Transports) != 1 {
		t.Fatalf("ожидался 1 транспорт, получено %d", len(rep.Transports))
	}
}

func TestAddThing_Success(t *testing.T) {
	svc, rep, _ := newTestService()

	err := svc.AddThing(domain.NewHelmet("Шлем", "INV-10"))
	if err != nil {
		t.Fatalf("ожидался успех, получена ошибка: %v", err)
	}
	if len(rep.Things) != 1 {
		t.Fatalf("ожидалась 1 вещь, получено %d", len(rep.Things))
	}
}

func TestAddThing_DuplicateAcrossCollections(t *testing.T) {
	svc, rep, _ := newTestService()
	rep.AddTransport(domain.NewScooter("Самокат", "INV-1", 1, 5))

	err := svc.AddThing(domain.NewHelmet("Шлем", "INV-1"))
	if !errors.Is(err, domain.ErrDuplicateInventory) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrDuplicateInventory, err)
	}
	if len(rep.Things) != 0 {
		t.Fatalf("вещь не должна попадать в репозиторий")
	}
}

func TestAddThing_EmptyName(t *testing.T) {
	svc, rep, _ := newTestService()

	err := svc.AddThing(domain.NewDockingStation("", "INV-2"))
	if !errors.Is(err, domain.ErrEmptyName) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrEmptyName, err)
	}
	if len(rep.Things) != 0 {
		t.Fatalf("репозиторий должен оставаться пустым")
	}
}

func TestService_DelegatesQueries(t *testing.T) {
	svc, rep, _ := newTestService()
	rep.TotalEnergy = 42.5
	rep.NoviceList = []domain.Transport{domain.NewBicycle("Велосипед", "INV-1", 1, 7)}
	rep.ExtraItems = []domain.InventoryItem{domain.NewHelmet("Шлем", "INV-2")}

	if got := svc.TotalDailyEnergyKWh(); got != 42.5 {
		t.Fatalf("ожидалось 42.5, получено %f", got)
	}
	novice := svc.NoviceSuitableTransports()
	if len(novice) != 1 || novice[0].Name() != "Велосипед" {
		t.Fatalf("неожиданный список новичков: %v", novice)
	}
	if items := svc.AllItems(); len(items) != 1 {
		t.Fatalf("ожидалась 1 вещь из ExtraItems, получено %d", len(items))
	}
}
