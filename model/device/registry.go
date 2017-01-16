package device

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"sync"
)

type Registry struct {
	sync.Mutex
	devices      []Device
	devicesIndex map[string]int
	baseTopic    string
}

func NewRegistry(baseTopic string) Registry {
	return Registry{sync.Mutex{}, []Device{}, map[string]int{}, baseTopic}
}

func addToIndex(index map[string]int, key string, offset int) {
	index[key] = offset
}

func (d *Registry) Append(device Device) {
	d.Lock()
	defer d.Unlock()
	d.devices = append(d.devices, device)
	addToIndex(d.devicesIndex, device.Id, len(d.devices)-1)
}

func (d Registry) Get(id string) Device {
	d.Lock()
	defer d.Unlock()
	offset, ok := d.devicesIndex[id]
	if ok {
		return d.devices[offset]
	}
	return Device{}

}

func (d *Registry) Set(id string, path string, value string) {
	d.Lock()
	defer d.Unlock()
	offset, ok := d.devicesIndex[id]
	if ok {
		d.devices[offset].Set(path, value)
	}
}
func (d *Registry) DeviceOnlineCallback(client MQTT.Client, mqttMessage MQTT.Message) {
	message, err := homieMessage.New(mqttMessage, d.baseTopic)
	if err != nil {
		log.Warn("received an invalid message")
		return
	}
	if message.Payload == "true" {
		log.Debug("discovered a new device: " + message.Id)
		device := New(message.Id, d.baseTopic)
		d.Append(device)
		for _, prop := range homieMessage.Properties {
			client.Subscribe(d.baseTopic+message.Id+"/"+prop, 1, device.MQTTNodeHandler)
		}
		client.Subscribe(d.baseTopic+message.Id+"/+/$type", 1, device.MQTTNodeHandler)
	} else {
		log.Debug("a device has disconnected: " + message.Id)
		for _, prop := range homieMessage.Properties {
			client.Unsubscribe(d.baseTopic + message.Id + "/" + prop)
		}
		client.Unsubscribe(d.baseTopic + message.Id + "/+/$type")
	}
}
