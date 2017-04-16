package system

import (
	"github.com/jbonachera/homie-controller/model/node"
	"time"
)

type systemFactory struct{}

func (systemFactory) New(name string, parent string, baseTopic string) node.Type {
	return &SystemNode{name, baseTopic, "system", 0, time.Time{}, parent}
}

func init() {
	node.RegisterNodeTypeFactory("system", systemFactory{})
}
