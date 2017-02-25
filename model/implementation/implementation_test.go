package implementation

import (
	"github.com/jbonachera/homie-controller/model/homieMessage"
	"testing"
)

type mockImplementation struct {
}

func (e *mockImplementation) GetName() string {
	return "mock"
}
func (e *mockImplementation) Do(action string) {
}
func (e *mockImplementation) Set(property string, value string) {
}
func (e *mockImplementation) GetProperties() []string {
	return []string{}
}

func (e *mockImplementation) MessageHandler(message homieMessage.HomieMessage) {}

type mockFactory struct {
}

func (m *mockFactory) New(parent string, baseTopic string) Implementation {
	return &mockImplementation{}
}

func TestRegisterImplementation(t *testing.T) {
	RegisterImplementation("mock", &mockFactory{})
	_, exist := implementations["mock"]
	if !exist {
		t.Error("could not found inserted implementation")
	}
}
