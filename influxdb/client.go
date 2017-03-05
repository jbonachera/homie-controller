package influxdb

import (
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/jbonachera/homie-controller/config"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/metric"
)

var dbClient influxdb.Client
var ready bool = false
var logOnlyMode = false

func Start(config influxdb.HTTPConfig, logOnly bool) {
	var err error
	dbClient, err = influxdb.NewHTTPClient(config)
	if err != nil {
		log.Error("Error: " + err.Error())
		return
	}
	logOnlyMode = logOnly
	ready = true
}

func Ready() bool {
	return ready
}

func metricToPoint(metric *metric.Metric) *influxdb.Point {
	point, err := influxdb.NewPoint(metric.Name, metric.Tags, metric.Fields, metric.Time)
	if err != nil {
		log.Error("could not create point")
		return nil
	} else {
		return point
	}
}

func PublishMetric(metric *metric.Metric) {
	PublishPoint(metricToPoint(metric))
}

func PublishPoint(metric *influxdb.Point) {

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
			log.Debug("metrics sent to influxdb server: " + metric.Name())
		}
	} else {
		log.Info("would have send metric to InfluxDB: " + metric.Name())
	}
}
