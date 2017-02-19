package metric

import (
	"time"
)

type Metric struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
	Time   time.Time
}

func New(name string, tags map[string]string, fields map[string]interface{}) *Metric {
	point := &Metric{name, tags, fields, time.Now()}
	return point
}

type Metricable interface {
	GetPoint() *Metric
}
