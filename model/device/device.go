package device

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/node"
	"strconv"
	"strings"
	"github.com/jbonachera/homie-controller/model/metric"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jbonachera/homie-controller/influxdb"


)

type DeviceStats struct {
	Uptime   int `json:"uptime"`
	Signal   int `json:"signal"`
	Interval int `json:"interval"`
}

type DeviceFirmware struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
}

type Device struct {
	Id             string               `json:"id"`
	Online         bool                 `json:"online"`
	Name           string               `json:"name"`
	Localip        string               `json:"localip"`
	Mac            string               `json:"mac"`
	Stats          DeviceStats          `json:"stats"`
	Fw             DeviceFirmware       `json:"fw"`
	Implementation string               `json:"implementation"`
	Nodes          map[string]node.Type `json:"nodes"`
	BaseTopic      string               `json:"base_topic"`
}

func New(id string, baseTopic string) *Device {
	return &Device{id, false, "", "", "", DeviceStats{0, 0, 0}, DeviceFirmware{"", "", ""}, "", map[string]node.Type{}, baseTopic}
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

func (d *Device) MQTTNodeHandler(mqttClient MQTT.Client, mqttMessage MQTT.Message) {
	// Will be bound to devices/<id/>+/$type
	message, err := homieMessage.New(mqttMessage, d.BaseTopic)
	if err != nil {
		return
	}
	topicComponents := strings.Split(message.Path, "/")
	if len(topicComponents) == 1 {
		log.Debug("updating "+ topicComponents[0]+" for "+d.Id)
		d.Set(topicComponents[0], message.Payload)
		log.Debug(topicComponents[0]+" for "+d.Id + "set to "+message.Payload)
	} else if len(topicComponents) == 2 {
		nodeName, property := topicComponents[0], topicComponents[1]
		switch nodeName {
		case "$stats":
			d.Set(message.Path, message.Payload)
			if influxdb.Ready() {
				influxdb.PublishPoint(d.GetPoint())
			}
		case "$fw":
			d.Set(message.Path, message.Payload)
		default:
			if property == "$type" {
				newNode, err := node.New(message.Payload, nodeName, d.BaseTopic)
				if err == nil {
					d.Nodes[nodeName] = newNode
					properties := newNode.GetProperties()
					log.Debug("adding node " + nodeName + " for device " + message.Id)
					for _, property := range properties {
						log.Debug("adding property " + property + " to node " + nodeName + " for device " + message.Id)
						mqttClient.Subscribe(d.BaseTopic+d.Id+"/"+nodeName+"/"+property, 1, d.Nodes[nodeName].MQTTHandler)
					}
				} else {
					log.Warn("adding node failed: " + err.Error())
				}
			}
		}
	} else {
		return
	}

}

func (d *Device) GetPoint() *client.Point {
	return metric.New("system", map[string]string{"id": d.Id, "name": d.Name}, map[string]interface{}{"signal": d.Stats.Signal, "uptime": d.Stats.Uptime, "sensors_count": len(d.Nodes)})
}
