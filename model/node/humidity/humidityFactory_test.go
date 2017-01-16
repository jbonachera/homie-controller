package humidity

import (
	"github.com/jbonachera/homie-controller/model/node"
	"testing"
)

func TestNew(t *testing.T) {
	humidity, _ := node.New("humidity", "devices/")
	if humidity.GetType() != "humidity" {
		t.Error("could not get a tempetature node type")
	}
}
