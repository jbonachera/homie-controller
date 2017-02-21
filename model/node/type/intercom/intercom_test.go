package intercom

import (
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"testing"
)

var message homieMessage.HomieMessage = homieMessage.HomieMessage{
	"devices/u1234/intercom/ringing",
	"true",
	"devices/",
	"u1234",
	"intercom/ringing",
}

func TestMQTTHandler(t *testing.T) {
	intercom := IntercomNode{"intercom", "devices/", "intercom", "bool", false, "living"}
	intercom.MessageHandler(message)
	if !intercom.Ringing {
		t.Error("setting intercom via MQTTHandler failed: wanted true, got", intercom.Ringing)
	}
}

func TestGetPoint(t *testing.T) {
	intercom := IntercomNode{"intercom", "devices/", "intercom", "bool", false, "living"}
	point := intercom.GetPoint()

	if point != nil {
		t.Error("Invalid point retrieved: wanted nil, got", point)
	}
}
