package device

import (
	"errors"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"sync"
)

type RegistrationHandler struct {
	add func(topic string, callback messaging.CallbackHandler)
	del func(topic string)
}

type Registry struct {
	sync.Mutex
	devices             map[string]*Device
	baseTopic           string
	registrationManager RegistrationHandler
}

var registry Registry

func List() []string {
	keys := make([]string, 0, len(registry.devices))
	for k := range registry.devices {
		keys = append(keys, k)
	}

	return keys
}

func NewRegistry(baseTopic string) {
	registry = Registry{sync.Mutex{}, map[string]*Device{}, baseTopic, RegistrationHandler{
		add: messaging.AddHandler,
		del: messaging.DelSubscription,
	}}
}

func SetRegistrationManager(manager RegistrationHandler) {
	registry.registrationManager = manager
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
func OnlineCallback(message homieMessage.HomieMessage) {
	device, err := Get(message.Id)
	if message.Payload == "true" {
		if err != nil {
			log.Debug("discovered a new device: " + message.Id)
			newDevice := New(message.Id, registry.baseTopic)
			Append(newDevice)
			for _, prop := range homieMessage.Properties {
				registry.registrationManager.add(registry.baseTopic+message.Id+"/"+prop, newDevice.MQTTNodeHandler)
			}
			registry.registrationManager.add(registry.baseTopic+message.Id+"/+/$type", newDevice.MQTTNodeHandler)
			newDevice.Set("$online", "true")
		} else {
			log.Debug("device " + message.Id+ " came back to life")
			for _, prop := range homieMessage.Properties {
				registry.registrationManager.add(registry.baseTopic+message.Id+"/"+prop, device.MQTTNodeHandler)
			}
			registry.registrationManager.add(registry.baseTopic+message.Id+"/+/$type", device.MQTTNodeHandler)
			device.Set("$online", "true")
		}
	} else {
		if err == nil {
			log.Debug("a device has disconnected: " + message.Id)
			device.Set("$online", "false")
			for _, prop := range homieMessage.Properties {
				registry.registrationManager.del(registry.baseTopic + message.Id + "/" + prop)
			}
			registry.registrationManager.del(registry.baseTopic + message.Id + "/+/$type")
		}
	}
}

func GetAll() map[string]*Device {
	return registry.devices
}
