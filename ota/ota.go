package ota

import (
	"errors"
	"github.com/jbonachera/homie-controller/config"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/ota/persistentCache"
	"github.com/mcuadros/go-version"
	"strconv"
	"time"
)

var firmwares map[string]FirmwareProvider
var factories map[string]FirmwareFactory

var done chan bool
var bootTime time.Time

type Firmware interface {
	Name() string
	Version() string
	Checksum() string
	Payload() []byte
	Announce()
}

type FirmwareProvider interface {
	Id() string
	GetLatest() Firmware
	GetLastFive()
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
	if firmware, present := factories[provider]; present {
		firmwares[name] = firmware.New(name)
		firmwares[name].GetLatest()
		firmwares[name].GetLastFive()
	}
}

func IsUpToDate(firmware string, current_version string) (bool, error) {
	if firmwareProvider, ok := firmwares[firmware]; ok {
		return IsVersionGreater(current_version, firmwareProvider.GetLatest().Version()), nil

	} else {
		return false, errors.New("firmware not found in OTA")
	}
}

func IsVersionGreater(local string, remote string) bool {
	return version.Compare(local, remote, ">=")
}

func LastVersion(firmware string) string {
	if firmwareProvider, ok := firmwares[firmware]; ok {
		return firmwareProvider.GetLatest().Version()

	} else {
		return "unknown"
	}
}

func LastFirmware(firmware string) (Firmware, error) {
	if firmwareProvider, ok := firmwares[firmware]; ok {
		return firmwareProvider.GetLatest(), nil

	} else {
		return nil, errors.New("firmware not found in OTA")
	}
}

func Refresh() {
	for _, provider := range firmwares {
		log.Info("fetching last version of firmware " + provider.Id())
		go provider.GetLatest()
		go provider.GetLastFive()
	}
}

func Stop() {
	done <- true
}

func PublishStats() {
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/$stats/uptime", Payload: strconv.Itoa(int(time.Since(bootTime).Seconds()))})
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/$stats/signal", Payload: "100"})
}

func Start() {
	done = make(chan bool, 1)
	persistentCache.Start()
	bootTime = time.Now()
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/$online", Payload: "true"})
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/$name", Payload: "controller"})
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/$homie", Payload: "2.0.0"})
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/$stats/interval", Payload: "10"})
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/$fw/name", Payload: "homie-controller"})
	messaging.PublishState(homieMessage.HomieMessage{Topic: "devices/controller/$fw/version", Payload: config.Get("version")})

	go func() {
		run := true
		for run {
			select {
			case <-done:
				run = false
				break
			case <-time.After(10 * time.Second):
				go PublishStats()
			case <-time.After(24 * time.Hour):
				go Refresh()
				break
			}
		}
		log.Debug("Finished OTA routing")
	}()

}
