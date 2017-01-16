package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/device"
	"os"
	"time"
)

var baseTopic string = "devices/"
var c MQTT.Client

func Start(broker string) {
	device.NewRegistry(baseTopic)
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + broker + ":1883")
	opts.SetClientID("homie-controller")
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
	if token := c.Subscribe("devices/+/$online", 1, device.OnlineCallback); token.Wait() && token.Error() != nil {
		log.Error("Could not subscribe to devices/+/$online")
		os.Exit(1)
	}

}