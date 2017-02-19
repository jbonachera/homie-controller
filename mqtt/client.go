package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"time"
	"github.com/jbonachera/homie-controller/model/device"
	"github.com/jbonachera/homie-controller/model/homieMessage"
)

var baseTopic string = "devices/"
var c MQTT.Client

type CallbackHandler func(message homieMessage.HomieMessage)

func Start(broker string, client_id string) {
	device.NewRegistry(baseTopic)
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + broker + ":1883")
	opts.SetClientID(client_id)
	c = MQTT.NewClient(opts)
	connected := false
	for !connected {
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			log.Error("could not connect to MQTT.")
			time.Sleep(5 * time.Second)
		} else {
			log.Info("connected to MQTT broker")
			connected = true
		}
	}
	AddSubscription("devices/+/$online", 1, device.OnlineCallback)
}

func AddSubscription(topic string, qos byte, callback MQTT.MessageHandler){
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