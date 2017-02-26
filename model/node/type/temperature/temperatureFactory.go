package temperature

import (
	"github.com/jbonachera/homie-controller/model/node"
	"time"
)

type TemperatureFactory struct{}

func (TemperatureFactory) New(name string, parent string, baseTopic string) node.Type {
	return &TemperatureNode{name, baseTopic, "temperature", "c", 0.0, "", time.Time{}, parent}
}

func init() {
	node.RegisterNodeTypeFactory("temperature", TemperatureFactory{})
}
