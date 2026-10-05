package domain

type Transport interface {
	InventoryItem
	EnergyConsumer
	GetSimplicity() int // 1...10
}
