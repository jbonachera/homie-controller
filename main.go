package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jbonachera/homie-controller/config"
	"github.com/jbonachera/homie-controller/http"
	"github.com/jbonachera/homie-controller/influxdb"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/messaging"
	_ "github.com/jbonachera/homie-controller/model/implementation/esp8266"
	_ "github.com/jbonachera/homie-controller/model/node/type/humidity"
	_ "github.com/jbonachera/homie-controller/model/node/type/intercom"
	_ "github.com/jbonachera/homie-controller/model/node/type/system"
	_ "github.com/jbonachera/homie-controller/model/node/type/temperature"
	_ "github.com/jbonachera/homie-controller/model/node/type/weather_sensor"
	_ "github.com/jbonachera/homie-controller/ota/github_releases"
	"os"
	"os/signal"

	"github.com/jbonachera/homie-controller/model/device"
	"github.com/jbonachera/homie-controller/ota"
)

var VERSION = "0.0.2"
var baseTopic string = "devices/"

func main() {
	log.SetLogLevel(config.Get("LOG_LEVEL"))
	config.Set("version", VERSION)
	log.Info("starting homie-controller version " + VERSION)
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

	influxdb_server := config.Get("influxdb_server")
	if influxdb_server != "" {
		log.Debug("Starting log export to influxdb at " + influxdb_server)
		influxdb.Start(client.HTTPConfig{
			Addr: "http://" + influxdb_server + ":8086",
		}, false)
	}

	device.NewRegistry(baseTopic)
	messaging.AddHandler("devices/+/$online", device.OnlineCallback)
	go http.Start()
	select {
	case <-sigc:
		log.Debug("received interrupt - aborting operations")
		ota.Stop()
		messaging.Stop()
		break
	}
	log.Debug("main process finished")
}
