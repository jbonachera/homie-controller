package humidity

import (
	"github.com/jbonachera/homie-controller/model/node"
	"testing"
)

func TestNew(t *testing.T) {
	humidity, _ := node.New("humidity", "humidity", "u1", "devices/")
	if humidity.GetType() != "humidity" {
		t.Error("could not get a tempetature node type")
	}
}
