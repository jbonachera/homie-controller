package system

import (
	"github.com/jbonachera/homie-controller/model/node"
	"testing"
)

func TestNew(t *testing.T) {
	system, _ := node.New("system", "system", "u1", "devices/")
	if system.GetType() != "system" {
		t.Error("could not get a system node type")
	}
}
