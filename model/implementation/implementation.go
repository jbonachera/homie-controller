package implementation

import (
	"errors"
	"github.com/jbonachera/homie-controller/model/homieMessage"
)

var implementations = make(map[string]Factory)

type Implementation interface {
	GetName() string
	Set(property string, value string)
	GetProperties() []string
	MessageHandler(message homieMessage.HomieMessage)
}

type Factory interface {
	New(baseTopic string) Implementation
}

func RegisterImplementation(name string, constructor Factory) {
	_, exist := implementations[name]
	if !exist {
		implementations[name] = constructor
	}
}

func New(name string, baseTopic string) (Implementation, error) {
	if _, exist := implementations[name]; exist {
		return implementations[name].New(baseTopic), nil
	}
	return nil, errors.New("Invalid type requested: " + name)
}