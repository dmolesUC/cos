package logging

import (
	"fmt"
	"regexp"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

const (
	nsPerSecond = int64(time.Second)
	nsPerMinute = int64(time.Minute)
	nsPerHour   = int64(time.Hour)
)

type Pretty interface {
	Pretty() string
}

func Prettify(a ...interface{}) []interface{} {
	var pretty []interface{}
	for _, v := range a {
		if p, ok := v.(Pretty); ok {
			pretty = append(pretty, p.Pretty())
		} else {
			pretty = append(pretty, v)
		}
	}
	return pretty
}

func Untabify(text string, indent string) string {
	// TODO: support multi-level indent
	return regexp.MustCompile(`(?m)^[\t ]+`).ReplaceAllString(text, indent)
}

func PrettyStrP(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0B"
	}
	return bytefmt.ByteSize(uint64(bytes))
}

func FormatNanos(ns int64) string {
	hours := ns / nsPerHour
	remainder := ns % nsPerHour
	minutes := remainder / nsPerMinute
	remainder = ns % nsPerMinute
	seconds := remainder / nsPerSecond
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
