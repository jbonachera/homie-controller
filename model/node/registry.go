package node

import (
	"errors"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/model/metric"
)

var Factories = make(map[string]Factory)

type Type interface {
	GetName() string
	GetType() string
	GetProperties() []string
	GetPoint() metric.Metric
	MQTTHandler(mqttClient MQTT.Client, message MQTT.Message)
}

type Factory interface {
	New(name string, baseTopic string) Type
}

func RegisterNodeTypeFactory(name string, nodeType Factory) {
	_, exist := Factories[name]
	if !exist {
		Factories[name] = nodeType
	}
}

func New(nodeType string, baseTopic string) (Type, error) {
	if _, exist := Factories[nodeType]; exist {
		node := Factories[nodeType].New(nodeType, baseTopic)
		return node, nil
	}
	return nil, errors.New("Invalid type requested: " + nodeType)
}
