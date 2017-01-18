package temperature

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"github.com/jbonachera/homie-controller/influxdb"
	"strconv"
	"strings"
	"github.com/influxdata/influxdb/client/v2"
)

var nodeType string = "temperature"

type TemperatureNode struct {
	name      string `json:"name"`
	baseTopic string `json:"base_topic"`
	Unit      string `json:"unit"`
	Degrees   float64 `json:"degrees"`
	Room      string `json:"room"`
}

func (t TemperatureNode) GetName() string {
	return t.name
}
func (t TemperatureNode) GetType() string {
	return nodeType
}
func (t TemperatureNode) GetProperties() []string {
	return []string{"degrees", "unit", "room"}
}
func (t TemperatureNode) GetPoint() *client.Point {
	return metric.New("temperature", map[string]string{"room": t.Room, "sensor": t.name}, map[string]interface{}{"degrees": t.Degrees})
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
	log.Debug("set property " + property + " to " + message.Payload + " for node" + t.name)
	switch property {
	case "unit":
		t.Unit = message.Payload
	case "room":
		t.Room = message.Payload
	case "degrees":
		degrees, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.Degrees = degrees
			if influxdb.Ready() {
				influxdb.PublishPoint(t.GetPoint())
			}
		}
	}

}
