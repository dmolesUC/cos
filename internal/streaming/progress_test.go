package streaming

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

func (s *ProgressSuite) TestProgress(c *C) {
	// TODO: move this into fixture, & split assertions

	totalBytes := int64(2048 * 1024 * 16)
	nsElapsed := int64(16) * int64(time.Second)
	contentLength := totalBytes * 32

	progress := Progress { NsElapsed: nsElapsed, TotalBytes: totalBytes, ContentLength: contentLength }

	expectedKiBps := float64(2048)
	expectedBps := float64(2048 * 1024)
	expectedNsRemaining := int64(31) * nsElapsed

	c.Assert(progress.EstimatedKibPerSecond(), Equals, expectedKiBps)
	c.Assert(progress.EstimatedBytesPerSecond(), Equals, expectedBps)
	c.Assert(progress.EstimatedNsRemaining(), Equals, expectedNsRemaining)
}

func (s *ProgressSuite) TestInfoTo(c *C) {
	totalBytes := int64(2048 * 1024 * 16)
	nsElapsed := int64(16) * int64(time.Second)
	contentLength := totalBytes * 32

	progress := Progress { NsElapsed: nsElapsed, TotalBytes: totalBytes, ContentLength: contentLength }
	progress.InfoTo(s.logger)
	c.Assert(len(s.logger.Infos), Equals, 1)

	expectedKiBps := float64(2048)
	expectedNsRemaining := int64(31) * nsElapsed

	expected := fmt.Sprintf(
		"read %d of %d bytes (%0.f KiB/s; %v elapsed, %v remaining)\n",
		totalBytes, contentLength, expectedKiBps, formatNanos(nsElapsed), formatNanos(expectedNsRemaining),
	)
	c.Assert(s.logger.Infos[0], Equals, expected)
}

type CapturingLogger struct {
	Infos   []string
	Details []string
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

