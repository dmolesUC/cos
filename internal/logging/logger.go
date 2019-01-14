package logging

import (
	"fmt"
	"io"
	"log"
	"os"
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
}

// NewLogger returns a new logger, either verbose or not, as specified
func NewLogger(verbose bool) Logger {
	if verbose {
		return verboseLogger{ infoLogger {out: os.Stderr} }
	}
	return terseLogger{ infoLogger {out: os.Stderr} }
}


// ------------------------------------------------------------
// Unexported symbols

// ------------------------------
// infoLogger

// Partial base Logger implementation
type infoLogger struct {
	out io.Writer
}

// Logger.Info() implementation: log to stderr
func (l infoLogger) Info(a ...interface{}) {
	_, err := fmt.Fprintln(l.out, a...)
	if err != nil {
		log.Fatal(err)
	}
}

// Logger.Infof() implementation: log to stderr
func (l infoLogger) Infof(format string, a ...interface{}) {
	_, err := fmt.Fprintf(l.out, format, a...)
	if err != nil {
		log.Fatal(err)
	}
}

// ------------------------------
// terseLogger

// Logger implementation with no-op Detail()
type terseLogger struct {
	infoLogger
}

// No-op Logger.Detail() impelementation
func (l terseLogger) Detail(a ...interface{}) {
	// does nothing
}

// No-op Logger.Detailf() impelementation
func (l terseLogger) Detailf(format string, a ...interface{}) {
	// does nothing
}

// Logger.Verbose() implementation
func (l terseLogger) Verbose() bool {
	return false
}

func (l terseLogger) String() string {
	return "terse"
}

// ------------------------------
// verboseLogger

// Logger implementation forwarding Detail() to Info()
type verboseLogger struct {
	infoLogger
}

// Logger.Detail() implementation: forward to Info()
func (l verboseLogger) Detail(a ...interface{}) {
	l.Info(a...)
}

// Logger.Detailf() implementation: forward to Infof()
func (l verboseLogger) Detailf(format string, a ...interface{}) {
	l.Infof(format, a...)
}

// Logger.Verbose() implementation
func (l verboseLogger) Verbose() bool {
	return true
}

func (l verboseLogger) String() string {
	return "verbose"
}