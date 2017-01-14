package model

import (
	"testing"
)

type MessageMock struct {
	topic   string
	payload string
}

func (m MessageMock) Topic() string {
	return m.topic
}
func (m MessageMock) Payload() []byte {
	return []byte(m.payload)
}

func TestNew(t *testing.T) {
	messages := []MessageMock{
		MessageMock{"devices/u1234/$online", "true"},
	}
	for _, message := range messages {
		homieMessage := NewHomieMessage(message)
		if homieMessage.Topic != message.topic {
			t.Error("Expected ", message.topic, ", got ", homieMessage.Topic)
		}
		if homieMessage.Payload != message.payload {
			t.Error("Expected ", message.payload, ", got ", homieMessage.Payload)
		}
	}
}
