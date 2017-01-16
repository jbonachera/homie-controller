package node

import (
	"errors"
)

var Factories = make(map[string]Factory)

type Factory interface {
	New(name string, baseTopic string) Type
}


func RegisterNodeTypeFactory(name string, nodeType Factory) {
	_, exist := Factories[name]
	if !exist {
		Factories[name] = nodeType
	}
}

func New(nodeType string, baseTopic string) (Type, error) {
	if _, exist := Factories[nodeType]; exist {
		node := Factories[nodeType].New(nodeType, baseTopic)
		return node, nil
	}
	return nil, errors.New("Invalid type requested: " + nodeType)
}