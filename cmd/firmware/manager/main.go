package main

import (
	"github.com/jbonachera/homie-controller/config"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/messaging"
	_ "github.com/jbonachera/homie-controller/ota/github_releases"
	"os"
	"os/signal"

	"github.com/jbonachera/homie-controller/ota"
)

var VERSION = "0.0.1"
var baseTopic string = "devices/"

func main() {
	log.SetLogLevel(config.Get("LOG_LEVEL"))
	log.Info("starting homie-firmware-manager version " + VERSION)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)
	broker := config.Get("mqtt_broker")
	if broker != "" {
		log.Debug("starting connection to MQTT broker at " + broker)
		messaging.Start(broker, config.Get("mqtt_client_id")+"-"+VERSION, baseTopic)
	}
	ota.Start()

	ota.AddFirmware("intercom", "github_release")
	ota.AddFirmware("vx-temperature-sensor", "github_release")
	ota.AddFirmware("temperature-sensor", "github_release")
	log.Debug("main process initialiszed")

	select {
	case <-sigc:
		log.Debug("received interrupt - aborting operations")
		ota.Stop()
		messaging.Stop()
		break
	}
	log.Debug("main process finished")
}
