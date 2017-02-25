package humidity

import (
	"github.com/jbonachera/homie-controller/influxdb"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"strconv"
	"strings"
	"time"
)

var nodeType string = "humidity"

type HumidityNode struct {
	name      string
	baseTopic string
	Nodetype  string  `json:"type"`
	Unit      string  `json:"unit"`
	Percent   float64 `json:"percent"`
	Room      string  `json:"room"`
	LastUpdate time.Time `json:"last_update"`
	ParentId string `json:"-"`
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
func (t HumidityNode) GetPoint() *metric.Metric {
	return metric.New("humidity", map[string]string{"room": t.Room, "sensor": t.name}, map[string]interface{}{"percent": t.Percent})
}
func (t *HumidityNode) MessageHandler(message homieMessage.HomieMessage) {
	topicComponents := strings.Split(message.Path, "/")
	node, property := topicComponents[0], topicComponents[1]
	if node != t.name {
		log.Debug("received message for " + node + " but we are " + t.name)
		return
	}
	switch property {
	case "unit":
		t.Unit = message.Payload
	case "room":
		t.Room = message.Payload
	case "percent":
		percent, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.Percent = percent
			t.LastUpdate = time.Now()
			if influxdb.Ready() {
				influxdb.PublishMetric(t.GetPoint())
			}
		}
	}

}
