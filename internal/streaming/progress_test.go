package streaming

import (
	"time"

	. "gopkg.in/check.v1"

	"github.com/dmolesUC3/cos/internal/logging"
)

type ProgressSuite struct {
	logger *logging.CapturingLogger
}

var _ = Suite(&ProgressSuite{})

func (s *ProgressSuite) SetUpTest(c *C) {
	s.logger = &logging.CapturingLogger{}
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
	c.Assert(s.logger.Infos[0], Equals, "")
}