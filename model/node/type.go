package node

import (
	"github.com/jbonachera/homie-controller/model/metric"
	"github.com/jbonachera/homie-controller/model/homieMessage"
)

type Type interface {
	GetName() string
	GetType() string
	GetProperties() []string
	GetPoint() *metric.Metric
	MessageHandler(message homieMessage.HomieMessage)
}
