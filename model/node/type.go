package node

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/influxdata/influxdb/client/v2"
)

type Type interface {
	GetName() string
	GetType() string
	GetProperties() []string
	GetPoint() *client.Point
	MQTTHandler(mqttClient MQTT.Client, message MQTT.Message)
}
