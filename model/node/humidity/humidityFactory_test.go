package humidity

import (
	"testing"
	"github.com/jbonachera/homie-controller/model/node"
)

func TestNew(t *testing.T) {
	humidity, _ := node.New("humidity", "devices/")
	if humidity.GetType() != "humidity" {
		t.Error("could not get a tempetature node type")
	}
}
