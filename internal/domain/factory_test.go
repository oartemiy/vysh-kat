package domain_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"vysh-kat/internal/domain"
)

func TestNewTransport(t *testing.T) {
	tests := []struct {
		kind     string
		wantType string
	}{
		{"самокат", "domain.Scooter"},
		{" Самокат ", "domain.Scooter"},
		{"велосипед", "domain.Bicycle"},
		{"ВЕЛОСИПЕД", "domain.Bicycle"},
		{"электровелосипед", "domain.EBike"},
	}
	for _, tt := range tests {
		tr, err := domain.NewTransport(tt.kind, "Название", "INV-1", 1.5, 7)
		if err != nil {
			t.Fatalf("%q: не ожидалась ошибка: %v", tt.kind, err)
		}
		if got := fmt.Sprintf("%T", tr); got != tt.wantType {
			t.Fatalf("%q: ожидался тип %s, получен %s", tt.kind, tt.wantType, got)
		}
		if tr.Name() != "Название" || tr.InventoryNumber() != "INV-1" {
			t.Fatalf("%q: неверно переданы поля: %s / %s", tt.kind, tr.Name(), tr.InventoryNumber())
		}
	}
}

func TestNewTransport_UnknownKind(t *testing.T) {
	_, err := domain.NewTransport("саамокат", "Название", "INV-1", 1, 5)
	if !errors.Is(err, domain.ErrInvalidTransportName) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrInvalidTransportName, err)
	}
	if !strings.Contains(err.Error(), "саамокат") {
		t.Fatalf("в сообщении должно быть введённое значение: %v", err)
	}
}

func TestNewThing(t *testing.T) {
	tests := []struct {
		kind     string
		wantType string
	}{
		{"шлем", "domain.Helmet"},
		{"Шлем", "domain.Helmet"},
		{"док-станция", "domain.DockingStation"},
		{"докстанция", "domain.DockingStation"},
	}
	for _, tt := range tests {
		th, err := domain.NewThing(tt.kind, "Название", "INV-1")
		if err != nil {
			t.Fatalf("%q: не ожидалась ошибка: %v", tt.kind, err)
		}
		if got := fmt.Sprintf("%T", th); got != tt.wantType {
			t.Fatalf("%q: ожидался тип %s, получен %s", tt.kind, tt.wantType, got)
		}
	}
}

func TestNewThing_UnknownKind(t *testing.T) {
	_, err := domain.NewThing("перчатки", "Название", "INV-1")
	if !errors.Is(err, domain.ErrInvalidThingName) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrInvalidThingName, err)
	}
	if !strings.Contains(err.Error(), "перчатки") {
		t.Fatalf("в сообщении должно быть введённое значение: %v", err)
	}
}
