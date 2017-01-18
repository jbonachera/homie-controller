package esp8266

import (
	"github.com/jbonachera/homie-controller/model/implementation"
)

type Factory struct{}

func (Factory) New(baseTopic string) implementation.Implementation{
	return New(baseTopic)
}

func init() {
	implementation.RegisterImplementation("esp8266", Factory{})
}
