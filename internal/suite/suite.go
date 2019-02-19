package suite

import (
	"fmt"
	"time"

	"github.com/dmolesUC3/cos/internal/logging"

	"github.com/dmolesUC3/cos/internal/objects"
)

type Suite interface {
	Execute() int64
}

func FileSizeSuite(sizeMax int64, target objects.Target, logLevel logging.LogLevel, dryRun bool) Suite {
	return newSuite(FileSizeCases(sizeMax), target, logLevel, dryRun)
}

func FileCountSuite(countMax uint64, target objects.Target, logLevel logging.LogLevel, dryRun bool) Suite {
	return newSuite(FileCountCases(countMax), target, logLevel, dryRun)
}

func AllCases(sizeMax int64, countMax uint64, target objects.Target, logLevel logging.LogLevel, dryRun bool) Suite {
	cases := FileSizeCases(sizeMax)
	cases = append(cases, FileCountCases(countMax)...)
	return newSuite(cases, target, logLevel, dryRun)
}

// ------------------------------------------------------------
// Unexported types

func newSuite(cases []Case, target objects.Target, logLevel logging.LogLevel, dryRun bool) Suite {
	return &suite{
		cases:    cases,
		target:   target,
		logLevel: logLevel,
		dryRun:   dryRun,
	}
}

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
		detail := c.RunWithSpinner(index, target, dryRun)
		if detail != "" && logLevel > logging.Info {
			fmt.Println(detail)
		}
	}
	return time.Now().UnixNano() - startAll
}
