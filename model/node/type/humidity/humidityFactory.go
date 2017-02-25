package humidity

import (
	"github.com/jbonachera/homie-controller/model/node"
	"time"
)

type HumidityFactory struct{}

func (HumidityFactory) New(name string, parent string, baseTopic string) node.Type {
	return &HumidityNode{name, baseTopic, "humidity", "%", 0, "", time.Time{}, parent}
}

func init() {
	node.RegisterNodeTypeFactory("humidity", HumidityFactory{})
}
