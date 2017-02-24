package humidity

import (
	"github.com/jbonachera/homie-controller/model/node"
	"time"
)

type HumidityFactory struct{}

func (HumidityFactory) New(name string, baseTopic string) node.Type {
	return &HumidityNode{name, baseTopic, "humidity", "%", 0, "", time.Time{}}
}

func init() {
	node.RegisterNodeTypeFactory("humidity", HumidityFactory{})
}
