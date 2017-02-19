package esp8266

import (
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"strconv"
	"strings"
)

type esp8266 struct {
	version   string
	ota       bool
	baseTopic string
}

func New(baseTopic string) *esp8266 {
	return &esp8266{"", false, baseTopic}
}

func (e *esp8266) GetName() string {
	return "esp8266"
}

func (e *esp8266) Set(property string, value string) {
	switch property {
	case "version":
		e.version = value
	case "ota":
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			e.ota = boolValue
		}
	}
}
func (e *esp8266) GetProperties() []string {
	return []string{"version", "ota"}
}

func (e *esp8266) MessageHandler(message homieMessage.HomieMessage) {
	// Will be bound to devices/<id>/implementation/#

	topicComponents := strings.Split(message.Path, "/")
	if len(topicComponents) < 2 {
		return
	}
	log.Debug("received update: " + topicComponents[1] + " to " + message.Payload)
	propertyName := topicComponents[1]
	e.Set(propertyName, message.Payload)
}
