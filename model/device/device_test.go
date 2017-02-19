package device

import (
	"github.com/jbonachera/homie-controller/mocks/mqtt"
	MQTT "github.com/jbonachera/homie-controller/messaging"
	_ "github.com/jbonachera/homie-controller/model/node/type/temperature"
	"testing"
	"github.com/jbonachera/homie-controller/model/homieMessage"
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
	var message homieMessage.HomieMessage
	message, _ = homieMessage.New(mqtt.NewMessage(
		"devices/u1234/temperature/$type",
		"temperature",
	), "devices/")
	client := mqtt.NewMockClient(true, "old/topic")
	registrator := func(topic string, callback MQTT.CallbackHandler){
		client.Topic = topic
	}
	device := New("azertyuip", "devices/")
	device.SetRegistrator(registrator)
	device.MQTTNodeHandler( message)
	if len(device.Nodes) != 1 {
		t.Error("adding node failed")
	}
	if client.Topic != "devices/azertyuip/temperature/room" {
		t.Error("subscription to wrong topic: want devices/azertyuip/temperature/room, got ", client.Topic)
	}
}

func TestDevice_MQTTNodeHandler(t *testing.T) {
	var message homieMessage.HomieMessage
	updateMessage := mqtt.NewMessage(
		"devices/u1234/$localip",
		"127.0.0.1",
	)
	message, _ = homieMessage.New(updateMessage, "devices/")
	device := New("u1234", "devices/")
	device.MQTTNodeHandler(message)
	if device.Localip != "127.0.0.1"{
		t.Error("messaging message did not update property LocalIP: wanted 127.0.0.1, got", device.Localip)
	}

}
