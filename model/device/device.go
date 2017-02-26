package device

import (
	"github.com/jbonachera/homie-controller/influxdb"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/metric"
	"github.com/jbonachera/homie-controller/model/node"
	"strconv"
	"strings"

	"github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/model/implementation"
	"github.com/jbonachera/homie-controller/model/search"
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
	Id             string                        `json:"id"`
	Online         bool                          `json:"online"`
	Name           string                        `json:"name"`
	Localip        string                        `json:"localip"`
	Mac            string                        `json:"mac"`
	Stats          DeviceStats                   `json:"stats"`
	Fw             DeviceFirmware                `json:"fw"`
	Implementation implementation.Implementation `json:"implementation"`
	Nodes          map[string]node.Type          `json:"nodes"`
	BaseTopic      string                        `json:"base_topic"`
	registrator    HandlerRegistrator
}

type HandlerRegistrator func(topic string, callback messaging.CallbackHandler)

func New(id string, baseTopic string) *Device {
	return &Device{id, false, "", "", "",
		DeviceStats{0, 0, 0}, DeviceFirmware{"", "", ""},
		nil, map[string]node.Type{}, baseTopic, messaging.AddHandler}
}
func (d *Device) SetRegistrator(handler HandlerRegistrator) {
	d.registrator = handler
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
	}
}

func (d *Device) MQTTNodeHandler(message homieMessage.HomieMessage) {
	// TODO: split this into a map of handler functions, and call handler[topic](payload), then default to Set()
	// Will be bound to devices/<id/>+/$type
	topicComponents := strings.Split(message.Path, "/")
	if len(topicComponents) == 1 {
		if topicComponents[0] == "$implementation" {
			var err error
			d.Implementation, err = implementation.New(message.Payload, d.Id, d.BaseTopic)
			if err != nil {
				log.Error(err.Error())
			} else {
				d.registrator(d.BaseTopic+d.Id+"/$implementation/+", func(message homieMessage.HomieMessage) {
					log.Debug("(device: " + d.Name + ") (implementation: " + d.Implementation.GetName() + ") " + message.Path + " → " + message.Payload)
					d.Implementation.MessageHandler(message)
				})
			}

		} else {
			log.Debug("updating " + topicComponents[0] + " for " + d.Id)
			d.Set(topicComponents[0], message.Payload)
			log.Debug(topicComponents[0] + " for " + d.Id + " set to " + message.Payload)
		}
	} else if len(topicComponents) == 2 {
		nodeName, property := topicComponents[0], topicComponents[1]
		switch nodeName {
		case "$stats":
			d.Set(message.Path, message.Payload)
			if influxdb.Ready() {
				influxdb.PublishMetric(d.GetPoint())
			}
		case "$fw":
			d.Set(message.Path, message.Payload)
		default:
			if property == "$type" {
				newNode, err := node.New(message.Payload, nodeName, d.Id, d.BaseTopic)
				if err == nil {
					d.Nodes[nodeName] = newNode
					properties := newNode.GetProperties()
					log.Debug("adding node " + nodeName + " for device " + message.Id)
					for _, property := range properties {
						log.Debug("adding property " + property + " to node " + nodeName + " for device " + message.Id)
						d.registrator(d.BaseTopic+d.Id+"/"+nodeName+"/"+property, func(message homieMessage.HomieMessage) {
							log.Debug("(device: " + d.Name + ") (node: " + d.Nodes[nodeName].GetName() + ") " + message.Path + " → " + message.Payload)
							d.Nodes[nodeName].MessageHandler(message)
						})
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

func (d *Device) GetPoint() *metric.Metric {
	return metric.New("system", map[string]string{"id": d.Id, "name": d.Name}, map[string]interface{}{"signal": d.Stats.Signal, "uptime": d.Stats.Uptime, "sensors_count": len(d.Nodes)})
}

func (d *Device) Match(opts search.Options) bool {
	match := false
	if term, ok := opts.Terms["online"]; ok {
		if (term == "true") == d.Online {
			match = true
		} else {
			return false
		}
	}
	if term, ok := opts.Terms["type"]; ok {
		type_matched := false
		for name := range d.Nodes {
			if name == term {
				type_matched = true
				break
			}
		}
		if !type_matched {
			return false
		} else {
			match = type_matched
		}
	}

	return match
}
