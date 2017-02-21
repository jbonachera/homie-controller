package intercom

import (
	"github.com/jbonachera/homie-controller/model/node"
	"testing"
)

func TestNew(t *testing.T) {
	intercom, _ := node.New("intercom", "intercom", "devices/")
	if intercom.GetType() != "intercom" {
		t.Error("could not get a intercom node type")
	}
}
