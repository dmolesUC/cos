package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// ------------------------------------------------------------
// Exported symbols

type LogLevel int

const (
	Info LogLevel = iota
	Detail
	Trace
	Default = Info
)

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
	String() string
}

func NewLogger(maxLevel LogLevel) Logger {
	return NewLoggerTo(maxLevel, os.Stderr)
}

func NewLoggerTo(maxLevel LogLevel, out io.Writer) Logger {
	return &writeLogger{maxlevel: maxLevel, out: out}
}

func (l LogLevel) String() string {
	if l == Info {
		return "Info"
	}
	if l == Detail {
		return "Detail"
	}
	return "Trace"
}

// ------------------------------
// Implementation

type writeLogger struct {
	maxlevel LogLevel
	out      io.Writer
	mux      sync.Mutex
}

func (l *writeLogger) Info(a ...interface{}) {
	l.Log(Info, a...)
}

func (l *writeLogger) Infof(format string, a ...interface{}) {
	l.Logf(Info, format, a...)
}

func (l *writeLogger) Detail(a ...interface{}) {
	l.Log(Detail, a...)
}

func (l *writeLogger) Detailf(format string, a ...interface{}) {
	l.Logf(Detail, format, a...)
}

func (l *writeLogger) Trace(a ...interface{}) {
	l.Log(Trace, a...)
}

func (l *writeLogger) Tracef(format string, a ...interface{}) {
	l.Logf(Trace, format, a...)
}

func (l *writeLogger) Log(lvl LogLevel, a ...interface{}) {
	if lvl > l.maxlevel {
		return
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	pretty := Prettify(a...)
	_, err := fmt.Fprintln(l.out, pretty...)
	if err != nil {
		// TODO: is this the best we can do?
		println(err)
	}
}

func (l *writeLogger) Logf(lvl LogLevel, format string, a ...interface{}) {
	if lvl > l.maxlevel {
		return
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	pretty := Prettify(a...)
	_, err := fmt.Fprintf(l.out, format, pretty...)
	if err != nil {
		// TODO: is this the best we can do?
		println(err)
	}
}

func (l *writeLogger) MaxLevel() LogLevel {
	return l.maxlevel
}

func (l *writeLogger) String() string {
	return l.maxlevel.String()
}

