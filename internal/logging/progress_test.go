package logging

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"
)

type ProgressSuite struct {
	logger *CapturingLogger
}

var _ = Suite(&ProgressSuite{})

func (s *ProgressSuite) SetUpTest(c *C) {
	s.logger = &CapturingLogger{}
}

func (s *ProgressSuite) TearDownTest(c *C) {
	s.logger = nil
}

func (s *ProgressSuite) TestInfoTo(c *C) {
	totalBytes := int64(2048 * 1024 * 16)
	nsElapsed := int64(16) * int64(time.Second)
	contentLength := totalBytes * 32

	progress := Progress { NsElapsed: nsElapsed, TotalBytes: totalBytes, ContentLength: contentLength }
	progress.InfoTo(s.logger)
	c.Assert(len(s.logger.Infos), Equals, 1)

	expected :=  "32M of 1G (2M/s; 16s elapsed, 8m 16s remaining)\n"
	c.Assert(s.logger.Infos[0], Equals, expected)
}

// ------------------------------------------------------------
// Helper types

type CapturingLogger struct {
	Infos   []string
	Details []string
}

func (logger *CapturingLogger) String() string {
	return "capturing"
}

func (logger *CapturingLogger) Info(a ...interface{}) {
	logger.Infos = append(logger.Infos, fmt.Sprintln(a...))
}

func (logger *CapturingLogger) Infof(format string, a ...interface{}) {
	logger.Infos = append(logger.Infos, fmt.Sprintf(format, a...))
}

func (logger *CapturingLogger) Detail(a ...interface{}) {
	logger.Details = append(logger.Details, fmt.Sprintln(a...))
}

func (logger *CapturingLogger) Detailf(format string, a ...interface{}) {
	logger.Details = append(logger.Details, fmt.Sprintf(format, a...))
}

func (logger *CapturingLogger) Verbose() bool {
	return true
}



