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
			if path == "$properties" {
				d.addNode(nodeName, value)
			}
		} else {
			if path == "$properties" {
				delete(d.Nodes, nodeName)
				d.addNode(nodeName, value)
			}
			d.Nodes[nodeName].Properties[path] = value
		}

	}
}

func (d *Device) addNode(node string, properties string) {
	d.Nodes[node] = DeviceNode{node, "", map[string]string{}}
}
