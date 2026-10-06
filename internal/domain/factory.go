package domain

import (
	"fmt"
	"strings"
)

func normalizeKind(kind string) string {
	return strings.ToLower(strings.TrimSpace(kind))
}

func NewTransport(kind, name, inv string, energy float64, simplicity int) (Transport, error) {
	switch normalizeKind(kind) {
	case "самокат":
		return NewScooter(name, inv, energy, simplicity), nil
	case "электровелосипед":
		return NewEBike(name, inv, energy, simplicity), nil
	case "велосипед":
		return NewBicycle(name, inv, energy, simplicity), nil
	default:
		return nil, fmt.Errorf("введено %q: %w", kind, ErrInvalidTransportName)
	}
}

func NewThing(kind, name, inv string) (Thing, error) {
	switch normalizeKind(kind) {
	case "шлем":
		return NewHelmet(name, inv), nil
	case "док-станция", "докстанция":
		return NewDockingStation(name, inv), nil
	default:
		return nil, fmt.Errorf("введено %q: %w", kind, ErrInvalidThingName)
	}
}
