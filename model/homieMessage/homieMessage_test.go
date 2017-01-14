package homieMessage

import (
	"testing"
)

type MessageMock struct {
	deviceId  string
	topic     string
	payload   string
	baseTopic string
	path      string
}

type TestIsPropertyStruct struct {
	name   string
	result bool
}

func (m MessageMock) Topic() string {
	return m.topic
}
func (m MessageMock) Payload() []byte {
	return []byte(m.payload)
}

var messages []MessageMock = []MessageMock{
	MessageMock{"u1234", "devices/u1234/$online", "true", "devices/", "$online"},
	MessageMock{"u123", "devices/u123/$online", "true", "devices/", "$online"},
	MessageMock{"u123", "homie/u123/$online", "true", "homie/", "$online"},
	MessageMock{"u123", "devices/foo/bar/u123/$online", "true", "devices/foo/bar/", "$online"},
}
var invalidMessages []MessageMock = []MessageMock{
	MessageMock{"", "devices", "true", "devices/foo/bar/", ""},
	MessageMock{"", "devices/foor/bar/", "true", "devices/foo/bar/", ""},
}

func TestNew(t *testing.T) {

	for _, message := range messages {
		homieMessage, err := New(message, message.baseTopic)
		if err != nil {
			t.Error("Error thrown: ", err)
		}
		if homieMessage.Topic != message.topic {
			t.Error("Expected ", message.topic, ", got ", homieMessage.Topic)
		}
		if homieMessage.Payload != message.payload {
			t.Error("Expected ", message.payload, ", got ", homieMessage.Payload)
		}

		if homieMessage.Path != message.path {
			t.Error("Expected ", message.path, ", got ", homieMessage.Path)
		}
	}

	for _, message := range invalidMessages {
		_, err := New(message, message.baseTopic)
		if err == nil {
			t.Error("Error not thrown")
		}
	}
}

func TestDeviceId(t *testing.T) {
	for _, message := range messages {
		homieMessage, err := New(message, message.baseTopic)
		if err != nil {
			t.Error("Error thrown: ", err)
		}
		if id := homieMessage.Id; id != message.deviceId {
			t.Error("Expected ", message.deviceId, ", got ", id)
		}
	}
}

func TestIsProperty(t *testing.T) {
	properties := []TestIsPropertyStruct{
		{"$stats/uptime", true},
		{"internet", false},
		{"$internet", false},
		{"$implementation/config", true},
	}
	for _, testProp := range properties {
		if IsProperty(testProp.name) != testProp.result {
			t.Error("Invalid property matching : ", testProp.name)
		}
	}
}
