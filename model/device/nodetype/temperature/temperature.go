package temperature

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"strconv"
	"strings"
)

var nodeType string = "temperature"

type TemperatureNode struct {
	name      string
	baseTopic string
	unit      string
	degrees   float64
	room      string
}

func (t TemperatureNode) GetName() string {
	return t.name
}
func (t TemperatureNode) GetType() string {
	return nodeType
}
func (t TemperatureNode) GetProperties() []string {
	return []string{"degres", "unit", "room"}
}
func (t TemperatureNode) GetPoint() metric.Metric {
	return metric.New("temperature", map[string]string{"room": t.room, "sensor": t.name}, map[string]interface{}{"degrees": t.degrees})
}
func (t *TemperatureNode) MQTTHandler(mqttClient MQTT.Client, mqttMessage MQTT.Message) {
	message, err := homieMessage.New(mqttMessage, t.baseTopic)
	if err != nil {
		return
	}
	topicComponents := strings.Split(message.Path, "/")
	node, property := topicComponents[0], topicComponents[1]
	if node != t.name {
		// Message was not for us
		return
	}
	switch property {
	case "unit":
		t.unit = message.Payload
	case "room":
		t.room = message.Payload
	case "degrees":
		degrees, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.degrees = degrees
		}
	}

}
