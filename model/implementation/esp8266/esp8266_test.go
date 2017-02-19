package esp8266

import (
	"testing"
	"github.com/jbonachera/homie-controller/mocks/mqtt"
)

func TestEsp8266_GetName(t *testing.T) {
	esp := New("devices/")
	if esp.GetName() != "esp8266"{
		t.Error("did create an invalid esp8266")
	}
}
func TestEsp8266_Set(t *testing.T) {
	esp := New("devices/")
	esp.Set("ota", "true")
	if !esp.ota {
		t.Error("setting ota property failed: wanted true, go", esp.ota)
	}
	esp.Set("version", "cookiebotOS v902.3-testing")
	if esp.version != "cookiebotOS v902.3-testing"{
		t.Error("setting version failed: wanted 'cookiebotOS v902.3-testing', got", esp.version)
	}
}

func TestEsp8266_GetProperties(t *testing.T) {
	esp := New("devices/")
	if len(esp.GetProperties()) != 2 {
		t.Error("wrong list of properties returned")
	}
}

func TestMQTTHandler(t *testing.T) {
	message := mqtt.NewMessage(
		"devices/u1234/implementation/version",
		"cookiebotOS_3.1",
	)
	client := mqtt.NewMockClient(true, "devices/")
	esp := New("devices/")
	esp.MQTTHandler(client, message)
	if esp.version != "cookiebotOS_3.1"{
		t.Error("setting version failed: wanted 'cookiebotOS_3.1', got", esp.version)
	}
}