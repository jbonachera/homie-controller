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
	MQTTHandler(mqttClient interface{}, message homieMessage.HomieExtractableMessage)
}

type NodeTypeFactory interface {
	New(name string, baseTopic string) NodeType
}

func RegisterNodeTypeFactory(name string, nodeType NodeTypeFactory) {
	_, exist := nodeTypeFactories[name]
	if !exist {
		nodeTypeFactories[name] = nodeType
	}
}

func New(nodeType string, baseTopic string) NodeType {
	return nodeTypeFactories[nodeType].New(nodeType, baseTopic)
}
