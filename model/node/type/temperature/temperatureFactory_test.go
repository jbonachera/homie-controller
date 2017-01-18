package temperature

import (
	"github.com/jbonachera/homie-controller/model/node"
	"testing"
)

func TestNew(t *testing.T) {
	temperature, _ := node.New("temperature", "temperature", "devices/")
	if temperature.GetType() != "temperature" {
		t.Error("could not get a tempetature node type")
	}
}
