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
}
