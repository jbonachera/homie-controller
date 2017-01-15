package metric

import (
	"time"
)

type Metric struct {
	Name      string
	Tags      map[string]string
	Fields    map[string]string
	Timestamp time.Time
}

func New(name string, tags map[string]string, fields map[string]string) Metric {
	return Metric{name, tags, fields, time.Now()}
}

type Metricable interface {
	GetPoint() Metric
}
