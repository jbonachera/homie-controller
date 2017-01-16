package log

import (
	"log"
	"sync"
)

const (
	DEBUG = iota
	INFO  = iota
	WARN  = iota
	ERROR = iota
)

type customLogger struct {
	sync.Mutex
	logLevel int
}

var defaultLogger *customLogger

func init() {
	defaultLogger = &customLogger{sync.Mutex{}, WARN}
}

func (c *customLogger) println(message string) {
	c.Lock()
	defer c.Unlock()
	log.Println(message)
}

func (c *customLogger) getLogLevel() int {
	c.Lock()
	defer c.Unlock()
	return c.logLevel
}

func (c *customLogger) setLogLevel(level int) {
	c.Lock()
	defer c.Unlock()
	if level < ERROR && level > DEBUG {
		c.logLevel = level
	}
}
func SetLogLevel(loglevel string) {
	switch loglevel {
	case "DEBUG":
		defaultLogger.setLogLevel(DEBUG)
	case "INFO":
		defaultLogger.setLogLevel(INFO)
	case "WARN":
		defaultLogger.setLogLevel(WARN)
	case "ERROR":
		defaultLogger.setLogLevel(ERROR)
	}
}

func Debug(message string) {
	if defaultLogger.getLogLevel() >= DEBUG {
		defaultLogger.println("DEBUG: " + message)
	}
}

func Info(message string) {
	if defaultLogger.getLogLevel() >= INFO {
		defaultLogger.println("INFO: " + message)
	}
}

func Warn(message string) {
	if defaultLogger.getLogLevel() >= WARN {
		defaultLogger.println("WARN: " + message)
	}
}

func Error(message string) {
	if defaultLogger.getLogLevel() >= ERROR {
		defaultLogger.println("ERROR: " + message)
	}
}
