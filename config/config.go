package config

import (
	"os"
	"strings"
)

var storage map[string]string

func init() {
	storage = map[string]string{}
	// We may load storage from file later
}

func Get(key string) string {
	value, exist := storage[key]
	if !exist {
		return os.Getenv(strings.ToUpper(key))
	} else {
		return value
	}
}
