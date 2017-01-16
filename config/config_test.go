package config

import (
	"testing"
	"os"
)

func TestGet(t *testing.T) {
	start()
	storage["key"] = "value"
	testKey := Get("key")
	if testKey != "value"{
		t.Error("did not get a valid value: got ", testKey, ", wanted key")
	}
}
func TestGetEnv(t *testing.T) {
	start()
	os.Setenv("KEY", "value2")
	testKey := Get("key")
	if testKey != "value2"{
		t.Error("did not get a valid value: got ", testKey, ", wanted value2")
	}
}

func TestGetEnvOverride(t *testing.T) {
	start()
	storage["key"] = "value"
	os.Setenv("key", "value2")
	testKey := Get("key")
	if testKey != "value"{
		t.Error("did not get a valid value: got ", testKey, ", wanted value")
	}
}