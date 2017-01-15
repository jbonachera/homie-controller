package device

import (
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

type DeviceNode struct {
	Name       string
	Type       string
	Properties map[string]string
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
	Nodes          map[string]DeviceNode
}

func New(id string) Device {
	return Device{id, false, "", "", "", DeviceStats{0, 0, 0}, DeviceFirmware{"", "", ""}, "", map[string]DeviceNode{}}
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
	default:
		// We suppose we are dealing with a Node
		splitted_path := strings.Split(prop, "/")
		if len(splitted_path) != 2 {
			// invalid request
			return
		}
		nodeName, path := splitted_path[0], splitted_path[1]
		_, exists := d.Nodes[nodeName]
		if !exists {
			d.addNode(nodeName)
		}

		node := d.Nodes[nodeName]
		switch path {
		case "$properties":
			node.setProperties(strings.Split(value, ","))
		case "$type":
			node.Type = value
			d.Nodes[nodeName] = node
		default:
			d.Nodes[nodeName].Properties[path] = value
		}
	}
}

func (d *Device) deleteNode(name string) {
	delete(d.Nodes, name)
}

func (d *Device) addNode(nodeName string) {
	d.Nodes[nodeName] = DeviceNode{nodeName, "", map[string]string{}}
}

func (n *DeviceNode) setProperties(properties []string) {
	// Handles deletion of already-present properties we don't want
	for property, _ := range n.Properties {
		for _, wantedProperty := range properties {
			if wantedProperty == property {
				break
			}
		}
		delete(n.Properties, property)
	}
	// Handles creation of new properties
	for _, wantedProperty := range properties {
		_, ok := n.Properties[wantedProperty]
		if !ok {
			n.Properties[wantedProperty] = ""
		}
	}
}

func (d *Device) ListTypes() []string {
	types := []string{}
	dedup := make(map[string]struct{})
	for _, node := range d.Nodes {
		_, present := dedup[node.Type]
		if !present {
			types = append(types, node.Type)
		}
	}
	return types
}
