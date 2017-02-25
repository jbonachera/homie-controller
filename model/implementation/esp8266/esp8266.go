package esp8266

import (
	"github.com/jbonachera/homie-controller/log"
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"strconv"
	"strings"
	"github.com/jbonachera/homie-controller/messaging"
)

type esp8266 struct {
	parentId         string
	Version          string `json:"version"`
	Ota              bool `json:"ota"`
	Actions		 []string `json:"actions"`
	baseTopic        string
	MessagePublisher messagePublisherHandler `json:"-"`
	ActionHandlers   map[string]func() `json:"-"`
}

type messagePublisherHandler func(message homieMessage.HomieMessage)

func New(parent string, baseTopic string) *esp8266 {
	esp := &esp8266{parent, "", false, []string{},baseTopic, messaging.PublishMessage,map[string]func(){}}
	actionHandlers := map[string]func(){
		"reset": esp.Reset,
	}
	actions := make([]string, len(actionHandlers))

	i := 0
	for k := range actionHandlers {
		actions[i] = k
		i++
	}
	esp.ActionHandlers = actionHandlers
	esp.Actions = actions
	return esp
}

func (e *esp8266) Do(action string) {
	return
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
	case "ota":
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

	topicComponents := strings.Split(message.Path, "/")
	if len(topicComponents) < 2 {
		return
	}
	log.Debug("received update: " + topicComponents[1] + " to " + message.Payload)
	propertyName := topicComponents[1]
	e.Set(propertyName, message.Payload)
}
