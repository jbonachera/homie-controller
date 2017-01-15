package temperature

import (
	"github.com/jbonachera/homie-controller/model/device/nodetype"
)

type TemperatureFactory struct{}

func (TemperatureFactory) New(name string, baseTopic string) nodetype.NodeType {
	return &TemperatureNode{name, baseTopic, "c", 0.0, ""}
}

func init() {
	nodetype.RegisterNodeTypeFactory("temperature", TemperatureFactory{})
}
