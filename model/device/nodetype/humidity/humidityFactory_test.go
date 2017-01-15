package humidity

import (
	"github.com/jbonachera/homie-controller/model/device/nodetype"
	"testing"
)

func TestNew(t *testing.T) {
	humidity, _ := nodetype.New("humidity", "devices/")
	if humidity.GetType() != "humidity" {
		t.Error("could not get a tempetature node type")
	}
}
