package ota

import (
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/mcuadros/go-version"
)

var firmwares map[string]FirmwareProvider
var factories map[string]FirmwareFactory

type Firmware interface {
	Name() string
	Version() string
	Checksum() string
	Payload() []byte
}

type FirmwareProvider interface {
	Id() string
	GetLatest() Firmware
}

type FirmwareFactory interface {
	New(name string) FirmwareProvider
	Id() string
}

func init() {
	factories = map[string]FirmwareFactory{}
	firmwares = map[string]FirmwareProvider{}
}

func RegisterFactory(name string, provider FirmwareFactory) {
	factories[name] = provider
}

func AddFirmware(name string, provider string) {
	firmwares[name] = factories[provider].New(name)
}

func IsUpToDate(firmware string, current_version string) bool {
	return version.Compare(current_version, firmwares[firmware].GetLatest().Version(), ">=")
}
