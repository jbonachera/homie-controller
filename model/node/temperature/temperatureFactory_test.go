package temperature

import (
	"testing"
	"github.com/jbonachera/homie-controller/model/node"
)

func TestNew(t *testing.T) {
	temperature, _ := node.New("temperature", "devices/")
	if temperature.GetType() != "temperature" {
		t.Error("could not get a tempetature node type")
	}
}
