package config

import "testing"

func TestGet(t *testing.T) {
	storage["key"] = "value"
	testKey := Get("key")
	if testKey != "value"{
		t.Error("did not get a valid value: got ", testKey, ", wanted key")
	}
}
