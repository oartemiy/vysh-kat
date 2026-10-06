package app

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	"vysh-kat/internal/domain"
)

func newTestApp(input string) (*App, *bytes.Buffer) {
	out := &bytes.Buffer{}
	return &App{
		container: NewContainer(),
		reader:    bufio.NewReader(strings.NewReader(input)),
		writer:    bufio.NewWriter(out),
	}, out
}

func TestAddTransport_Success(t *testing.T) {
	a, out := newTestApp("самокат\nКрасный самокат\nINV-1\n1.5\n8\n")

	if err := a.addTransport(); err != nil {
		t.Fatalf("ожидался успех, получена ошибка: %v", err)
	}

	items := a.container.Service.AllItems()
	if len(items) != 1 {
		t.Fatalf("ожидался 1 элемент, получено %d", len(items))
	}
	if items[0].Name() != "Красный самокат" {
		t.Fatalf("имя должно сохраняться как введено, получено %q", items[0].Name())
	}
	if !strings.Contains(out.String(), "Транспорт принят в парк") {
		t.Fatalf("нет подтверждения добавления: %q", out.String())
	}
}

func TestAddTransport_UnknownKind(t *testing.T) {
	a, _ := newTestApp("саамокат\nX\nINV-1\n1\n1\n")

	err := a.addTransport()
	if !errors.Is(err, domain.ErrInvalidTransportName) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrInvalidTransportName, err)
	}
	if items := a.container.Service.AllItems(); len(items) != 0 {
		t.Fatalf("ничего не должно добавиться, получено %d", len(items))
	}
}

func TestAddTransport_InvalidEnergy(t *testing.T) {
	a, _ := newTestApp("самокат\nX\nINV-1\nabc\n")

	err := a.addTransport()
	if !errors.Is(err, domain.ErrEnergyNotANumber) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrEnergyNotANumber, err)
	}
	if !strings.Contains(err.Error(), "abc") {
		t.Fatalf("в сообщении должно быть введённое значение: %v", err)
	}
}

func TestAddTransport_InvalidSimplicity(t *testing.T) {
	a, _ := newTestApp("самокат\nX\nINV-1\n1.5\nвосемь\n")

	err := a.addTransport()
	if !errors.Is(err, domain.ErrSimplicityNotANumber) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrSimplicityNotANumber, err)
	}
}

func TestAddThing_Success(t *testing.T) {
	a, out := newTestApp("шлем\nСтроительный\nINV-2\n")

	if err := a.addThing(); err != nil {
		t.Fatalf("ожидался успех, получена ошибка: %v", err)
	}
	if items := a.container.Service.AllItems(); len(items) != 1 {
		t.Fatalf("ожидалась 1 вещь, получено %d", len(items))
	}
	if !strings.Contains(out.String(), "Вещь добавлена на баланс") {
		t.Fatalf("нет подтверждения добавления: %q", out.String())
	}
}

func TestAddThing_UnknownKind(t *testing.T) {
	a, _ := newTestApp("перчатки\nX\nINV-2\n")

	err := a.addThing()
	if !errors.Is(err, domain.ErrInvalidThingName) {
		t.Fatalf("ожидалась ошибка %v, получена: %v", domain.ErrInvalidThingName, err)
	}
}

func TestPrintAllItems(t *testing.T) {
	a, out := newTestApp("")
	a.container.Service.AddTransport(domain.NewScooter("Самокат", "INV-1", 1, 5))
	a.container.Service.AddThing(domain.NewHelmet("Шлем", "INV-9"))

	a.printAllItems()

	text := out.String()
	if !strings.Contains(text, "[Транспорт] Самокат (INV: INV-1)") {
		t.Fatalf("нет транспорта в выводе: %q", text)
	}
	if !strings.Contains(text, "[Вещь]      Шлем (INV: INV-9)") {
		t.Fatalf("нет вещи в выводе: %q", text)
	}
}

func TestPrintNoviceTransports(t *testing.T) {
	a, out := newTestApp("")
	a.container.Service.AddTransport(domain.NewBicycle("Подходящий", "INV-1", 1, 7))
	a.container.Service.AddTransport(domain.NewBicycle("Сложный", "INV-2", 1, 3))

	a.printNoviceTransports()

	text := out.String()
	if !strings.Contains(text, "Подходящий") {
		t.Fatalf("новичковый транспорт должен попасть в вывод: %q", text)
	}
	if strings.Contains(text, "Сложный") {
		t.Fatalf("неподходящий транспорт не должен попадать в вывод: %q", text)
	}
}

func TestPrintErrorAndContinue_Format(t *testing.T) {
	a, out := newTestApp("")

	a.printErrorAndContinue(domain.ErrInvalidCommand)

	if !strings.HasPrefix(out.String(), "Ошибка: ") || !strings.HasSuffix(out.String(), "\n") {
		t.Fatalf("ошибка должна выводиться с префиксом и переводом строки: %q", out.String())
	}
}
