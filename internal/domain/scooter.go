package domain

type Scooter struct {
	name             string
	inventoryNumber  string
	dailyEnergyKWh   float64
	noviceSimplicity int
}

func NewScooter(name, inv string, energy float64, simplicity int) Scooter {
	return Scooter{name, inv, energy, simplicity}
}

func (k Scooter) Name() string            { return k.name }
func (k Scooter) InventoryNumber() string { return k.inventoryNumber }
func (k Scooter) DailyEnergyKWh() float64 { return k.dailyEnergyKWh }
func (k Scooter) NoviceSimplicity() int   { return k.noviceSimplicity }
