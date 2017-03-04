package esp8266

import (
	"encoding/json"
	"errors"
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/messaging"
	"github.com/jbonachera/homie-controller/model/device"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"github.com/jbonachera/homie-controller/ota"
	"strconv"
	"strings"
	"time"
)

const (
	OTANotRunning        = iota
	OTARequested         = iota
	OTAReceivingFirmware = iota
)

type esp8266 struct {
	parentId         string
	Name             string
	Version          string   `json:"version"`
	Ota              bool     `json:"ota"`
	Actions          []string `json:"actions"`
	Config           config   `json:"config"`
	baseTopic        string
	MessagePublisher messagePublisherHandler `json:"-"`
	ActionHandlers   map[string]func()       `json:"-"`
	OTARunning       bool                    `json:"-"`
	OTAStep          int                     `json:"-"`
	OTALastUpdate    time.Time               `json:"-"`
}

type messagePublisherHandler func(message homieMessage.HomieMessage)

func New(parent string, baseTopic string) *esp8266 {
	esp := &esp8266{parent, "", "", false, []string{}, config{}, baseTopic, messaging.PublishMessage, map[string]func(){}, false, OTANotRunning, time.Time{}}
	actionHandlers := map[string]func(){
		"reset":   esp.Reset,
		"upgrade": esp.StartOTA,
	}
	actions := make([]string, len(actionHandlers))

	i := 0
	for k := range actionHandlers {
		actions[i] = k
		i++
	}
	esp.Name = esp.GetName()
	esp.ActionHandlers = actionHandlers
	esp.Actions = actions
	return esp
}

func (e *esp8266) Do(action string) error {
	parentDevice, err := device.Get(e.parentId)
	if err != nil {
		return errors.New("device not found. aborting command execution.")
	}
	if !parentDevice.Online {
		return errors.New("device is offline")
	}
	if handler := e.ActionHandlers[action]; handler != nil {
		handler()
		return nil
	} else {
		return errors.New("method not found: " + action)
	}
}

func (e *esp8266) GetName() string {
	return "esp8266"
}

func (e *esp8266) Reset() {
	message, err := homieMessage.New(e.parentId, e.baseTopic, "$implementation/reset", "true")
	if err != nil {
		log.Error("failed to create message for reset")
		return
	}
	e.MessagePublisher(message)
}

func (e *esp8266) Set(property string, value string) {
	switch property {
	case "version":
		e.Version = value
	case "config":
		if err := json.Unmarshal([]byte(value), &e.Config); err != nil {
			log.Error("error while parsing device config: " + err.Error())
		}
	case "ota/enabled":
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			e.Ota = boolValue
		}
	}
}
func (e *esp8266) GetProperties() []string {
	return []string{"Version", "Ota"}
}

func (e *esp8266) MessageHandler(message homieMessage.HomieMessage) {
	// Will be bound to devices/<id>/implementation/#

	topicComponents := strings.SplitN(message.Path, "/", 2)
	if len(topicComponents) < 2 {
		return
	}
	propertyName := topicComponents[1]
	e.Set(propertyName, message.Payload)
}

func (e *esp8266) checkOTATimeout() {
	for e.OTARunning {
		time.Sleep(5 * time.Second)
		log.Debug("checking if OTA timed out...")
		if time.Since(e.OTALastUpdate) > 5*time.Second {
			log.Error("OTA timeout")
			e.AbortOTA()
			break
		}
	}
}

func (e *esp8266) AbortOTA() {
	log.Error("aborting OTA")
	messaging.DelSubscription(e.baseTopic + e.parentId + "/$implementation/ota/status")
	e.OTARunning = false
	e.OTAStep = OTANotRunning
}

func (e *esp8266) StartOTA() {
	if e.OTARunning {
		log.Error("OTA is already running")
		return
	}
	e.OTARunning = true
	e.OTALastUpdate = time.Now()
	go e.checkOTATimeout()
	parentDevice, err := device.Get(e.parentId)
	if err != nil {
		log.Error("parent device not found.")
		e.AbortOTA()
		return
	}
	firmware, err := ota.LastFirmware(parentDevice.Fw.Name)
	if err != nil {
		log.Error("error while fetching firmware: " + err.Error())
		e.AbortOTA()
		return
	}
	if ota.IsVersionGreater(parentDevice.Fw.Version, firmware.Version()) {
		log.Error("local version is more up-to-date than remote version.")
		e.AbortOTA()
		return
	}

	messaging.AddHandler(e.baseTopic+e.parentId+"/$implementation/ota/status", func(message homieMessage.HomieMessage) {
		parts := strings.SplitN(message.Payload, " ", 2)
		statusCode, payload := parts[0], ""
		if len(parts) > 1 {
			payload = parts[1]
		}
		e.OTALastUpdate = time.Now()
		switch statusCode {
		case "200":
			{
				if e.OTAStep >= OTAReceivingFirmware {
					log.Info("OTA successfull")
					messaging.DelSubscription(e.baseTopic + e.parentId + "/$implementation/ota/status")
					e.OTARunning = false
					e.OTAStep = OTANotRunning
				}
			}
		case "403":
			{
				log.Error("OTA is disabled on device")
				e.AbortOTA()
			}
		case "400":
			{
				log.Error("device " + e.parentId + "said the checksum was malformed.")
				e.AbortOTA()
			}
		case "202":
			{
				if e.OTAStep >= OTARequested {
					log.Info("device " + e.parentId + " accepted to start OTA")
					firmwareTopic := e.baseTopic + e.parentId + "/$implementation/ota/firmware"
					log.Debug("publishing firmware to " + firmwareTopic)
					e.OTAStep = OTAReceivingFirmware
					messaging.PublishFile(firmwareTopic, firmware.Payload())
				}

			}
		case "304":
			{
				log.Debug("device is already up to date")
				e.AbortOTA()
			}
		case "206":
			{
				log.Debug("device is receveing firmware: " + payload)
			}
		}
	})
	message, err := homieMessage.New(e.parentId, e.baseTopic, "$implementation/ota/checksum", firmware.Checksum())
	if err != nil {
		log.Error("error while creating message for device: " + err.Error())
		e.AbortOTA()
	} else {
		e.OTAStep = OTARequested
		messaging.PublishMessage(message)
	}
}
