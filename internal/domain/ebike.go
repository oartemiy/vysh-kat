package domain

type EBike struct {
	name             string
	inventoryNumber  string
	dailyEnergyKWh   float64
	noviceSimplicity int
}

func (k EBike) GetSimplicity() int {
	return k.noviceSimplicity
}

func NewEBike(name, inv string, energy float64, simplicity int) EBike {
	return EBike{name, inv, energy, simplicity}
}

func (k EBike) Name() string            { return k.name }
func (k EBike) InventoryNumber() string { return k.inventoryNumber }
func (k EBike) DailyEnergyKWh() float64 { return k.dailyEnergyKWh }
func (k EBike) NoviceSimplicity() int   { return k.noviceSimplicity }
