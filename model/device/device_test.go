package device

import (
	MQTT "github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/mocks/mqtt"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/node"
	_ "github.com/jbonachera/homie-controller/model/node/type/temperature"
	"github.com/jbonachera/homie-controller/model/search"
	"testing"
)

func TestNew(t *testing.T) {
	device := New("azertyuip", "devices/")
	if device.Id != "azertyuip" {
		t.Error("Wrong device id: expected azertyuip, got ", device.Id)
	}
	if device.Online {
		t.Error("Extract device is online, and should not")
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
	message, _ = homieMessage.Extract(mqtt.NewMessage(
		"devices/u1234/temperature/$type",
		"temperature",
	), "devices/")
	client := mqtt.NewMockClient(true, "old/topic")
	registrator := func(topic string, callback MQTT.CallbackHandler) {
		client.Topic = topic
	}
	device := New("azertyuip", "devices/")
	device.SetRegistrator(registrator)
	device.MQTTNodeHandler(message)
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
	message, _ = homieMessage.Extract(updateMessage, "devices/")
	device := New("u1234", "devices/")
	device.MQTTNodeHandler(message)
	if device.Localip != "127.0.0.1" {
		t.Error("messaging message did not update property LocalIP: wanted 127.0.0.1, got", device.Localip)
	}

}

func TestDevice_Match(t *testing.T) {
	searchTerms := map[string]string{
		"online": "true",
	}
	opts := search.Options{Terms: searchTerms}
	device := New("u1234", "devices/")
	device.Online = false
	if device.Match(opts) {
		t.Error("device should not match the requested search")
	}
	device.Online = true
	if !device.Match(opts) {
		t.Error("device should match the requested search")
	}

}
func TestDevice_Match2(t *testing.T) {
	searchTerms := map[string]string{
		"type": "temperature",
	}
	opts := search.Options{Terms: searchTerms}
	device := New("u1234", "devices/")
	if device.Match(opts) {
		t.Error("device should not match the requested search")
	}
	device.Nodes["temperature"], _ = node.New("temperature", "temperature", "u1", "devices/")
	if !device.Match(opts) {
		t.Error("device should match the requested search")
	}

}

func TestDevice_Match3(t *testing.T) {
	searchTerms := map[string]string{
		"type":   "temperature",
		"online": "true",
	}
	opts := search.Options{Terms: searchTerms}
	device := New("u1234", "devices/")
	device.Online = false
	if device.Match(opts) {
		t.Error("device should not match the requested search")
	}
	device.Online = true
	if device.Match(opts) {
		t.Error("device should not match the requested search")
	}
	device.Nodes["temperature"], _ = node.New("temperature", "temperature", "u1", "devices/")
	device.Online = false
	if device.Match(opts) {
		t.Error("device should not match the requested search")
	}
	device.Online = true
	if !device.Match(opts) {
		t.Error("device should match the requested search")
	}

}
