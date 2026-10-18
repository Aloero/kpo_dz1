package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"vihkat/fleet"
)

type consoleApp struct {
	fleet  *fleet.Fleet
	input  *bufio.Scanner
	output io.Writer
}

func newConsoleApp(park *fleet.Fleet, input io.Reader, output io.Writer) *consoleApp {
	return &consoleApp{fleet: park, input: bufio.NewScanner(input), output: output}
}

func (a *consoleApp) Run() error {
	fmt.Fprintln(a.output, "ВышКат. Учёт парка")
	for {
		fmt.Fprintln(a.output, "\n1. Добавить транспорт\n2. Добавить вещь\n3. Отчёт\n4. Транспорт для новичков\n5. Инвентаризация\n0. Выход")
		command, err := a.readText("Команда: ")
		if err == nil {
			switch command {
			case "1":
				err = a.addTransport()
			case "2":
				err = a.addThing()
			case "3":
				a.showReport()
			case "4":
				a.showBeginnerTransport()
			case "5":
				a.showInventory()
			case "0":
				return nil
			default:
				fmt.Fprintln(a.output, "Неизвестная команда. Введите число от 0 до 5.")
			}
		}
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(a.output, "\nВвод завершён.")
			return nil
		}
		if err != nil {
			if a.input.Err() != nil {
				return err
			}
			fmt.Fprintf(a.output, "Ошибка: %v.\n", err)
		}
	}
}

func (a *consoleApp) addTransport() error {
	kind, err := a.readInteger("Вид (1 — электросамокат, 2 — велосипед): ", 1, 2)
	if err != nil {
		return err
	}
	name, err := a.readRequiredText("Название: ")
	if err != nil {
		return err
	}
	number, err := a.readRequiredText("Инвентарный номер: ")
	if err != nil {
		return err
	}
	ease, err := a.readInteger("Простота для новичка (1,10): ", 1, 10)
	if err != nil {
		return err
	}
	operational, err := a.readInteger("Исправность по результатам осмотра (1 — исправен, 0 — неисправен): ", 0, 1)
	if err != nil {
		return err
	}

	var transport fleet.Transport
	if kind == 1 {
		energy, readErr := a.readEnergy()
		if readErr != nil {
			return readErr
		}
		transport, err = fleet.NewScooter(name, number, ease, operational == 1, energy)
	} else {
		transport, err = fleet.NewBicycle(name, number, ease, operational == 1)
	}
	if err != nil {
		return err
	}
	accepted, err := a.fleet.AddTransport(transport)
	if err != nil {
		return err
	}
	if accepted {
		fmt.Fprintln(a.output, "Техника принята в парк.")
	} else {
		fmt.Fprintln(a.output, "Техосмотр не пройден. Техника отклонена.")
	}
	return nil
}

func (a *consoleApp) addThing() error {
	kind, err := a.readInteger("Вид (1 — шлем, 2 — зарядная станция): ", 1, 2)
	if err != nil {
		return err
	}
	name, err := a.readRequiredText("Название: ")
	if err != nil {
		return err
	}
	number, err := a.readRequiredText("Инвентарный номер: ")
	if err != nil {
		return err
	}

	var item fleet.InventoryItem
	if kind == 1 {
		item, err = fleet.NewHelmet(name, number)
	} else {
		energy, readErr := a.readEnergy()
		if readErr != nil {
			return readErr
		}
		item, err = fleet.NewChargingStation(name, number, energy)
	}
	if err != nil {
		return err
	}
	if err := a.fleet.AddThing(item); err != nil {
		return err
	}
	fmt.Fprintln(a.output, "Вещь поставлена на баланс.")
	return nil
}

func (a *consoleApp) showReport() {
	total, transport := len(a.fleet.Inventory()), a.fleet.TransportCount()
	fmt.Fprintf(a.output, "Транспорт: %d\nВещи: %d\nВсего на балансе: %d\n", transport, total-transport, total)
	fmt.Fprintf(a.output, "Энергия за сутки: %.2f кВт·ч\n", a.fleet.DailyEnergyConsumption())
}

func (a *consoleApp) showBeginnerTransport() {
	transport := a.fleet.BeginnerTransport()
	if len(transport) == 0 {
		fmt.Fprintln(a.output, "Подходящего транспорта пока нет.")
	}
	for _, item := range transport {
		fmt.Fprintf(a.output, "%s | инв. # %s | простота %d/10\n", item.Name(), item.InventoryNumber(), item.BeginnerEase())
	}
}

func (a *consoleApp) showInventory() {
	inventory := a.fleet.Inventory()
	if len(inventory) == 0 {
		fmt.Fprintln(a.output, "На балансе пока нет объектов.")
	}
	for _, item := range inventory {
		fmt.Fprintf(a.output, "%s | инв. № %s\n", item.Name(), item.InventoryNumber())
	}
}

func (a *consoleApp) readText(prompt string) (string, error) {
	fmt.Fprint(a.output, prompt)
	if !a.input.Scan() {
		if err := a.input.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(a.input.Text()), nil
}

func (a *consoleApp) readRequiredText(prompt string) (string, error) {
	for {
		text, err := a.readText(prompt)
		if err != nil {
			return "", err
		}
		if text != "" {
			return text, nil
		}
		fmt.Fprintln(a.output, "Значение не должно быть пустым.")
	}
}

func (a *consoleApp) readInteger(prompt string, minimum, maximum int) (int, error) {
	for {
		text, err := a.readText(prompt)
		if err != nil {
			return 0, err
		}
		value, err := strconv.Atoi(text)
		if err == nil && value >= minimum && value <= maximum {
			return value, nil
		}
		fmt.Fprintf(a.output, "Введите целое число от %d до %d.\n", minimum, maximum)
	}
}

func (a *consoleApp) readEnergy() (float64, error) {
	for {
		text, err := a.readText("Суточный расход энергии, кВт ч: ")
		if err != nil {
			return 0, err
		}
		value, err := strconv.ParseFloat(strings.ReplaceAll(text, ",", "."), 64)
		if err == nil && value >= 0 && !math.IsNaN(value) && !math.IsInf(value, 0) {
			return value, nil
		}
		fmt.Fprintln(a.output, "Введите конечное неотрицательное число, например 1,5 или 1.5.")
	}
}
