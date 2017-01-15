package registry

import (
	"github.com/jbonachera/homie-controller/model/device"
	"strconv"
	"sync"
	"testing"
)

func TestAppend(t *testing.T) {
	registry := Registry{}
	registry.Append(device.New("u1"))
	if registry.devices[0].Id != "u1" {
		t.Error("didn't get the device we just inserted")
	}
	registry.Append(device.New("u2"))
	if registry.devices[0].Id != "u1" {
		t.Error("existing device disapeared after insertion")
	}
	if registry.devices[1].Id != "u2" {
		t.Error("second device were not inserted")
	}
}

func appendList(count int, r *Registry, wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0
	for i < count {
		r.Append(device.New("u" + strconv.Itoa(i)))
		i += 1
	}
}
func TestParralellAppend(t *testing.T) {
	var wg sync.WaitGroup
	count := 100
	registry := Registry{}
	wg.Add(count)
	i := 0
	for i < count {
		go appendList(200, &registry, &wg)
		i += 1
	}
	wg.Wait()
	if len(registry.devices) != 200*count {
		t.Error("some elements were not inserted: wanted ", 200*count, " , got ", len(registry.devices))
	}
}

func populate(count int) *Registry {
	registry := Registry{}
	i := 0
	for i < 30 {
		registry.Append(device.New("u" + strconv.Itoa(i)))
		i += 1
	}
	return &registry
}
func TestGet(t *testing.T) {
	registry := populate(30)
	device := registry.Get("u17")
	if device.Id != "u17" {
		t.Error("could not Get() a created device: got ", device.Id, " , wanted u17")
	}

}
func TestSet(t *testing.T) {
	registry := populate(30)
	registry.Set("u17", "$online", "true")
	if registry.devices[17].Online != true {
		t.Error("could not Set() a property: wanted true, got", registry.devices[17].Online)
	}
}

func setProp(count int, prop string, id string, value string, r *Registry, wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0
	for i < count {
		r.Set(id, prop, value)
		i += 1
	}
}
func TestParrallelSet(t *testing.T) {
	r := populate(30)
	var wg sync.WaitGroup
	i := 0
	wg.Add(30 * 3)
	for i < 30 {
		go setProp(20, "$online", "u17", "false", r, &wg)
		go setProp(20, "$online", "u17", "true", r, &wg)
		go setProp(20, "$stats/signal", "u17", "87", r, &wg)
		i += 1
	}
	wg.Wait()
	device := r.Get("u17")
	if device.Stats.Signal != 87 {
		t.Error("could not set property: wanted 87, got", device.Stats.Signal)
	}
}
