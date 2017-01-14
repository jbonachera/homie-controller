package model

import (
	"testing"
)

type MessageMock struct {
	deviceId  string
	topic     string
	payload   string
	baseTopic string
}

func (m MessageMock) Topic() string {
	return m.topic
}
func (m MessageMock) Payload() []byte {
	return []byte(m.payload)
}

var messages []MessageMock = []MessageMock{
	MessageMock{"u1234", "devices/u1234/$online", "true", "devices/"},
	MessageMock{"u123", "devices/u123/$online", "true", "devices/"},
	MessageMock{"u123", "homie/u123/$online", "true", "homie/"},
	MessageMock{"u123", "devices/foo/bar/u123/$online", "true", "devices/foo/bar/"},
	MessageMock{"", "devices", "true", "devices/foo/bar/"},
	MessageMock{"", "devices/foor/bar/", "true", "devices/foo/bar/"},
}

func TestNew(t *testing.T) {

	for _, message := range messages {
		homieMessage, err := NewHomieMessage(message, message.baseTopic)
		if homieMessage.Topic != message.topic {
			t.Error("Expected ", message.topic, ", got ", homieMessage.Topic)
		}
		if homieMessage.Payload != message.payload {
			t.Error("Expected ", message.payload, ", got ", homieMessage.Payload)
		}
	}
}

func TestDeviceId(t *testing.T) {
	for _, message := range messages {
		homieMessage, err := NewHomieMessage(message, message.baseTopic)
		if id := homieMessage.DeviceId(); id != message.deviceId {
			t.Error("Expected ", message.deviceId, ", got ", id)
		}
	}
}
