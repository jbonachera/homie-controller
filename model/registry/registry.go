package registry

import (
	"github.com/jbonachera/homie-controller/model/device"
	"sync"
)

type Registry struct {
	sync.Mutex
	devices []device.Device
}

func (d *Registry) Append(device device.Device) {
	d.Lock()
	d.devices = append(d.devices, device)
	d.Unlock()
}

func (d Registry) Get(id string) device.Device {
	d.Lock()
	for _, device := range d.devices {
		if device.Id == id {
			return device
		}
	}
	d.Unlock()
	return device.Device{}

}

func (d *Registry) Set(id string, path string, value string) {
	d.Lock()
	for idx, device := range d.devices {
		if device.Id == id {
			d.devices[idx].Set(path, value)
		}
	}
	d.Unlock()
}
