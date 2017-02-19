package mqtt

import (
	"testing"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/mocks/mqtt"
)

func callback(_ MQTT.Client, _ MQTT.Message) {}


func TestAddSubscription(t *testing.T) {
	mock := mqtt.NewMockClient(true, "old/topic")
	c = mock
	AddSubscription("devices/bah/+", 0, callback)
	if mock.Topic != "devices/bah/+" {
		t.Error("could not add subscription")
	}
}

func TestAddHandler(t *testing.T) {
	mock := mqtt.NewMockClient(true, "old/topic")
	c = mock
	AddHandler("devices/bah/+", func(topic string, payload string){})
	if mock.Topic != "devices/bah/+" {
		t.Error("could not add handler")
	}

}