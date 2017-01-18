package temperature

import (
	"github.com/jbonachera/homie-controller/mocks/mqtt"
	"testing"
)

var message mqtt.MessageMock = mqtt.NewMessage(
	"devices/u1234/temperature/degrees",
	"23.9",
)

func TestMQTTHandler(t *testing.T) {
	temperature := TemperatureNode{"temperature", "devices/", "temperature", "c", 21.0, "living"}
	temperature.MQTTHandler(nil, message)
	if temperature.Degrees != 23.9 {
		t.Error("setting temperature via MQTTHandler failed: wanted 23.9, got", temperature.Degrees)
	}
}

func TestGetPoint(t *testing.T) {
	temperature := TemperatureNode{"temperature", "devices/", "temperature", "c", 21.0, "living"}
	point := temperature.GetPoint()
	fields, err := point.Fields()
	if err != nil{
		t.Error("error while retrieving Point: ", err)
	}
	if fields["degrees"] != 21.0 {
		t.Error("Invalid point retrieved: wanted 21.0, got", fields["degrees"])
	}
}
