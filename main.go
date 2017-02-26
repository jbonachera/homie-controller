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
	_ "github.com/jbonachera/homie-controller/model/node/type/temperature"
	_ "github.com/jbonachera/homie-controller/ota/github_releases"
	"os"
	"os/signal"

	"github.com/jbonachera/homie-controller/model/device"
	"github.com/jbonachera/homie-controller/ota"
)

var VERSION = "0.1"
var baseTopic string = "devices/"

func main() {
	log.SetLogLevel(config.Get("LOG_LEVEL"))
	log.Info("starting homie-controller version " + VERSION)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	ota.AddFirmware("vx-intercom-sensor", "github_release")

	broker := config.Get("mqtt_broker")
	influxdb_server := config.Get("influxdb_server")
	if influxdb_server != "" {
		log.Debug("Starting log export to influxdb at " + influxdb_server)
		influxdb.Start(client.HTTPConfig{
			Addr: "http://" + influxdb_server + ":8086",
		}, false)
	}
	if broker != "" {
		log.Debug("starting connection to MQTT broker at " + broker)
		go messaging.Start(broker, config.Get("mqtt_client_id")+"-"+VERSION, baseTopic)
	}
	device.NewRegistry(baseTopic)
	messaging.AddHandler("devices/+/$online", device.OnlineCallback)
	go http.Start()
	<-sigc
}
