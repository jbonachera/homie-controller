package mqtt

import (
	"os"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"time"
	"github.com/jbonachera/homie-controller/model/device"
)
var baseTopic string = "devices/"
var devices device.Registry = device.NewRegistry(baseTopic)
var c MQTT.Client

func Start() {

	opts := MQTT.NewClientOptions().AddBroker("tcp://" + os.Getenv("MQTT_BROKER") + ":1883")
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
	if token := c.Subscribe("devices/+/$online", 0, devices.DeviceOnlineCallback); token.Wait() && token.Error() != nil {
		log.Error("Could not subscribe to devices/+/$online")
		os.Exit(1)
	}

}

func Client() MQTT.Client{
	return c
}