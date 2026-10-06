package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"vysh-kat/internal/domain"
)

type App struct {
	container *Container
	reader    *bufio.Reader
	writer    *bufio.Writer
}

func NewApp(container *Container) *App {
	return &App{container: container, reader: bufio.NewReader(os.Stdin), writer: bufio.NewWriter(os.Stdout)}
}

func (a *App) printMenu() {
	a.writer.WriteString("------- ВышКат -------\n")
	a.writer.WriteString("0. Выход\n")
	a.writer.WriteString("1. Добавить транспорт\n")
	a.writer.WriteString("2. Добавить вещь\n")
	a.writer.WriteString("3. Суммарное суточное потребление энергии\n")
	a.writer.WriteString("4. Транспорт для новичков\n")
	a.writer.WriteString("5. Вывести все единицы и вещи\n")
	a.writer.Flush()
}

func (a *App) readInput(prompt string) (string, error) {
	a.writer.WriteString(prompt)
	a.writer.Flush()
	text, err := a.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func (a *App) finishGame(exitCode int, err error) {
	if err != nil {
		a.writer.WriteString(fmt.Sprintf("Ошибка: %s\n", err))
		a.writer.Flush()
	}
	os.Exit(exitCode)
}

func (a *App) printErrorAndContinue(err error) {
	a.writer.WriteString(fmt.Sprintf("Ошибка: %s\n", err))
	a.writer.Flush()
}

func (a *App) addTransport() error {
	a.writer.WriteString("\n--- Добавление транспорта ---\n")
	kind, err := a.readInput("Тип транспорта (самокат / велосипед / электровелосипед): ")
	if err != nil {
		return err
	}
	name, err := a.readInput("Имя средства: ")
	if err != nil {
		return err
	}
	inv, err := a.readInput("Инвентарный номер: ")
	if err != nil {
		return err
	}
	energyStr, err := a.readInput("Суточный расход энергии (кВт·ч): ")
	if err != nil {
		return err
	}
	energy, err := strconv.ParseFloat(energyStr, 64)
	if err != nil {
		return fmt.Errorf("%w (введено %q)", domain.ErrEnergyNotANumber, energyStr)
	}
	simplicityStr, err := a.readInput("Простота для новичка (1-10): ")
	if err != nil {
		return err
	}
	simplicity, err := strconv.Atoi(simplicityStr)
	if err != nil {
		return fmt.Errorf("%w (введено %q)", domain.ErrSimplicityNotANumber, simplicityStr)
	}
	transport, err := domain.NewTransport(kind, name, inv, energy, simplicity)
	if err != nil {
		return err
	}
	if err := a.container.Service.AddTransport(transport); err != nil {
		return err
	}
	a.writer.WriteString("Транспорт принят в парк\n")
	a.writer.Flush()
	return nil
}

func (a *App) addThing() error {
	a.writer.WriteString("\n--- Добавление вещи ---\n")
	kind, err := a.readInput("Тип вещи (шлем / док-станция): ")
	if err != nil {
		return err
	}
	name, err := a.readInput("Имя вещи: ")
	if err != nil {
		return err
	}
	inv, err := a.readInput("Инвентарный номер: ")
	if err != nil {
		return err
	}
	thing, err := domain.NewThing(kind, name, inv)
	if err != nil {
		return err
	}
	if err := a.container.Service.AddThing(thing); err != nil {
		return err
	}
	a.writer.WriteString("Вещь добавлена на баланс\n")
	a.writer.Flush()
	return nil
}

func (a *App) printTotalEnergy() {
	total := a.container.Service.TotalDailyEnergyKWh()
	a.writer.WriteString("\n--- Суммарное суточное потребление энергии ---\n")
	a.writer.WriteString(fmt.Sprintf("%.2f кВт·ч\n", total))
	a.writer.Flush()
}

func (a *App) printNoviceTransports() {
	transports := a.container.Service.NoviceSuitableTransports()
	a.writer.WriteString("\n--- Транспорт для новичков ---\n")
	if len(transports) == 0 {
		a.writer.WriteString("Нет подходящего транспорта\n")
		a.writer.Flush()
		return
	}
	for _, t := range transports {
		a.writer.WriteString(fmt.Sprintf("- %s (INV: %s), простота: %d\n",
			t.Name(), t.InventoryNumber(), t.GetSimplicity()))
	}
	a.writer.Flush()
}

func (a *App) printAllItems() {
	items := a.container.Service.AllItems()
	a.writer.WriteString("\n--- Все единицы и вещи на балансе ---\n")
	if len(items) == 0 {
		a.writer.WriteString("Список пуст\n")
		a.writer.Flush()
		return
	}
	for _, item := range items {
		switch v := item.(type) {
		case domain.Transport:
			a.writer.WriteString(fmt.Sprintf("[Транспорт] %s (INV: %s)\n",
				v.Name(), v.InventoryNumber()))
		case domain.Thing:
			a.writer.WriteString(fmt.Sprintf("[Вещь]      %s (INV: %s)\n",
				v.Name(), v.InventoryNumber()))
		default:
			a.writer.WriteString(fmt.Sprintf("[Неизвестно] %s (INV: %s)\n",
				item.Name(), item.InventoryNumber()))
		}
	}
	a.writer.Flush()
}

func (a *App) Run() {
	for {
		a.printMenu()
		choice, err := a.readInput("\nВведите номер действия: ")
		if err != nil {
			a.finishGame(-1, err)
		}
		switch choice {
		case "0":
			a.finishGame(0, nil)
		case "1":
			err = a.addTransport()
			if err != nil {
				a.printErrorAndContinue(err)
			}
		case "2":
			err = a.addThing()
			if err != nil {
				a.printErrorAndContinue(err)
			}
		case "3":
			a.printTotalEnergy()
		case "4":
			a.printNoviceTransports()
		case "5":
			a.printAllItems()
		default:
			a.printErrorAndContinue(domain.ErrInvalidCommand)
		}
	}
}
