package ota

import (
	"errors"
	"github.com/jbonachera/homie-controller/log"
	"github.com/mcuadros/go-version"
)

var firmwares map[string]map[string]FirmwareProvider
var factories map[string]map[string]FirmwareFactory

type Firmware interface {
	Name() string
	Brand() string
	Version() string
	Checksum() string
	Payload() []byte
	Announce()
}

type FirmwareProvider interface {
	Id() string
	Brand() string
	GetLatest() Firmware
}

type FirmwareFactory interface {
	New(name string, brand string) FirmwareProvider
	Id() string
}

func init() {
	factories = map[string]map[string]FirmwareFactory{}
	firmwares = map[string]map[string]FirmwareProvider{}
}

func RegisterFactory(name string, brand string, provider FirmwareFactory) {
	if _, present := factories[brand]; !present {
		factories[brand] = map[string]FirmwareFactory{}
	}
	factories[brand][name] = provider
}

func AddFirmware(name string, brand string, provider string) {
	if _, present := firmwares[brand]; !present {
		firmwares[brand] = map[string]FirmwareProvider{}
	}
	if firmware, present := factories[brand][provider]; present {
		firmwares[brand][name] = firmware.New(name, brand)
	}
}

func IsUpToDate(firmware string, brand string, current_version string) (bool, error) {
	if _, ok := firmwares[brand]; !ok {
		return false, errors.New("brand not found in OTA")
	}
	if firmwareProvider, ok := firmwares[brand][firmware]; ok {
		return IsVersionGreater(current_version, firmwareProvider.GetLatest().Version()), nil

	} else {
		return false, errors.New("firmware not found in OTA")
	}
}

func IsVersionGreater(local string, remote string) bool {
	return version.Compare(local, remote, ">=")
}

func LastVersion(firmware string, brand string) string {
	if _, ok := firmwares[brand]; !ok {
		return "unknown"
	}
	if firmwareProvider, ok := firmwares[brand][firmware]; ok {
		return firmwareProvider.GetLatest().Version()

	} else {
		return "unknown"
	}
}

func LastFirmware(firmware string, brand string) (Firmware, error) {
	if _, ok := firmwares[brand]; !ok {
		return nil, errors.New("brand not found in OTA")
	}
	if firmwareProvider, ok := firmwares[brand][firmware]; ok {
		return firmwareProvider.GetLatest(), nil

	} else {
		return nil, errors.New("firmware not found in OTA")
	}
}

func Refresh() {
	for brand := range firmwares {
		for _, provider := range firmwares[brand] {
			log.Info("fetching last version of firmware " + provider.Id())
			go provider.GetLatest()
		}
	}
}
