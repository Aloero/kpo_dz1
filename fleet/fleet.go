package fleet

import (
	"fmt"
	"math"
	"strings"
)

type Inspector interface {
	Inspect(Transport) bool
}

type ServiceCenter struct{}

func (ServiceCenter) Inspect(transport Transport) bool {
	return transport.IsOperational()
}

type Fleet struct {
	inspector Inspector
	inventory []InventoryItem
}

func NewFleet(inspector Inspector) *Fleet {
	return &Fleet{inspector: inspector}
}

func (f *Fleet) AddTransport(transport Transport) (bool, error) {
	if err := f.checkItem(transport); err != nil {
		return false, err
	}
	if transport.BeginnerEase() < 1 || transport.BeginnerEase() > 10 {
		return false, fmt.Errorf("простота должна быть от 1 до 10")
	}
	if !f.inspector.Inspect(transport) {
		return false, nil
	}
	f.inventory = append(f.inventory, transport)
	return true, nil
}

func (f *Fleet) AddThing(item InventoryItem) error {
	if _, isTransport := item.(Transport); isTransport {
		return fmt.Errorf("транспорт нужно добавлять через техосмотр")
	}
	if err := f.checkItem(item); err != nil {
		return err
	}
	f.inventory = append(f.inventory, item)
	return nil
}

func (f *Fleet) Inventory() []InventoryItem {
	return append([]InventoryItem(nil), f.inventory...)
}

func (f *Fleet) TransportCount() int {
	count := 0
	for _, item := range f.inventory {
		if _, ok := item.(Transport); ok {
			count++
		}
	}
	return count
}

func (f *Fleet) DailyEnergyConsumption() float64 {
	total := 0.0
	for _, item := range f.inventory {
		if consumer, ok := item.(EnergyConsumer); ok {
			total += consumer.DailyEnergyConsumption()
		}
	}
	return total
}

func (f *Fleet) BeginnerTransport() []Transport {
	var result []Transport
	for _, item := range f.inventory {
		if transport, ok := item.(Transport); ok && transport.BeginnerEase() >= 6 {
			result = append(result, transport)
		}
	}
	return result
}

func (f *Fleet) checkItem(item InventoryItem) error {
	if item == nil {
		return fmt.Errorf("объект не задан")
	}
	if _, err := newAsset(item.Name(), item.InventoryNumber()); err != nil {
		return err
	}
	for _, existing := range f.inventory {
		if strings.EqualFold(existing.InventoryNumber(), item.InventoryNumber()) {
			return fmt.Errorf("инвентарный номер %s уже занят", item.InventoryNumber())
		}
	}
	if consumer, ok := item.(EnergyConsumer); ok {
		if err := validateEnergy(consumer.DailyEnergyConsumption()); err != nil {
			return err
		}
		// защита от переполнения суммы в +Inf
		if math.IsInf(f.DailyEnergyConsumption()+consumer.DailyEnergyConsumption(), 0) {
			return fmt.Errorf("суммарный расход энергии слишком велик")
		}
	}
	return nil
}
