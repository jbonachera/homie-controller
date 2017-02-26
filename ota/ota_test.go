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

func TestIsUpToDate(t *testing.T) {
	firmwares["mock"] = &MockProvider{id: "mock"}
	if uptodate, _ := IsUpToDate("mock", "1.0.1"); !uptodate {
		t.Error("should have detected current version were the latest")
	}
	if uptodate, _ := IsUpToDate("mock", "1.0.0"); uptodate {
		t.Error("should have detected current version were not the latest: got 1.0.1 > 1.0.0")
	}
	if uptodate, _ := IsUpToDate("mock", "1.0.2"); !uptodate {
		t.Error("should have detected version were greater than the latest: got 1.0.1 > 1.0.2")
	}
}

func TestLastVersion(t *testing.T) {
	firmwares["mock"] = &MockProvider{id: "mock"}
	if LastVersion("mock") != "1.0.1" {
		t.Error("could not get last version")
	}
}

func TestLastFirmware(t *testing.T) {
	firmwares["mock"] = &MockProvider{id: "mock"}
	if firmware, _ := LastFirmware("mock"); firmware.Version() != "1.0.1" {
		t.Error("could not get last version")
	}
}
