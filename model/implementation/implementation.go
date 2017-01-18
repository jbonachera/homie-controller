package implementation

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"

	"errors"
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

func New(name string, baseTopic string) (Implementation, error) {
	if _, exist := implementations[name]; exist {
		return implementations[name].New(baseTopic), nil
	}
	return nil, errors.New("Invalid type requested: " + name)
}