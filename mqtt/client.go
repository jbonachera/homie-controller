package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/device"
	"time"
)

var baseTopic string = "devices/"
var c MQTT.Client

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