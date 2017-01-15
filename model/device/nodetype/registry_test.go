package nodetype

import (
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"

	"testing"
)

type dummyNodeType struct {
	name     string
	nodeType string
}

func (d dummyNodeType) GetName() string {
	return d.name
}
func (d dummyNodeType) GetType() string {
	return d.nodeType
}
func (d dummyNodeType) GetPoint() metric.Metric {
	return metric.New("temperature", map[string]string{"room": "living", "sensor": "sensor01"}, map[string]interface{}{"degrees": 24.0})
}
func (d dummyNodeType) MQTTHandler(mqttClient interface{}, message homieMessage.HomieExtractableMessage) {
}

type dummyFactory struct{}

func (dummyFactory) New(name string, baseTopic string) NodeType {
	return dummyNodeType{name, "dummyType"}
}

func TestRegisterNodeTypeFactory(t *testing.T) {
	RegisterNodeTypeFactory("dummy", dummyFactory{})
}

func TestNew(t *testing.T) {
	RegisterNodeTypeFactory("dummy", dummyFactory{})
	node := New("dummy", "devices/")
	if node.GetName() != "dummy" {
		t.Error("could not create a node")
	}
}
