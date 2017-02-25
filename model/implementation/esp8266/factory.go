package esp8266

import (
	"github.com/jbonachera/homie-controller/model/implementation"
)

type Factory struct{}

func (Factory) New(parent string, baseTopic string) implementation.Implementation {
	return New(parent, baseTopic)
}

func init() {
	implementation.RegisterImplementation("esp8266", Factory{})
}
