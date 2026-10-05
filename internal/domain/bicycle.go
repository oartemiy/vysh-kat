package domain

type Bicycle struct {
	name             string
	inventoryNumber  string
	dailyEnergyKWh   float64
	noviceSimplicity int
}

func NewBicycle(name, inv string, energy float64, simplicity int) Bicycle {
	return Bicycle{name, inv, energy, simplicity}
}

func (k Bicycle) Name() string            { return k.name }
func (k Bicycle) InventoryNumber() string { return k.inventoryNumber }
func (k Bicycle) DailyEnergyKWh() float64 { return k.dailyEnergyKWh }
func (k Bicycle) NoviceSimplicity() int   { return k.noviceSimplicity }
