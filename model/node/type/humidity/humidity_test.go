package humidity

import (
	"github.com/jbonachera/homie-controller/mocks/mqtt"
	"testing"
)

var message mqtt.MessageMock = mqtt.NewMessage(
	"devices/u1234/humidity/percent",
	"23.0",
)

func TestMQTTHandler(t *testing.T) {
	humidity := HumidityNode{"humidity", "devices/", "%", 21, "living"}
	humidity.MQTTHandler(nil, message)
	if humidity.Percent != 23.0 {
		t.Error("setting humidity via MQTTHandler failed: wanted 23, got", humidity.Percent)
	}
}

func TestGetPoint(t *testing.T) {
	humidity := HumidityNode{"humidity", "devices/", "%", 21.0, "living"}
	point := humidity.GetPoint()
	fields, err := point.Fields()
	if err != nil{
		t.Error("error while retrieving Point: ", err)
	}
	if fields["percent"] != 21.0 {
		t.Error("Invalid point retrieved: wanted 21.0, got", fields["percent"])
	}
}
