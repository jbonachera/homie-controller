package messaging

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"time"
	"github.com/jbonachera/homie-controller/model/homieMessage"
)

var c MQTT.Client
var baseTopic string
type CallbackHandler func(message homieMessage.HomieMessage)
var connected bool = false

func Start(broker string, client_id string, mqttBaseTopic string) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + broker + ":1883")
	opts.SetClientID(client_id)
	c = MQTT.NewClient(opts)
	baseTopic = mqttBaseTopic
	connected = false
	for !connected {
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			log.Error("could not connect to MQTT.")
			time.Sleep(5 * time.Second)
		} else {
			log.Info("connected to MQTT broker")
			connected = true
		}
	}
}
func asyncDelSubscription(topic string){
	for !connected {
		log.Info("waiting for MQTT connection to start..")
		time.Sleep(2 * time.Second)
	}
	log.Debug("Unsubscribing to "+topic)
	c.Unsubscribe(topic)
}
func DelSubscription(topic string){
	go asyncDelSubscription(topic)
}

func AddSubscription(topic string, qos byte, callback MQTT.MessageHandler) {
	go asyncAddSubscription(topic, qos, callback)
}
func asyncAddSubscription(topic string, qos byte, callback MQTT.MessageHandler){
	for !connected {
		log.Info("waiting for MQTT connection to start..")
		time.Sleep(2 * time.Second)
	}
	log.Debug("Subscribing to "+topic)
	c.Subscribe(topic, qos, callback)
}

func AddHandler(topic string, callback CallbackHandler){
	AddSubscription(topic, 0, func(mqttClient MQTT.Client, mqttMessage MQTT.Message){
		message, err := homieMessage.New(mqttMessage, baseTopic)
		if err != nil {
			log.Error("malformed message")
		} else {
			callback(message)
		}
	})
}