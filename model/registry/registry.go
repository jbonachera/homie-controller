package registry

import (
	"github.com/jbonachera/homie-controller/model/device"
	"sync"
)

type Registry struct {
	sync.Mutex
	devices      []device.Device
	devicesIndex map[string]int
}

func New() Registry {
	return Registry{sync.Mutex{}, []device.Device{}, map[string]int{}}
}

func (d *Registry) Append(device device.Device) {
	d.Lock()
	d.devices = append(d.devices, device)
	d.devicesIndex[device.Id] = len(d.devices) - 1
	d.Unlock()
}

func (d Registry) Get(id string) device.Device {
	d.Lock()
	defer d.Unlock()
	for _, device := range d.devices {
		if device.Id == id {
			return device
		}
	}
	return device.Device{}

}

func (d *Registry) Set(id string, path string, value string) {
	d.Lock()
	defer d.Unlock()
	offset, ok := d.devicesIndex[id]
	if ok {
		d.devices[offset].Set(path, value)
	}
}
