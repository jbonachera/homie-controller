package main

import (
	"github.com/jbonachera/homie-controller/log"
	_ "github.com/jbonachera/homie-controller/model/node/type/humidity"
	_ "github.com/jbonachera/homie-controller/model/node/type/temperature"
	_ "github.com/jbonachera/homie-controller/model/implementation/esp8266"
	"github.com/jbonachera/homie-controller/mqtt"
	"github.com/jbonachera/homie-controller/http"
	"os"
	"os/signal"
	"github.com/jbonachera/homie-controller/config"
	"github.com/jbonachera/homie-controller/influxdb"
	"github.com/influxdata/influxdb/client/v2"

)
var VERSION = "0.1"

func main() {
	log.SetLogLevel("DEBUG")
	log.Info("starting homie-controller version "+VERSION)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	broker := config.Get("mqtt_broker")
	influxdb_server := config.Get("influxdb_server")
	if influxdb_server != "" {
		log.Debug("Starting log export to influxdb at " + influxdb_server)
		influxdb.Start(client.HTTPConfig{
			Addr: "http://"+influxdb_server+":8086",
		}, true)
	}
	if broker != "" {
		log.Debug("starting connection to MQTT broker at "+broker)
		go mqtt.Start(broker, config.Get("mqtt_client_id"))
	}
	go http.Start()
	<-sigc
}
