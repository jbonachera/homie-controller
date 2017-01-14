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
