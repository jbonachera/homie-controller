package device

import (
	MQTT "github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/mocks/mqtt"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/model/search"
	"strconv"
	"sync"
	"testing"
)

var baseTopic string = "devices/"

func TestAppend(t *testing.T) {
	NewRegistry(baseTopic)
	Append(New("u1", baseTopic))
	if registry.devices["u1"].Id != "u1" {
		t.Error("didn't get the device we just inserted")
	}
}

func appendList(count int, start int, wg *sync.WaitGroup) {
	defer wg.Done()
	i := start
	for i < count+start {
		Append(New("u"+strconv.Itoa(i), baseTopic))
		i += 1
	}
}
func TestParralellAppend(t *testing.T) {
	var wg sync.WaitGroup
	count := 100
	NewRegistry(baseTopic)
	wg.Add(count)
	i := 0
	for i < count {
		go appendList(200, i*200, &wg)
		i += 1
	}
	wg.Wait()
	if len(registry.devices) != 200*count {
		t.Error("some elements were not inserted: wanted ", 200*count, " , got ", len(registry.devices))
	}
}

func populate(count int) *Registry {
	NewRegistry(baseTopic)
	i := 0
	for i < count {
		Append(New("u"+strconv.Itoa(i), "devices/"))
		i += 1
	}
	return &registry
}
func TestGet(t *testing.T) {
	populate(30)
	device, err := Get("u17")
	if err != nil {
		t.Error("got error when calling Get():", err)
	}
	if device.Id != "u17" {
		t.Error("could not Get() a created device: got ", device.Id, " , wanted u17")
	}

}
func TestRegistrySet(t *testing.T) {
	populate(30)
	Set("u17", "$online", "true")
	if registry.devices["u17"].Online != true {
		t.Error("could not Set() a property: wanted true, got", registry.devices["u17"].Online)
	}
}

func setProp(count int, prop string, id string, value string, wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0
	for i < count {
		Set(id, prop, value)
		i += 1
	}
}
func TestParallelSet(t *testing.T) {
	populate(30)
	var wg sync.WaitGroup
	i := 0
	wg.Add(30 * 3)
	for i < 30 {
		go setProp(20, "$online", "u17", "false", &wg)
		go setProp(20, "$online", "u17", "true", &wg)
		go setProp(20, "$stats/signal", "u17", "87", &wg)
		i += 1
	}
	wg.Wait()
	myDevice, err := Get("u17")
	if err != nil {
		t.Error("got error when calling Get():", err)
	}
	if myDevice.Stats.Signal != 87 {
		t.Error("could not set property: wanted 87, got", myDevice.Stats.Signal)
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
	populate(j)
	for i := 0; i < b.N; i++ {
		Set("u988", "$online", "true")
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
	populate(j)
	for i := 0; i < b.N; i++ {
		Get("u988")
	}
}

func TestList(t *testing.T) {
	populate(30)
	list := List()
	if len(list) != 30 {
		t.Error("could not get a list of devices")
	}
	found := false
	for _, item := range list {
		if item == "u0" {
			found = true
		}
	}
	if !found {
		t.Error("item 0 sould be u0, got", list[0])
	}
}
func BenchmarkList(b *testing.B) {
	benchmarkList(5, b)
}
func BenchmarkList10(b *testing.B) {
	benchmarkList(10, b)
}
func BenchmarkList100(b *testing.B) {
	benchmarkList(100, b)
}
func BenchmarkList1000(b *testing.B) {
	benchmarkList(1000, b)
}
func benchmarkList(j int, b *testing.B) {
	populate(j)
	for i := 0; i < b.N; i++ {
		List()
	}
}

func TestOnlineCallback(t *testing.T) {
	manager := RegistrationHandler{
		add: func(topic string, callback MQTT.CallbackHandler) {},
		del: func(topic string) {},
	}
	NewRegistry(baseTopic)
	SetRegistrationManager(manager)
	var message homieMessage.HomieMessage
	message, _ = homieMessage.Extract(mqtt.NewMessage(
		"devices/u0/$online",
		"true",
	), "devices/")
	OnlineCallback(message)
	device, err := Get("u0")
	if err != nil {
		t.Error("test device was not created")
	}
	message, err = homieMessage.Extract(mqtt.NewMessage(
		"devices/u0/$name",
		"sensor-test",
	), "devices/")
	device.MQTTNodeHandler(message)
	if registry.devices["u0"].Name != "sensor-test" {
		t.Error("test device was not updated")
	}
}

func TestGetAll(t *testing.T) {
	populate(30)
	deviceList := GetAll()
	if len(deviceList) != 30 {
		t.Error("GetAll did not return all the devices")
	}
}

func TestSearch(t *testing.T) {
	populate(30)
	registry.devices["u1"].Online = true
	registry.devices["u3"].Online = true
	searchTerms := map[string]string{
		"online": "true",
	}
	matched := Search(search.Options{Terms: searchTerms})
	if len(matched) != 2 {
		t.Error("device search failed")
	}
}

func BenchmarkSearch1000(b *testing.B) {
	benchmarkSearch(1000, b)
}
func BenchmarkSearch10000(b *testing.B) {
	benchmarkSearch(10000, b)
}
func benchmarkSearch(j int, b *testing.B) {
	populate(j)
	registry.devices["u1"].Online = true
	registry.devices["u3"].Online = true
	searchTerms := map[string]string{
		"online": "true",
	}
	for i := 0; i < b.N; i++ {
		Search(search.Options{Terms: searchTerms})
	}
}
