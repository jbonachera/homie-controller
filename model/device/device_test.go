package device

import "testing"

func TestNew(t *testing.T) {
	device := New("azertyuip")
	if device.Id != "azertyuip" {
		t.Error("Wrong device id: expected azertyuip, got ", device.Id)
	}
	if device.Online {
		t.Error("New device is online, and should not")
	}
}

func TestSet(t *testing.T) {
	device := New("azertyuip")
	device.Set("$online", "true")
	if device.Online != true {
		t.Error("Setting $online failed")
	}
	device.Set("$stats/signal", "80")
	if device.Stats.Signal != 80 {
		t.Error("Setting $stats/signal failed")
	}
	device.Set("temperature/$properties", "degrees,unit")
	device.Set("temperature/$type", "temperature")
	device.Set("temperature/degrees", "24.3")
	device.Set("temperature/unit", "c")
	if device.Nodes["temperature"].Properties["degrees"] != "24.3" {
		t.Error("Setting temperature/ failed")
	}
	if device.Nodes["temperature"].Properties["unit"] != "c" {
		t.Error("Setting temperature/unit failed")
	}
	if device.Nodes["temperature"].Type != "temperature" {
		t.Error("Setting temperature/$type failed; wanted temperature, got", device.Nodes)
	}
}
func TestSetNodePropertyRemoval(t *testing.T) {
	device := New("azertyuip")
	device.Set("temperature/$properties", "degrees,unit,killMe")
	device.Set("temperature/killMe", "bar")
	device.Set("temperature/$properties", "degrees,unit")

	if device.Nodes["temperature"].Properties["killMe"] == "bar" {
		t.Error("Node properties update failed")
	}
}

func TestListTypes(t *testing.T) {
	device := New("azertyuip")
	device.Set("temperature/$properties", "degrees")
	device.Set("temperature/degrees", "24.3")
	device.Set("temperature/$type", "temperature")

	types := device.ListTypes()
	if types[0] != "temperature" {
		t.Error("could not list exposed types")
	}
}
func TestListTypesMultiple(t *testing.T) {
	device := New("azertyuip")
	device.Set("temperature/$properties", "degrees")
	device.Set("temperature/degrees", "24.3")
	device.Set("temperature/$type", "temperature")

	device.Set("humidity/$properties", "percent")
	device.Set("humidity/percent", "49")
	device.Set("humiduty/$type", "humidity")

	types := device.ListTypes()
	for _, wantedString := range []string{"humidity", "temperature"} {
		found := false
		for _, presentType := range types {
			if presentType == wantedString {
				found = true
			}
		}
		if !found {
			t.Error("could not find type ", wantedString)
		}
	}
}

func TestListTypesDedup(t *testing.T) {
	device := New("azertyuip")
	device.Set("temperature/$properties", "degrees")
	device.Set("temperature/degrees", "24.3")
	device.Set("temperature/$type", "temperature")

	device.Set("temperature2/$properties", "degrees")
	device.Set("temperature2/degrees", "24.3")
	device.Set("temperature2/$type", "temperature")

	types := device.ListTypes()
	if len(types) != 1 {
		t.Error("got duplicates types")
	}
}
