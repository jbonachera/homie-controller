package system

import (
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"testing"
	"time"
)

var message homieMessage.HomieMessage = homieMessage.HomieMessage{
	"devices/u1234/humidity/percent",
	"23.0",
	"devices/",
	"u1234",
	"humidity/percent",
}

func TestMQTTHandler(t *testing.T) {
	system := SystemNode{"system", "devices/", "system", 6557.0, time.Time{}, "u1"}
	system.MessageHandler(message)
	if system.Vcc != 6557.0 {
		t.Error("setting VCC via MQTTHandler failed: wanted 6557, got", system.Vcc)
	}
}

func TestGetPoint(t *testing.T) {
	system := SystemNode{"system", "devices/", "system", 6557.0, time.Time{}, "u1"}
	point := system.GetPoint()
	fields := point.Fields

	if fields["vcc"] != 6557.0 {
		t.Error("Invalid point retrieved: wanted 6557, got", fields["vcc"])
	}
}
