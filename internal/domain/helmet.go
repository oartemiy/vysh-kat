package domain

type Helmet struct {
	name            string
	inventoryNumber string
}

func NewHelmet(name, inventoryNumber string) Helmet {
	return Helmet{name, inventoryNumber}
}

func (h Helmet) Name() string {
	return h.name
}

func (h Helmet) InventoryNumber() string {
	return h.inventoryNumber
}
