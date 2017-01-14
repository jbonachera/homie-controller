package model

type HomieExtractableMessage interface {
	Topic() string
	Payload() []byte
}

type HomieMessage struct {
	Topic   string
	Payload string
}

func NewHomieMessage(m HomieExtractableMessage) HomieMessage {
	return HomieMessage{m.Topic(), string(m.Payload())}
}
