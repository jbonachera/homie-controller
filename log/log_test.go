package log

import "testing"

func TestInfo(t *testing.T) {
	Info("test")
}

func TestSetLogLevel(t *testing.T) {
	SetLogLevel("DEBUG")
	if defaultLogger.getLogLevel() != DEBUG {
		t.Error("could not set log level to DEBUG")
	}
	if defaultLogger.getLogLevel() < INFO {
		t.Error("INFO loglevel is lower than DEBUG")
	}
	SetLogLevel("ERROR")
	if defaultLogger.getLogLevel() != ERROR {
		t.Error("could not set log level to ERROR")
	}
	if defaultLogger.getLogLevel() > INFO {
		t.Error("INFO loglevel is bigger than ERROR")
	}
}