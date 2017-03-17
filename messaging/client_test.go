package messaging

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/mocks/mqtt"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"testing"
	"time"
)

func callback(_ MQTT.Client, _ MQTT.Message) {}

func TestAddSubscription(t *testing.T) {
	mock := mqtt.NewMockClient(true, "old/topic")
	broker = messagingBroker{mock, "devices/", true}
	startChannels()
	AddSubscription("devices/bah/+", 0, callback)
	time.Sleep(300 * time.Millisecond)
	if mock.Topic != "devices/bah/+" {
		t.Error("could not add subscription")
	}
}

func TestDelSubscription(t *testing.T) {
	mock := mqtt.NewMockClient(true, "old/topic")
	broker = messagingBroker{mock, "devices/", true}
	startChannels()
	AddSubscription("devices/bah/+", 0, callback)
	DelSubscription("devices/bah/+")
	if mock.Topic == "devices/bah/+" {
		t.Error("could not del subscription")
	}
}
func TestAddHandler(t *testing.T) {
	mock := mqtt.NewMockClient(true, "old/topic")
	broker = messagingBroker{mock, "devices/", true}
	startChannels()
	AddHandler("devices/bah/+", func(message homieMessage.HomieMessage) {})
	time.Sleep(10 * time.Millisecond)
	if mock.Topic != "devices/bah/+" {
		t.Error("could not add handler")
	}

}

func TestPublishMessage(t *testing.T) {
	mock := mqtt.NewMockClient(true, "old/topic")
	broker = messagingBroker{mock, "devices/", true}
	startChannels()
	topic := "devices/u1/implementation/reset"
	PublishMessage(homieMessage.HomieMessage{Id: "u1", Payload: "true", BaseTopic: "devices/", Path: "$implementation/reset", Topic: topic})
	time.Sleep(10 * time.Millisecond)
	if mock.PublishedMessage[0].Topic() != topic {
		t.Error("did not published to MQTT broker")
	}
	if mock.PublishedMessage[0].Retained() != false {
		t.Error("MQTT message was flagged for retention, and should not")
	}
}
func TestPublishState(t *testing.T) {
	mock := mqtt.NewMockClient(true, "old/topic")
	broker = messagingBroker{mock, "devices/", true}
	startChannels()
	topic := "devices/u1/implementation/reset"
	PublishState(homieMessage.HomieMessage{Id: "u1", Payload: "true", BaseTopic: "devices/", Path: "$implementation/reset", Topic: topic})
	time.Sleep(10 * time.Millisecond)
	if mock.PublishedMessage[0].Topic() != topic {
		t.Error("did not published to MQTT broker")
	}
	if mock.PublishedMessage[0].Retained() != true {
		t.Error("MQTT message was not flagged for retention, and should be")
	}
}
