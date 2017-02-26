package messaging

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"time"
)

type CallbackHandler func(message homieMessage.HomieMessage)

type messagingBroker struct {
	c         MQTT.Client
	baseTopic string
	connected bool
}

var broker messagingBroker

func Start(brokerHost string, client_id string, mqttBaseTopic string) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + brokerHost + ":1883")
	opts.SetClientID(client_id)
	broker = messagingBroker{MQTT.NewClient(opts), mqttBaseTopic, false}
	for !broker.connected {
		if token := broker.c.Connect(); token.Wait() && token.Error() != nil {
			log.Error("could not connect to MQTT.")
			time.Sleep(5 * time.Second)
		} else {
			log.Info("connected to MQTT broker")
			broker.connected = true
		}
	}
}

func mqttPublish(topic string, qos byte, retained bool, payload string) {
	broker.c.Publish(topic, qos, retained, payload)
}

func PublishMessage(message homieMessage.HomieMessage) {
	mqttPublish(message.Topic, 1, false, message.Payload)
}
func PublishState(message homieMessage.HomieMessage) {
	mqttPublish(message.Topic, 1, true, message.Payload)
}
func DelSubscription(topic string) {
	for !broker.connected {
		log.Info("waiting for MQTT connection to start..")
		time.Sleep(2 * time.Second)
	}
	log.Debug("Unsubscribing to " + topic)
	broker.c.Unsubscribe(topic)
}
func AddSubscription(topic string, qos byte, callback MQTT.MessageHandler) {
	for !broker.connected {
		log.Info("waiting for MQTT connection to start..")
		time.Sleep(2 * time.Second)
	}
	log.Debug("Subscribing to " + topic)
	broker.c.Subscribe(topic, qos, callback)
}

func AddHandler(topic string, callback CallbackHandler) {
	AddSubscription(topic, 0, func(mqttClient MQTT.Client, mqttMessage MQTT.Message) {
		message, err := homieMessage.Extract(mqttMessage, broker.baseTopic)
		if err != nil {
			log.Error("malformed message")
		} else {
			callback(message)
		}
	})
}
