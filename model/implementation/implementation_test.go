package implementation

import (
	"testing"
	MQTT "github.com/eclipse/paho.mqtt.golang"

)

type mockImplementation struct {

}


func (e *mockImplementation) GetName() string{
	return "mock"
}

func (e *mockImplementation) Set(property string, value string){
}
func (e *mockImplementation) GetProperties() []string{
	return []string{}
}

func (e *mockImplementation) MQTTHandler(mqttClient MQTT.Client, mqttMessage MQTT.Message){}

type mockFactory struct {
}

func (m* mockFactory) New(baseTopic string)Implementation{
	return &mockImplementation{}
}

func TestRegisterImplementation(t *testing.T) {
	RegisterImplementation("mock", &mockFactory{})
	_, exist := implementations["mock"]; if !exist{
		t.Error("could not found inserted implementation")
	}
}