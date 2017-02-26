package ota

import (
	"errors"
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

func IsUpToDate(firmware string, current_version string) (bool, error) {
	firmwareProvider, ok := firmwares[firmware]
	if !ok {
		return false, errors.New("firmware not found in OTA")
	}
	return version.Compare(current_version, firmwareProvider.GetLatest().Version(), ">="), nil
}

func LastVersion(firmware string) string {
	firmwareProvider, ok := firmwares[firmware]
	if !ok {
		return "unknown"
	}
	return firmwareProvider.GetLatest().Version()
}
