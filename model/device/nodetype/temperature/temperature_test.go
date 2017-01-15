package temperature

import (
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
	"devices/u1234/temperature/degrees",
	"23.9",
}

func TestMQTTHandler(t *testing.T) {
	temperature := TemperatureNode{"temperature", "devices/", "c", 21.0, "living"}
	temperature.MQTTHandler(nil, message)
	if temperature.degrees != 23.9 {
		t.Error("setting temperature via MQTTHandler failed: wanted 23.9, got", temperature.degrees)
	}
}

func TestGetPoint(t *testing.T) {
	temperature := TemperatureNode{"temperature", "devices/", "c", 21.0, "living"}
	point := temperature.GetPoint()
	if point.Fields["degrees"] != 21.0 {
		t.Error("Invalid point retrieved: wanted 21.0, got", point.Fields["degrees"])
	}
}
