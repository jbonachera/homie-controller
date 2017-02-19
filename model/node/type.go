package node

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/model/metric"
)

type Type interface {
	GetName() string
	GetType() string
	GetProperties() []string
	GetPoint() *metric.Metric
	MQTTHandler(mqttClient MQTT.Client, message MQTT.Message)
}
