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

func NewHomieMessage(m HomieExtractableMessage, baseTopic string) HomieMessage {
	return HomieMessage{m.Topic(), string(m.Payload()), baseTopic}
}

func (m HomieMessage) DeviceId() string {
	// Remove the baseTopic from the topic by reading the string
	// juste after it.
	strippedBase := m.Topic[len(m.BaseTopic):]
	return strings.Split(strippedBase, "/")[0]
}
