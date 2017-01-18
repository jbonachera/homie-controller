package esp8266

import (
	"strconv"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"strings"
	"github.com/jbonachera/homie-controller/log"
)
type esp8266 struct {
	version string
	ota bool
	baseTopic string
}

func New(baseTopic string) *esp8266{
	return &esp8266{"", false, baseTopic}
}

func (e *esp8266) GetName() string{
	return "esp8266"
}

func (e *esp8266) Set(property string, value string){
	switch property {
	case "version":
		e.version = value
	case "ota":
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			e.ota = boolValue
		}
	}
}
func (e *esp8266) GetProperties() []string{
	return []string{"version", "ota"}
}

func (e *esp8266) MQTTHandler(mqttClient MQTT.Client, mqttMessage MQTT.Message){
	// Will be bound to devices/<id>/implementation/#
	message, err := homieMessage.New(mqttMessage, e.baseTopic)
	if err != nil {
		return
	}
	topicComponents := strings.Split(message.Path, "/")
	if len(topicComponents) < 2 {
		return
	}
	log.Debug("received update: " + topicComponents[1] + " to "+message.Payload)
	propertyName := topicComponents[1]
	e.Set(propertyName, message.Payload)
}