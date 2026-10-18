package fleet

import (
	"fmt"
	"math"
	"strings"
)

type InventoryItem interface {
	Name() string
	InventoryNumber() string
}

type EnergyConsumer interface {
	DailyEnergyConsumption() float64
}

type Transport interface {
	InventoryItem
	BeginnerEase() int
	IsOperational() bool
}

type asset struct {
	name   string
	number string
}

func newAsset(name, number string) (asset, error) {
	name, number = strings.TrimSpace(name), strings.TrimSpace(number)
	if name == "" || number == "" {
		return asset{}, fmt.Errorf("название и инвентарный номер не должны быть пустыми")
	}
	return asset{name: name, number: number}, nil
}

func (a asset) Name() string            { return a.name }
func (a asset) InventoryNumber() string { return a.number }

type vehicle struct {
	asset
	ease        int
	operational bool
}

func newVehicle(name, number string, ease int, operational bool) (vehicle, error) {
	item, err := newAsset(name, number)
	if err != nil {
		return vehicle{}, err
	}
	if ease < 1 || ease > 10 {
		return vehicle{}, fmt.Errorf("простота должна быть от 1 до 10")
	}
	return vehicle{asset: item, ease: ease, operational: operational}, nil
}

func (v vehicle) BeginnerEase() int   { return v.ease }
func (v vehicle) IsOperational() bool { return v.operational }

type Scooter struct {
	vehicle
	energy float64
}

func NewScooter(name, number string, ease int, operational bool, energy float64) (Scooter, error) {
	transport, err := newVehicle(name, number, ease, operational)
	if err != nil {
		return Scooter{}, err
	}
	if err := validateEnergy(energy); err != nil {
		return Scooter{}, err
	}
	return Scooter{vehicle: transport, energy: energy}, nil
}

func (s Scooter) DailyEnergyConsumption() float64 { return s.energy }

type Bicycle struct {
	vehicle
}

func NewBicycle(name, number string, ease int, operational bool) (Bicycle, error) {
	transport, err := newVehicle(name, number, ease, operational)
	if err != nil {
		return Bicycle{}, err
	}
	return Bicycle{vehicle: transport}, nil
}

type Helmet struct {
	asset
}

func NewHelmet(name, number string) (Helmet, error) {
	item, err := newAsset(name, number)
	if err != nil {
		return Helmet{}, err
	}
	return Helmet{asset: item}, nil
}

type ChargingStation struct {
	asset
	energy float64
}

func NewChargingStation(name, number string, energy float64) (ChargingStation, error) {
	item, err := newAsset(name, number)
	if err != nil {
		return ChargingStation{}, err
	}
	if err := validateEnergy(energy); err != nil {
		return ChargingStation{}, err
	}
	return ChargingStation{asset: item, energy: energy}, nil
}

func (s ChargingStation) DailyEnergyConsumption() float64 { return s.energy }

func validateEnergy(energy float64) error {
	if energy < 0 || math.IsNaN(energy) || math.IsInf(energy, 0) {
		return fmt.Errorf("расход энергии должен быть разумным неотрицательным числом")
	}
	return nil
}
