package nodetype

import (
	"errors"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
)

var nodeTypeFactories = make(map[string]NodeTypeFactory)

type NodeType interface {
	GetName() string
	GetType() string
	GetPoint() metric.Metric
	MQTTHandler(mqttClient homieMessage.SubscriptibleClient, message homieMessage.HomieExtractableMessage)
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

func New(nodeType string, baseTopic string) (NodeType, error) {
	if _, exist := nodeTypeFactories[nodeType]; exist {
		node := nodeTypeFactories[nodeType].New(nodeType, baseTopic)
		return node, nil
	}
	return nil, errors.New("Invalid type requested")
}
