package domain

type DockingStation struct {
	name            string
	inventoryNumber string
}

func NewDockingStation(name, inventoryNumber string) DockingStation {
	return DockingStation{name, inventoryNumber}
}

func (h DockingStation) Name() string {
	return h.name
}

func (h DockingStation) InventoryNumber() string {
	return h.inventoryNumber
}
