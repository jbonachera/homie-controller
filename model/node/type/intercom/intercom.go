package intercom

import (
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"strings"
	"github.com/jbonachera/homie-controller/log"
)

var nodeType string = "intercom"

type IntercomNode struct {
	name      string
	baseTopic string
	Nodetype  string  `json:"type"`
	Unit      string  `json:"unit"`
	Ringing   bool `json:"ringing"`
	Room      string  `json:"room"`
	ParentId string `json:"-"`
}

func (t IntercomNode) GetName() string {
	return t.name
}
func (t IntercomNode) GetType() string {
	return nodeType
}
func (t IntercomNode) GetProperties() []string {
	return []string{"ringing", "unit", "room"}
}
func (t IntercomNode) GetPoint() *metric.Metric {
	return nil
}
func (t *IntercomNode) MessageHandler(message homieMessage.HomieMessage) {
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
	case "ringing":
		t.Ringing = message.Payload == "true"
	}

}
