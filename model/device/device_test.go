package device

import (
	"github.com/jbonachera/homie-controller/mocks/mqtt"
	_ "github.com/jbonachera/homie-controller/model/node/type/temperature"
	"testing"
)

var message mqtt.MessageMock = mqtt.NewMessage(
	"devices/u1234/temperature/$type",
	"temperature",
)

func TestNew(t *testing.T) {
	device := New("azertyuip", "devices/")
	if device.Id != "azertyuip" {
		t.Error("Wrong device id: expected azertyuip, got ", device.Id)
	}
	if device.Online {
		t.Error("New device is online, and should not")
	}
}

func TestSet(t *testing.T) {
	device := New("azertyuip", "devices/")
	device.Set("$online", "true")
	if device.Online != true {
		t.Error("Setting $online failed")
	}
	device.Set("$stats/signal", "80")
	if device.Stats.Signal != 80 {
		t.Error("Setting $stats/signal failed")
	}
}

func TestAddNode(t *testing.T) {
	device := New("azertyuip", "devices/")
	client := mqtt.NewMockClient(true, "old/topic")
	device.MQTTNodeHandler(client, message)
	if len(device.Nodes) != 1 {
		t.Error("adding node failed")
	}
	if client.Topic != "devices/azertyuip/temperature/room" {
		t.Error("subscription to wrong topic: want devices/azertyuip/temperature/room, got ", client.Topic)
	}
}

func TestDevice_MQTTNodeHandler(t *testing.T) {
	updateMessage := mqtt.NewMessage(
		"devices/u1234/$localip",
		"127.0.0.1",
	)
	client := mqtt.NewMockClient(true, "old/topic")
	device := New("u1234", "devices/")

	device.MQTTNodeHandler(client, updateMessage)
	if device.Localip != "127.0.0.1"{
		t.Error("mqtt message did not update property LocalIP: wanted 127.0.0.1, got", device.Localip)
	}

}
