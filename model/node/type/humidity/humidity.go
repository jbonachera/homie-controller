package humidity

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"strconv"
	"strings"
	"github.com/jbonachera/homie-controller/influxdb"
	"github.com/influxdata/influxdb/client/v2"
)

var nodeType string = "humidity"

type HumidityNode struct {
	name      string `json:"name"`
	baseTopic string `json:"base_topic"`
	Nodetype  string `json:"type"`
	Unit      string `json:"unit"`
	Percent   float64 `json:"percent"`
	Room      string `json:"room"`
}

func (t HumidityNode) GetName() string {
	return t.name
}
func (t HumidityNode) GetType() string {
	return nodeType
}
func (t HumidityNode) GetProperties() []string {
	return []string{"percent", "unit", "room"}
}
func (t HumidityNode) GetPoint() *client.Point {
	return metric.New("humidity", map[string]string{"room": t.Room, "sensor": t.name}, map[string]interface{}{"percent": t.Percent})
}
func (t *HumidityNode) MQTTHandler(mqttClient MQTT.Client, mqttMessage MQTT.Message) {
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
	log.Debug("set property " + property + " to " + message.Payload + " for node " + t.name)
	switch property {
	case "unit":
		t.Unit = message.Payload
	case "room":
		t.Room = message.Payload
	case "percent":
		percent, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.Percent = percent
			if influxdb.Ready() {
				influxdb.PublishPoint(t.GetPoint())
			}
		}
	}

}
