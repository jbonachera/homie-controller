package influxdb

import (
	"github.com/jbonachera/homie-controller/log"
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/jbonachera/homie-controller/config"
)

var dbClient influxdb.Client
var ready bool = false
var logOnlyMode = false

func Start(config influxdb.HTTPConfig, logOnly bool){
	var err error
	dbClient, err = influxdb.NewHTTPClient(config)
	if err != nil {
		log.Error("Error: "+ err.Error())
		return
	}
	logOnlyMode = logOnly
	ready = true
}

func Ready() bool{
	return ready
}

func PublishPoint(metric *influxdb.Point){

	if !logOnlyMode {
		bp, _ := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Database:  config.Get("INFLUXDB_DATABASE"),
			Precision: "us",
		})
		bp.AddPoint(metric)
		err := dbClient.Write(bp)
		if err != nil {
			log.Error("Error: " + err.Error())
		} else {
			log.Debug("metrics sent to influxdb server")
		}
	} else {
		log.Info("would have send metric to InfluxDB: " + metric.Name())
	}
}