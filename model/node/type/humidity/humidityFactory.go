package humidity

import (
	"github.com/jbonachera/homie-controller/model/node"
)

type HumidityFactory struct{}

func (HumidityFactory) New(name string, baseTopic string) node.Type {
	return &HumidityNode{name, baseTopic, "humidity", "%", 0, ""}
}

func init() {
	node.RegisterNodeTypeFactory("humidity", HumidityFactory{})
}
