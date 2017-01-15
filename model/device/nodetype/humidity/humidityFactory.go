package humidity

import (
	"github.com/jbonachera/homie-controller/model/device/nodetype"
)

type HumidityFactory struct{}

func (HumidityFactory) New(name string, baseTopic string) nodetype.NodeType {
	return &HumidityNode{name, baseTopic, "%", 0, ""}
}

func init() {
	nodetype.RegisterNodeTypeFactory("humidity", HumidityFactory{})
}
