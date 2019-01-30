package logging

import (
	"io"
	"os"
)

// ------------------------------------------------------------
// LogLevel


// The Logger interface represents a minimalist logger, inspired by:
// - https://dave.cheney.net/2015/11/05/lets-talk-about-logging
// - https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern
type Logger interface {
	Info(a ...interface{})
	Infof(format string, a ...interface{})

	Detail(a ...interface{})
	Detailf(format string, a ...interface{})

	Trace(a ...interface{})
	Tracef(format string, a ...interface{})

	Log(lvl LogLevel, a ...interface{})
	Logf(lvl LogLevel, format string, a ...interface{})

	MaxLevel() LogLevel
	SetMaxLevel(lvl LogLevel)

	String() string
}

func NewLogger(maxLevel LogLevel) Logger {
	return NewLoggerTo(maxLevel, os.Stderr)
}

func NewLoggerTo(maxLevel LogLevel, out io.Writer) Logger {
	return &writeLogger{maxlevel: maxLevel, out: out}
}

