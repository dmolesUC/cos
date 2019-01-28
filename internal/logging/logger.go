package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// ------------------------------------------------------------
// Exported symbols

// The Logger interface represents a minimalist logger, inspired by:
// - https://dave.cheney.net/2015/11/05/lets-talk-about-logging
// - https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern
type Logger interface {
	Info(a ...interface{})
	Detail(a ...interface{})
	Infof(format string, a ...interface{})
	Detailf(format string, a ...interface{})
	Verbose() bool
	String() string
	TrapFatal(logFatal func (v ...interface{}))
}

// NewLogger returns a new logger, either verbose or not, as specified
func NewLogger(verbose bool) Logger {
	return NewLoggerTo(verbose, os.Stderr)
}

func NewLoggerTo(verbose bool, out io.Writer) Logger {
	if verbose {
		return &verboseLogger{ infoLogger {out: out} }
	}
	return &terseLogger{ infoLogger {out: out} }
}

// ------------------------------
// infoLogger

// Partial base Logger implementation
type infoLogger struct {
	out    io.Writer
	mux    sync.Mutex
	fatalP *func (v ...interface{})
}

func (l *infoLogger) fatal(v ...interface{}) {
	if l.fatalP == nil {
		log.Fatal(v...)
	} else {
		(*l.fatalP)(v...)
	}
}

func (l *infoLogger) TrapFatal(fatal func(v ...interface{})) {
	l.fatalP = &fatal
}

// Logger.Info() implementation: log to stderr
func (l *infoLogger) Info(a ...interface{}) {
	l.mux.Lock()
	defer l.mux.Unlock()

	pretty := Prettify(a...)
	_, err := fmt.Fprintln(l.out, pretty...)
	if err != nil {
		l.fatal(err)
	}
}

// Logger.Infof() implementation: log to stderr
func (l *infoLogger) Infof(format string, a ...interface{}) {
	l.mux.Lock()
	defer l.mux.Unlock()

	pretty := Prettify(a...)
	_, err := fmt.Fprintf(l.out, format, pretty...)
	if err != nil {
		l.fatal(err)
	}
}

// ------------------------------
// terseLogger

// Logger implementation with no-op Detail()
type terseLogger struct {
	infoLogger
}

// No-op Logger.Detail() impelementation
func (l *terseLogger) Detail(a ...interface{}) {
	// does nothing
}

// No-op Logger.Detailf() impelementation
func (l *terseLogger) Detailf(format string, a ...interface{}) {
	// does nothing
}

// Logger.Verbose() implementation
func (l *terseLogger) Verbose() bool {
	return false
}

func (l *terseLogger) String() string {
	return "terse"
}

// ------------------------------
// verboseLogger

// Logger implementation forwarding Detail() to Info()
type verboseLogger struct {
	infoLogger
}

// Logger.Detail() implementation: forward to Info()
func (l *verboseLogger) Detail(a ...interface{}) {
	l.Info(a...)
}

// Logger.Detailf() implementation: forward to Infof()
func (l *verboseLogger) Detailf(format string, a ...interface{}) {
	l.Infof(format, a...)
}

// Logger.Verbose() implementation
func (l *verboseLogger) Verbose() bool {
	return true
}

func (l *verboseLogger) String() string {
	return "verbose"
}

