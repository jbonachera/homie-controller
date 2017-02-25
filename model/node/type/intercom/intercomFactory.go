package intercom

import (
	"github.com/jbonachera/homie-controller/model/node"
)

type IntercomFactory struct{}

func (IntercomFactory) New(name string, parent string, baseTopic string) node.Type {
	return &IntercomNode{name, baseTopic, "intercom", "bool", false, "NoRoom"}
}

func init() {
	node.RegisterNodeTypeFactory("intercom", IntercomFactory{})
}
