package esp8266

import (
	"github.com/jbonachera/homie-controller/mocks/mqtt"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"testing"
)

func TestEsp8266_GetName(t *testing.T) {
	esp := New("u1", "devices/")
	if esp.GetName() != "esp8266" {
		t.Error("did create an invalid esp8266")
	}
}
func TestEsp8266_Set(t *testing.T) {
	esp := New("u1", "devices/")
	esp.Set("ota/enabled", "true")
	if !esp.Ota {
		t.Error("setting Ota property failed: wanted true, go", esp.Ota)
	}
	esp.Set("version", "cookiebotOS v902.3-testing")
	if esp.Version != "cookiebotOS v902.3-testing" {
		t.Error("setting Version failed: wanted 'cookiebotOS v902.3-testing', got", esp.Version)
	}
}

func TestEsp8266_GetProperties(t *testing.T) {
	esp := New("u1", "devices/")
	if len(esp.GetProperties()) != 2 {
		t.Error("wrong list of properties returned")
	}
}

func TestMQTTHandler(t *testing.T) {
	var message homieMessage.HomieMessage
	message, _ = homieMessage.Extract(mqtt.NewMessage(
		"devices/u1234/implementation/version",
		"cookiebotOS_3.1",
	), "devices/")
	esp := New("u1", "devices/")
	esp.MessageHandler(message)
	if esp.Version != "cookiebotOS_3.1" {
		t.Error("setting Version failed: wanted 'cookiebotOS_3.1', got", esp.Version)
	}
}
func TestEsp8266_Reset(t *testing.T) {
	esp := New("u1", "devices/")
	messages := []homieMessage.HomieMessage{}
	esp.MessagePublisher = func(message homieMessage.HomieMessage) {
		messages = append(messages, message)
	}
	esp.Reset()
	if messages[0].Payload != "true" {
		t.Error("message payload was wrong: got " + messages[0].Payload + ", wanted true")
	}
	topic := "devices/u1/$implementation/reset"
	if messages[0].Topic != topic {
		t.Error("message payload was wrong: got " + messages[0].Topic + ", wanted " + topic)
	}
}
