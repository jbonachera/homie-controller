package humidity

import (
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"strconv"
	"strings"
)

var nodeType string = "humidity"

type HumidityNode struct {
	name      string
	baseTopic string
	unit      string
	percent   float64
	room      string
}

func (t HumidityNode) GetName() string {
	return t.name
}
func (t HumidityNode) GetType() string {
	return nodeType
}
func (t HumidityNode) GetPoint() metric.Metric {
	return metric.New("humidity", map[string]string{"room": t.room, "sensor": t.name}, map[string]interface{}{"percent": t.percent})
}
func (t *HumidityNode) MQTTHandler(mqttClient homieMessage.SubscriptibleClient, mqttMessage homieMessage.HomieExtractableMessage) {
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
	case "percent":
		percent, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.percent = percent
		}
	}

}