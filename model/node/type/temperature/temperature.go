package temperature

import (
	"github.com/jbonachera/homie-controller/influxdb"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"strconv"
	"strings"
	"time"
)

var nodeType string = "temperature"

type TemperatureNode struct {
	name      string
	baseTopic string
	Nodetype  string  `json:"type"`
	Unit      string  `json:"unit"`
	Degrees   float64 `json:"degrees"`
	Room      string  `json:"room"`
	LastUpdate time.Time `json:"last_update"`
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
func (t TemperatureNode) GetPoint() *metric.Metric {
	return metric.New(t.name, map[string]string{"room": t.Room, "sensor": t.name}, map[string]interface{}{"degrees": t.Degrees})
}
func (t *TemperatureNode) MessageHandler(message homieMessage.HomieMessage) {
	topicComponents := strings.Split(message.Path, "/")
	node, property := topicComponents[0], topicComponents[1]
	if node != t.name {
		log.Debug("received message for " + node + " but we are " + t.name)
		return
	}
	log.Debug("set property " + property + " to " + message.Payload + " for node " + t.name)
	switch property {
	case "unit":
		t.Unit = message.Payload
	case "room":
		t.Room = message.Payload
	case "degrees":
		degrees, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.Degrees = degrees
			t.LastUpdate = time.Now()
			if influxdb.Ready() {
				influxdb.PublishMetric(t.GetPoint())
			}
		}
	}

}
