package nodetype

import (
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
)

var nodeTypeFactories = make(map[string]NodeTypeFactory)

type NodeType interface {
	GetName() string
	GetType() string
	GetPoint() metric.Metric
	MQTTHandler(message homieMessage.HomieMessage)
}

type NodeTypeFactory interface {
	New() NodeType
}

func RegisterNodeTypeFactory(name string, nodeType NodeTypeFactory) {
	_, exist := nodeTypeFactories[name]
	if !exist {
		nodeTypeFactories[name] = nodeType
	}
}

func New(nodeType string) NodeType {
	return nodeTypeFactories[nodeType].New()
}
