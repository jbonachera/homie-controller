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

var registry Registry


func NewRegistry(baseTopic string) {
	registry = Registry{sync.Mutex{}, []Device{}, map[string]int{}, baseTopic}
}

func addToIndex(index map[string]int, key string, offset int) {
	index[key] = offset
}

func Append(device Device) {
	registry.Lock()
	defer registry.Unlock()
	registry.devices = append(registry.devices, device)
	addToIndex(registry.devicesIndex, device.Id, len(registry.devices)-1)
}

func Get(id string) Device {
	registry.Lock()
	defer registry.Unlock()
	offset, ok := registry.devicesIndex[id]
	if ok {
		return registry.devices[offset]
	}
	return Device{}

}

func Set(id string, path string, value string) {
	registry.Lock()
	defer registry.Unlock()
	offset, ok := registry.devicesIndex[id]
	if ok {
		registry.devices[offset].Set(path, value)
	}
}
func OnlineCallback(client MQTT.Client, mqttMessage MQTT.Message) {
	message, err := homieMessage.New(mqttMessage, registry.baseTopic)
	if err != nil {
		log.Warn("received an invalid message")
		return
	}
	if message.Payload == "true" {
		log.Debug("discovered a new newDevice: " + message.Id)
		newDevice := New(message.Id, registry.baseTopic)
		Append(newDevice)
		for _, prop := range homieMessage.Properties {
			client.Subscribe(registry.baseTopic+message.Id+"/"+prop, 1, newDevice.MQTTNodeHandler)
		}
		client.Subscribe(registry.baseTopic+message.Id+"/+/$type", 1, newDevice.MQTTNodeHandler)
	} else {
		log.Debug("a device has disconnected: " + message.Id)
		for _, prop := range homieMessage.Properties {
			client.Unsubscribe(registry.baseTopic + message.Id + "/" + prop)
		}
		client.Unsubscribe(registry.baseTopic + message.Id + "/+/$type")
	}
}
