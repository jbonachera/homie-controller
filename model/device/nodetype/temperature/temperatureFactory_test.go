package temperature

import (
	"github.com/jbonachera/homie-controller/model/device/nodetype"
	"testing"
)

func TestNew(t *testing.T) {
	temperature, _ := nodetype.New("temperature", "devices/")
	if temperature.GetType() != "temperature" {
		t.Error("could not get a tempetature node type")
	}
}
