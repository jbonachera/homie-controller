package device

import (
	"fmt"
	"github.com/jbonachera/homie-controller/model/device/nodetype"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"strconv"
	"strings"
)

type DeviceStats struct {
	Uptime   int
	Signal   int
	Interval int
}

type DeviceFirmware struct {
	Name     string
	Version  string
	Checksum string
}

type Device struct {
	Id             string
	Online         bool
	Name           string
	Localip        string
	Mac            string
	Stats          DeviceStats
	Fw             DeviceFirmware
	Implementation string
	Nodes          map[string]nodetype.NodeType
	BaseTopic      string
}

func New(id string, baseTopic string) Device {
	return Device{id, false, "", "", "", DeviceStats{0, 0, 0}, DeviceFirmware{"", "", ""}, "", map[string]nodetype.NodeType{}, baseTopic}
}
func (d *Device) Set(prop string, value string) {
	switch prop {
	case "$online":
		d.Online, _ = strconv.ParseBool(value)
	case "$name":
		d.Name = value
	case "$localip":
		d.Localip = value
	case "$mac":
		d.Mac = value
	case "$stats/uptime":
		d.Stats.Uptime, _ = strconv.Atoi(value)
	case "$stats/signal":
		d.Stats.Signal, _ = strconv.Atoi(value)
	case "$stats/interval":
		d.Stats.Interval, _ = strconv.Atoi(value)
	case "$fw/name":
		d.Fw.Name = value
	case "$fw/version":
		d.Fw.Version = value
	case "$fw/checksum":
		d.Fw.Checksum = value
	case "$implementation":
		d.Implementation = value
	}
}

func (d *Device) MQTTNodeHandler(mqttClient homieMessage.SubscriptibleClient, mqttMessage homieMessage.HomieExtractableMessage) {
	// Will be bound to devices/<id/>+/$type
	message, err := homieMessage.New(mqttMessage, d.BaseTopic)
	if err != nil {
		return
	}
	topicComponents := strings.Split(message.Path, "/")
	if len(topicComponents) != 2 {
		return
	}
	node, property := topicComponents[0], topicComponents[1]
	if property == "$type" {
		newNode, err := nodetype.New(message.Payload, d.BaseTopic)
		if err == nil {
			d.Nodes[node] = newNode
			mqttClient.Subscribe(d.BaseTopic+d.Id+"/"+node+"/+", 1, d.Nodes[node].MQTTHandler)
		} else {
			fmt.Println("adding node failed: ", err)
		}
	}
}
