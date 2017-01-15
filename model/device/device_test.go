package device

import (
	_ "github.com/jbonachera/homie-controller/model/device/nodetype/temperature"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"testing"
)

type MessageMock struct {
	topic   string
	payload string
}

func (m MessageMock) Topic() string {
	return m.topic
}
func (m MessageMock) Payload() []byte {
	return []byte(m.payload)
}

var message MessageMock = MessageMock{
	"devices/u1234/temperature/$type",
	"temperature",
}

type dummyClient struct{}

func (d dummyClient) Subscribe(topic string, qos byte, callback func(homieMessage.SubscriptibleClient, homieMessage.HomieExtractableMessage)) interface{} {
	return nil
}

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
	device.MQTTNodeHandler(dummyClient{}, message)
	if len(device.Nodes) != 1 {
		t.Error("adding node failed")
	}
}
