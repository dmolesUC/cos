package logging

import (
	"os"
)

var defaultLogger Logger

func DefaultLogger() Logger {
	if defaultLogger == nil {
		SetDefaultLogger(&writeLogger{maxlevel: Default, out: os.Stderr})
	}
	return defaultLogger
}

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

func DefaultLoggerWithLevel(lvl LogLevel) Logger {
	logger := DefaultLogger()
	logger.SetMaxLevel(lvl)
	return logger
}