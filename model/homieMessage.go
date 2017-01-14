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
	Path      string
}

func NewHomieMessage(m HomieExtractableMessage, baseTopic string) (HomieMessage, error) {
	message := HomieMessage{m.Topic(), string(m.Payload()), baseTopic, "", ""}
	strippedPrefix := strings.TrimPrefix(message.Topic, message.BaseTopic)
	if strippedPrefix == message.Topic {
		return HomieMessage{}, errors.New("Topic does not start with BaseTopic")
	}
	message.Id = strings.Split(strippedPrefix, "/")[0]
	message.Path = strings.SplitN(strippedPrefix, "/", 2)[1]
	return message, nil
}
