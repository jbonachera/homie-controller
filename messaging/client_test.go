package messaging

import (
	"testing"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/mocks/mqtt"
	"github.com/jbonachera/homie-controller/model/homieMessage"
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

func TestDelSubscription(t *testing.T) {
	mock := mqtt.NewMockClient(true, "old/topic")
	c = mock
	AddSubscription("devices/bah/+", 0, callback)
	DelSubscription("devices/bah/+")
	if mock.Topic == "devices/bah/+" {
		t.Error("could not del subscription")
	}
}
func TestAddHandler(t *testing.T) {
	mock := mqtt.NewMockClient(true, "old/topic")
	c = mock
	AddHandler("devices/bah/+", func(message homieMessage.HomieMessage){})
	if mock.Topic != "devices/bah/+" {
		t.Error("could not add handler")
	}

}