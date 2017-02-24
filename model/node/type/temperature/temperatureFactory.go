package temperature

import (
	"github.com/jbonachera/homie-controller/model/node"
	"time"
)

type TemperatureFactory struct{}

func (TemperatureFactory) New(name string, baseTopic string) node.Type {
	return &TemperatureNode{name, baseTopic, "temperature", "c", 0.0, "",time.Time{}}
}

func init() {
	node.RegisterNodeTypeFactory("temperature", TemperatureFactory{})
}
