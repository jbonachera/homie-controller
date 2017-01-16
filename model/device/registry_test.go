package device

import (
	"strconv"
	"sync"
	"testing"
)

var baseTopic string = "devices/"

func TestAppend(t *testing.T) {
	registry := NewRegistry(baseTopic)
	registry.Append(New("u1", baseTopic))
	if registry.devices[0].Id != "u1" {
		t.Error("didn't get the device we just inserted")
	}
	registry.Append(New("u2", baseTopic))
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
		r.Append(New("u"+strconv.Itoa(i), baseTopic))
		i += 1
	}
}
func TestParralellAppend(t *testing.T) {
	var wg sync.WaitGroup
	count := 100
	registry := NewRegistry(baseTopic)
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
	registry := NewRegistry(baseTopic)
	i := 0
	for i < count {
		registry.Append(New("u"+strconv.Itoa(i), "devices/"))
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
func TestRegistrySet(t *testing.T) {
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

func BenchmarkSet10(b *testing.B) {
	benchmarkSet(10, b)
}

func BenchmarkSet100(b *testing.B) {
	benchmarkSet(100, b)
}
func BenchmarkSet1000(b *testing.B) {
	benchmarkSet(1000, b)
}
func BenchmarkSet10000(b *testing.B) {
	benchmarkSet(10000, b)
}
func benchmarkSet(j int, b *testing.B) {
	r := populate(j)
	for i := 0; i < b.N; i++ {
		r.Set("u988", "$online", "true")
	}
}

func BenchmarkGet10(b *testing.B) {
	benchmarkGet(10, b)
}

func BenchmarkGet100(b *testing.B) {
	benchmarkGet(100, b)
}
func BenchmarkGet1000(b *testing.B) {
	benchmarkGet(1000, b)
}
func BenchmarkGet10000(b *testing.B) {
	benchmarkGet(10000, b)
}

func benchmarkGet(j int, b *testing.B) {
	r := populate(j)
	for i := 0; i < b.N; i++ {
		r.Get("u988")
	}
}
