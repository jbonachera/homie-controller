package node

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
func (d dummyNodeType) GetProperties() []string {
	return []string{"temperature"}
}
func (d dummyNodeType) GetPoint() *metric.Metric {
	return metric.New("temperature", map[string]string{"room": "living", "sensor": "sensor01"}, map[string]interface{}{"degrees": 24.0})
}

type dummyClient struct{}

func (d dummyClient) Subscribe(topic string, qos byte, callback interface{}) interface{} {
	return nil
}

func (d dummyNodeType) MessageHandler(message homieMessage.HomieMessage) {
}

type dummyFactory struct{}

func (dummyFactory) New(name string, parent string, baseTopic string) Type {
	return dummyNodeType{name, "dummyType"}
}

func TestRegisterNodeTypeFactory(t *testing.T) {
	RegisterNodeTypeFactory("dummy", dummyFactory{})
}

func TestNew(t *testing.T) {
	RegisterNodeTypeFactory("dummy", dummyFactory{})
	node, err := New("dummy", "dummy", "u1", "devices/")
	if err != nil {
		t.Error("error occured:", err)
	}
	if node.GetName() != "dummy" {
		t.Error("could not create a node")
	}
}

func TestNew2(t *testing.T) {
	RegisterNodeTypeFactory("dummy", dummyFactory{})
	node, err := New("dummy", "dummy2", "u1", "devices/")
	if err != nil {
		t.Error("error occured:", err)
	}
	if node.GetName() != "dummy2" {
		t.Error("could not create a node with a different name from the type")
	}
}
func TestAbsentNodeType(t *testing.T) {
	_, err := New("invalid", "invalid", "u1", "devices/")
	if err == nil {
		t.Error("could create an invalid node")
	}
}
