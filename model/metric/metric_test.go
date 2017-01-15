package metric

import (
	"testing"
)

func TestNew(t *testing.T) {
	_ = New("temperature", map[string]string{"room": "living", "sensor": "sensor01"}, map[string]string{"degrees": "24.0"})
}
