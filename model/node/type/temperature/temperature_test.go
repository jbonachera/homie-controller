package temperature

import (
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"testing"
	"time"
)

var message homieMessage.HomieMessage = homieMessage.HomieMessage{
	"devices/u1234/temperature/degrees",
	"23.9",
	"devices/",
	"u1234",
	"temperature/degrees",
}

func TestMQTTHandler(t *testing.T) {
	temperature := TemperatureNode{"temperature", "devices/", "temperature", "c", 21.0, "living", time.Time{}}
	temperature.MessageHandler(message)
	if temperature.Degrees != 23.9 {
		t.Error("setting temperature via MessageHandler failed: wanted 23.9, got", temperature.Degrees)
	}
}

func TestGetPoint(t *testing.T) {
	temperature := TemperatureNode{"temperature", "devices/", "temperature", "c", 21.0, "living", time.Time{}}
	point := temperature.GetPoint()
	fields := point.Fields
	if fields["degrees"] != 21.0 {
		t.Error("Invalid point retrieved: wanted 21.0, got", fields["degrees"])
	}
}
