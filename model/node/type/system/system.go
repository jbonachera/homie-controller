package system

import (
	"github.com/jbonachera/homie-controller/influxdb"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"strconv"
	"strings"
	"time"
)

var nodeType string = "system"

type SystemNode struct {
	name       string
	baseTopic  string
	Nodetype   string    `json:"type"`
	Vcc        float64   `json:"vcc"`
	LastUpdate time.Time `json:"last_update"`
	ParentId   string    `json:"-"`
}

func (t SystemNode) GetName() string {
	return t.name
}
func (t SystemNode) GetType() string {
	return nodeType
}
func (t SystemNode) GetProperties() []string {
	return []string{"vcc"}
}
func (t SystemNode) GetPoint() *metric.Metric {
	return metric.New("system", map[string]string{"sensor": t.name}, map[string]interface{}{"vcc": t.Vcc})
}
func (t *SystemNode) MessageHandler(message homieMessage.HomieMessage) {
	topicComponents := strings.Split(message.Path, "/")
	node, property := topicComponents[0], topicComponents[1]
	if node != t.name {
		log.Debug("received message for " + node + " but we are " + t.name)
		return
	}
	switch property {
	case "vcc":
		vcc, err := strconv.ParseFloat(message.Payload, 64)
		if err == nil {
			t.Vcc = vcc
			t.LastUpdate = time.Now()
			if influxdb.Ready() {
				influxdb.PublishMetric(t.GetPoint())
			}
		}
	}

}
