package implementation

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"

)

var implementations = make(map[string]Factory)

type Implementation interface {
	GetName() string
	Set(property string, value string)
	GetProperties() []string
	MQTTHandler(mqttClient MQTT.Client, message MQTT.Message)
}

type Factory interface {
	New(baseTopic string) Implementation
}

func RegisterImplementation(name string, constructor Factory) {
	_, exist := implementations[name]
	if !exist {
		implementations[name] = constructor
	}
}
