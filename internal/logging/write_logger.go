package logging

import (
	"fmt"
	"io"
	"sync"
)

type writeLogger struct {
	maxlevel LogLevel
	out      io.Writer
	mux      sync.Mutex
}

func (l *writeLogger) SetMaxLevel(maxLevel LogLevel) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.maxlevel = maxLevel
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

