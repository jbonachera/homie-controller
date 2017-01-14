package model

import (
	"errors"
	"strings"
)

type HomieExtractableMessage interface {
	Topic() string
	Payload() []byte
}

type HomieMessage struct {
	Topic     string
	Payload   string
	BaseTopic string
	Id        string
}

func NewHomieMessage(m HomieExtractableMessage, baseTopic string) (HomieMessage, error) {
	message := HomieMessage{m.Topic(), string(m.Payload()), baseTopic, ""}
	strippedPrefix := strings.TrimPrefix(message.Topic, message.BaseTopic)
	if strippedPrefix == message.Topic {
		return HomieMessage{}, errors.New("Topic does not start with BaseTopic")
	}
	message.Id = strings.Split(strippedPrefix, "/")[0]
	return message, nil
}

func (m HomieMessage) DeviceId() string {
	// Remove the baseTopic from the topic by reading the string
	// juste after it.
	if strings.Contains(m.Topic, m.BaseTopic) {
		strippedBase := m.Topic[len(m.BaseTopic):]
		return strings.Split(strippedBase, "/")[0]
	}
	return ""
}
