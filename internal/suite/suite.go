package suite

import (
	"fmt"
	"log"
	"time"

	"github.com/dmolesUC3/cos/internal/logging"

	"github.com/dmolesUC3/cos/internal/objects"
)

type Suite interface {
	Execute() int64
}

func NewSuite(cases []Case, target objects.Target, logLevel logging.LogLevel, dryRun bool) Suite {
	return &suite{
		cases:    cases,
		target:   target,
		logLevel: logLevel,
		dryRun:   dryRun,
	}
}

// ------------------------------------------------------------
// Unexported types

type suite struct {
	cases    []Case
	target   objects.Target
	logLevel logging.LogLevel
	dryRun   bool
}

func (s *suite) Execute() int64 {
	cases := s.cases
	target := s.target
	logLevel := s.logLevel
	dryRun := s.dryRun

	startAll := time.Now().UnixNano()
	for index, c := range cases {
		if c == nil {
			log.Fatalf("nil case at index %d", index)
		}
		detail := c.RunWithSpinner(index, target, dryRun)
		if detail != "" && logLevel > logging.Info {
			fmt.Println(detail)
		}
	}
	return time.Now().UnixNano() - startAll
}
