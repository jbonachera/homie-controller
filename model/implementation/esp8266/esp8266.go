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
)

type esp8266 struct {
	parentId         string
	Name             string
	Version          string   `json:"version"`
	LastVersion      string   `json:"last_version"`
	Ota              bool     `json:"ota"`
	UpToDate         bool     `json:"up_to_date"`
	Actions          []string `json:"actions"`
	Config           config   `json:"config"`
	baseTopic        string
	MessagePublisher messagePublisherHandler `json:"-"`
	ActionHandlers   map[string]func()       `json:"-"`
}

type messagePublisherHandler func(message homieMessage.HomieMessage)

func New(parent string, baseTopic string) *esp8266 {
	esp := &esp8266{parent, "", "", "", false, true, []string{}, config{}, baseTopic, messaging.PublishMessage, map[string]func(){}}
	actionHandlers := map[string]func(){
		"reset": esp.Reset,
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

func (e *esp8266) checkOTA() {
	log.Debug("checking if parentDevice " + e.parentId + " is up to date")
	parentDevice, err := device.Get(e.parentId)
	if err != nil {
		log.Error("device " + e.parentId + " not found")
	} else if uptodate, err := ota.IsUpToDate(parentDevice.Fw.Name, parentDevice.Fw.Version); err != nil {
		log.Debug("device " + e.parentId + " is running a firmware (" + parentDevice.Fw.Name + ") which is not managed by OTA")
	} else {
		e.LastVersion = ota.LastVersion(parentDevice.Fw.Name)
		if uptodate {
			e.UpToDate = true
		} else {
			log.Info("device " + e.parentId + " is outdated!")
			e.UpToDate = false
		}

	}
}

func (e *esp8266) Set(property string, value string) {
	switch property {
	case "version":
		e.Version = value
		if e.Ota {
			e.checkOTA()
		}
	case "config":
		if err := json.Unmarshal([]byte(value), &e.Config); err != nil {
			log.Error("error while parsing device config: " + err.Error())
		}
	case "ota/enabled":
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			e.Ota = boolValue
			if e.Ota && e.Version != "" {
				e.checkOTA()
			}
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
