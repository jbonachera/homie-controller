package node

import (
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
)

type Type interface {
	GetName() string
	GetType() string
	GetProperties() []string
	GetPoint() *metric.Metric
	MessageHandler(message homieMessage.HomieMessage)
}
