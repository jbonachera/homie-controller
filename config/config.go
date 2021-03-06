package config

import (
	"os"
	"strings"
)

var storage map[string]string

func init() {
	start()
}
func start() {
	storage = map[string]string{}
	// We may load storage from file later
}

func Get(key string) string {
	// Priority: storage > env
	value, exist := storage[key]
	if !exist {
		return os.Getenv(strings.ToUpper(key))
	} else {
		return value
	}
}

func Set(key string, value string) {
	storage[key] = value
}
