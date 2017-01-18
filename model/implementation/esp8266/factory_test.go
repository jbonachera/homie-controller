package esp8266

import "testing"

func TestFactory_New(t *testing.T) {
	factory := new(Factory)
	esp := factory.New("devices/")
	if esp.GetName() != "esp8266"{
		t.Error("did create an invalid 8266 implementation")
	}
}

func BenchmarkFactory_New(b *testing.B) {
	factory := new(Factory)
	for i := 0; i < b.N; i++ {
		_ = factory.New("devices/")
	}
}