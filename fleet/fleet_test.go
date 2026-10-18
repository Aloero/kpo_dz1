package fleet_test

import (
	"math"
	"testing"

	"vihkat/fleet"
)

func TestRejectedTransportDoesNotChangeFleet(t *testing.T) {
	park := fleet.NewFleet(fleet.ServiceCenter{})
	if err := park.AddThing(station(t, "C-1", 2)); err != nil {
		t.Fatal(err)
	}
	transport := scooter(t, "S-1", 10, false, 9)

	accepted, err := park.AddTransport(transport)

	if err != nil || accepted {
		t.Fatalf("AddTransport() = %v, %v; want false, nil", accepted, err)
	}
	if len(park.Inventory()) != 1 || park.TransportCount() != 0 ||
		park.DailyEnergyConsumption() != 2 || len(park.BeginnerTransport()) != 0 {
		t.Fatal("rejected transport changed fleet reports")
	}
}

func TestReportsIncludeTransportAndThings(t *testing.T) {
	park := fleet.NewFleet(fleet.ServiceCenter{})
	accept(t, park, scooter(t, "S-1", 6, true, 1.25))
	accept(t, park, bicycle(t, "B-1", 5, true))
	for _, item := range []fleet.InventoryItem{helmet(t, "H-1"), station(t, "C-1", 2.5)} {
		if err := park.AddThing(item); err != nil {
			t.Fatal(err)
		}
	}

	if math.Abs(park.DailyEnergyConsumption()-3.75) > 1e-9 || park.TransportCount() != 2 {
		t.Fatal("incorrect count or daily energy")
	}
	inventory := park.Inventory()
	want := []string{"S-1", "B-1", "H-1", "C-1"}
	if len(inventory) != len(want) {
		t.Fatalf("inventory length = %d, want %d", len(inventory), len(want))
	}
	for i, number := range want {
		if inventory[i].InventoryNumber() != number || inventory[i].Name() == "" {
			t.Fatalf("incorrect inventory entry: %v", inventory[i])
		}
	}
}

func TestBeginnerBoundary(t *testing.T) {
	park := fleet.NewFleet(fleet.ServiceCenter{})
	for i, ease := range []int{1, 5, 6, 10} {
		accept(t, park, bicycle(t, []string{"B-1", "B-5", "B-6", "B-10"}[i], ease, true))
	}
	if _, err := park.AddTransport(scooter(t, "S-1", 10, false, 1)); err != nil {
		t.Fatal(err)
	}
	if err := park.AddThing(helmet(t, "H-1")); err != nil {
		t.Fatal(err)
	}

	beginners := park.BeginnerTransport()
	if len(beginners) != 2 || beginners[0].InventoryNumber() != "B-6" || beginners[1].InventoryNumber() != "B-10" {
		t.Fatalf("incorrect selection at the 5/6 boundary: %v", beginners)
	}
}

func TestInvalidAdditionDoesNotChangeFleet(t *testing.T) {
	park := fleet.NewFleet(fleet.ServiceCenter{})
	accept(t, park, bicycle(t, "B-1", 6, true))
	if err := park.AddThing(helmet(t, " b-1 ")); err == nil {
		t.Fatal("thing reused an existing inventory number")
	}
	if _, err := park.AddTransport(bicycle(t, "B-1", 6, true)); err == nil {
		t.Fatal("duplicate transport was accepted")
	}
	if err := park.AddThing(bicycle(t, "B-2", 6, false)); err == nil {
		t.Fatal("transport bypassed inspection through AddThing")
	}
	if len(park.Inventory()) != 1 || park.TransportCount() != 1 {
		t.Fatal("invalid addition changed the fleet")
	}
}

type electricBicycle struct {
	fleet.Bicycle
}

func (electricBicycle) DailyEnergyConsumption() float64 { return 0.5 }

func TestNewTransportType(t *testing.T) {
	park := fleet.NewFleet(fleet.ServiceCenter{})
	transport := electricBicycle{Bicycle: bicycle(t, "E-1", 6, true)}
	accept(t, park, transport)
	if park.TransportCount() != 1 || park.DailyEnergyConsumption() != 0.5 || len(park.BeginnerTransport()) != 1 {
		t.Fatal("existing operations do not support a new transport type")
	}
}

func accept(t *testing.T, park *fleet.Fleet, transport fleet.Transport) {
	t.Helper()
	if accepted, err := park.AddTransport(transport); err != nil || !accepted {
		t.Fatalf("AddTransport() = %v, %v; want true, nil", accepted, err)
	}
}

func bicycle(t *testing.T, number string, ease int, operational bool) fleet.Bicycle {
	t.Helper()
	item, err := fleet.NewBicycle("Велосипед", number, ease, operational)
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func scooter(t *testing.T, number string, ease int, operational bool, energy float64) fleet.Scooter {
	t.Helper()
	item, err := fleet.NewScooter("Самокат", number, ease, operational, energy)
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func helmet(t *testing.T, number string) fleet.Helmet {
	t.Helper()
	item, err := fleet.NewHelmet("Шлем", number)
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func station(t *testing.T, number string, energy float64) fleet.ChargingStation {
	t.Helper()
	item, err := fleet.NewChargingStation("Станция", number, energy)
	if err != nil {
		t.Fatal(err)
	}
	return item
}
