package homieMessage

import (
	"errors"
	"strings"
)

var Properties []string = []string{
	"$homie",
	"$online",
	"$name",
	"$localip",
	"$mac",
	"$stats/uptime",
	"$stats/signal",
	"$stats/interval",
	"$fw/name",
	"$fw/version",
	"$fw/checksum",
	"$implementation",
	"$implementation/+",
}

type SubscriptibleClient interface {
	Subscribe(topic string, qos byte, callback func(SubscriptibleClient, HomieExtractableMessage)) interface{}
}
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

func New(id string, baseTopic string, path string, payload string) (HomieMessage, error) {
	if !strings.HasSuffix(baseTopic, "/") {
		baseTopic = baseTopic + "/"
	}

	topic := baseTopic + id + "/" + path
	return HomieMessage{Topic: topic, Path: path, BaseTopic: baseTopic, Id: id, Payload: payload}, nil
}

func Extract(m HomieExtractableMessage, baseTopic string) (HomieMessage, error) {
	message := HomieMessage{m.Topic(), string(m.Payload()), baseTopic, "", ""}
	strippedPrefix := strings.TrimPrefix(message.Topic, message.BaseTopic)
	if strippedPrefix == message.Topic {
		return HomieMessage{}, errors.New("Topic does not start with BaseTopic")
	}
	message.Id = strings.Split(strippedPrefix, "/")[0]
	message.Path = strings.SplitN(strippedPrefix, "/", 2)[1]
	return message, nil
}

func IsProperty(prop string) bool {
	// https://github.com/marvinroger/homie#device-properties
	switch prop {
	case
		"$homie",
		"$online",
		"$name",
		"$localip",
		"$mac",
		"$stats/uptime",
		"$stats/signal",
		"$stats/interval",
		"$fw/name",
		"$fw/version",
		"$fw/checksum",
		"$implementation":
		return true
	}
	return strings.HasPrefix(prop, "$implementation/")
}
