package domain

import "errors"

var (
	ErrInvalidCommand       = errors.New("неизвестная команда, введите число от 0 до 5")
	ErrInvalidTransportName = errors.New("неизвестный тип транспорта, доступно: самокат, велосипед, электровелосипед")
	ErrInvalidThingName     = errors.New("неизвестный тип вещи, доступно: шлем, док-станция")
	ErrEmptyName            = errors.New("название не может быть пустым")
	ErrEmptyInventoryNumber = errors.New("инвентарный номер не может быть пустым")
	ErrDuplicateInventory   = errors.New("инвентарный номер уже занят")
	ErrEnergyNotANumber     = errors.New("суточный расход энергии должен быть числом, например 1.5")
	ErrSimplicityNotANumber = errors.New("простота должна быть целым числом от 1 до 10")
	ErrInspectionRejected   = errors.New("транспорт не прошёл проверку сервисного центра")
)
