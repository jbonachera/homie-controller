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
	device.Set("temperature/degrees", "24.3")
	device.Set("temperature/unit", "c")
	if device.Nodes["temperature"].Properties["degrees"] != "24.3" {
		t.Error("Setting temperature/ failed")
	}
}
