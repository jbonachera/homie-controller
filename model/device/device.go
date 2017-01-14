package device

import (
	"strconv"
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
	Type       string
	Properties [][]string
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
	Nodes          []DeviceNode
}

func New(id string) Device {
	return Device{id, false, "", "", "", DeviceStats{0, 0, 0}, DeviceFirmware{"", "", ""}, "", []DeviceNode{}}
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
