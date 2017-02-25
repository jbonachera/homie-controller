package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
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
func (m MessageMock) Duplicate() bool {
	return false
}
func (m MessageMock) Qos() byte {
	return byte(0)
}
func (m MessageMock) Retained() bool {
	return true
}
func (m MessageMock) MessageID() uint16 {
	return uint16(10)
}

type MockClient struct {
	isConnected bool
	Topic       string
	PublishedMessage []MessageMock
}

func (mc *MockClient) IsConnected() bool {
	return mc.isConnected
}

func (mc *MockClient) Connect() MQTT.Token {
	return nil
}

func (mc *MockClient) Disconnect(quiesce uint) {
	return
}

func (mc *MockClient) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	payload_string := payload.(string)
	mc.PublishedMessage = append(mc.PublishedMessage, MessageMock{topic:topic, payload:payload_string})
	return nil
}
func (mc *MockClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	mc.Topic = topic
	return nil
}

func (mc *MockClient) SubscribeMultiple(filters map[string]byte, callback MQTT.MessageHandler) MQTT.Token {
	return nil
}

func (mc *MockClient) Unsubscribe(topics ...string) MQTT.Token {
	mc.Topic = ""
	return nil
}

func NewMessage(topic string, payload string) MessageMock {
	return MessageMock{topic, payload}
}
func NewMockClient(connected bool, topic string) *MockClient {
	return &MockClient{connected, topic, []MessageMock{}}
}
