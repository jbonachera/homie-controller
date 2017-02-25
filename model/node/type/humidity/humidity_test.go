package humidity

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
	humidity := HumidityNode{"humidity", "devices/", "humidity", "%", 21, "living", time.Time{}, "u1"}
	humidity.MessageHandler(message)
	if humidity.Percent != 23.0 {
		t.Error("setting humidity via MQTTHandler failed: wanted 23, got", humidity.Percent)
	}
}

func TestGetPoint(t *testing.T) {
	humidity := HumidityNode{"humidity", "devices/", "humidity", "%", 21.0, "living", time.Time{}, "u1"}
	point := humidity.GetPoint()
	fields := point.Fields

	if fields["percent"] != 21.0 {
		t.Error("Invalid point retrieved: wanted 21.0, got", fields["percent"])
	}
}
