package ota

import "testing"

type MockFirmware struct {
}

func (f *MockFirmware) Name() string {
	return "mockFirmware"
}

func (f *MockFirmware) Version() string {
	return "1.0.1"
}

func (f *MockFirmware) Checksum() string {
	return ""
}

func (f *MockFirmware) Payload() []byte {
	return []byte{}
}

type MockProvider struct {
	id string
}

func (m *MockProvider) GetLatest() Firmware {
	return &MockFirmware{}
}
func (m *MockProvider) Id() string {
	return m.id
}

type MockFactory struct {
	id string
}

func (m *MockFactory) Id() string {
	return m.id
}

func (m *MockFactory) New(name string) FirmwareProvider {
	return &MockProvider{id: name}
}

func TestRegisterFactory(t *testing.T) {
	factory := &MockFactory{id: "mockFactory"}
	RegisterFactory("mock", factory)
	if factories["mock"].Id() != "mockFactory" {
		t.Error("could not register a new factory")
	}
}

func TestAddFirmware(t *testing.T) {
	factories["mock"] = &MockFactory{id: "mockFactory"}
	AddFirmware("mockFirmware", "mock")
	if _, ok := firmwares["mockFirmware"]; !ok {
		t.Error("firmware was not registered")
	} else if firmwares["mockFirmware"].Id() != "mockFirmware" {
		t.Error("firmware was not correctly registered")
	}
}
