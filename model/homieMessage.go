package model

import "strings"

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
	return HomieMessage{m.Topic(), string(m.Payload()), baseTopic}, nil
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
