package device

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"sync"
	"errors"
)

type Registry struct {
	sync.Mutex
	devices      map[string]*Device
	baseTopic    string
}

var registry Registry

func List() []string{
	keys := make([]string, 0, len(registry.devices))
	for k := range registry.devices {
		keys = append(keys, k)
	}

	return keys
}

func NewRegistry(baseTopic string) {
	registry = Registry{sync.Mutex{}, map[string]*Device{}, baseTopic}
}

func Append(device *Device) {
	registry.Lock()
	defer registry.Unlock()
	registry.devices[device.Id] = device
}

func Get(id string) (*Device, error) {
	registry.Lock()
	defer registry.Unlock()
	wantedDevice, ok := registry.devices[id]
	if ok {
		return wantedDevice, nil
	}
	return nil, errors.New("device not found")

}

func Set(id string, path string, value string) {
	registry.Lock()
	defer registry.Unlock()
	wantedDevice, ok := registry.devices[id]
	if ok {
		wantedDevice.Set(path, value)
		registry.devices[id] = wantedDevice
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
