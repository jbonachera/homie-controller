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
}

func NewHomieMessage(m HomieExtractableMessage, baseTopic string) (HomieMessage, error) {
	message := HomieMessage{m.Topic(), string(m.Payload()), baseTopic}
	if !strings.HasPrefix(message.Topic, message.BaseTopic) {
		return HomieMessage{}, errors.New("Topic does not start with BaseTopic")
	}
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
