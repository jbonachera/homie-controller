package metric

import (
	"time"
	"github.com/influxdata/influxdb/client/v2"
)

type Metric client.Point

func New(name string, tags map[string]string, fields map[string]interface{}) *client.Point {
	point, _ := client.NewPoint(name, tags, fields, time.Now())
	return point
}

type Metricable interface {
	GetPoint() Metric
}
