package test

import (
	"fmt"
	"strings"

	. "gopkg.in/check.v1"

	. "github.com/dmolesUC3/cos/internal/logging"
)

// ------------------------------------------------------------
// Fixture

type LoggerSuite struct {
	out StringableWriter
}

func (s *LoggerSuite) newInfoLogger() Logger {
	return NewLoggerTo(Info, s.out)
}

func (s *LoggerSuite) newDetailLogger() Logger {
	return NewLoggerTo(Detail, s.out)
}

var _ = Suite(&LoggerSuite{})

func (s *LoggerSuite) SetUpTest(c *C) {
	s.out = &strings.Builder{}
}

func (s *LoggerSuite) TearDownTest(c *C) {
	s.out = nil
}

// ------------------------------------------------------------
// Tests

func (s *LoggerSuite) TestMinLevel(c *C) {
	levels := []LogLevel {Info, Detail, Trace }
	levelStrs := []string {"Info", "Detail", "Trace" }
	for idx, level := range levels {
		logger := NewLogger(level)
		c.Assert(logger.MaxLevel(), Equals, level)

		strExpected := levelStrs[idx]
		c.Assert(logger.String(), Equals, strExpected)
	}
}

func (s *LoggerSuite) TestInfoInfo(c *C) {
	var msg = "I am a log message"
	var logger = s.newInfoLogger()
	logger.Info(msg)
	c.Assert(s.out.String(), Equals, msg+"\n")
}

func (s *LoggerSuite) TestInfoDetail(c *C) {
	var msg = "I am a log message"
	var logger = s.newInfoLogger()
	logger.Detail(msg)
	c.Assert(s.out.String(), Equals, "")
}

func (s *LoggerSuite) TestInfoInfof(c *C) {
	var format = "I am a log message: %v %d"
	var expected = fmt.Sprintf(format, "text", 123)
	var logger = s.newInfoLogger()
	logger.Infof(format, "text", 123)
	c.Assert(s.out.String(), Equals, expected)
}

func (s *LoggerSuite) TestInfoDetailf(c *C) {
	var format = "I am a log message: %v %d"
	var logger = s.newInfoLogger()
	logger.Detailf(format, "text", 123)
	c.Assert(s.out.String(), Equals, "")
}

func (s *LoggerSuite) TestDetailInfo(c *C) {
	var msg = "I am a log message"
	var logger = s.newDetailLogger()
	logger.Info(msg)
	c.Assert(s.out.String(), Equals, msg+"\n")
}

func (s *LoggerSuite) TestDetailDetail(c *C) {
	var msg = "I am a log message"
	var logger = s.newDetailLogger()
	logger.Detail(msg)
	c.Assert(s.out.String(), Equals, msg+"\n")
}

func (s *LoggerSuite) TestDetailInfof(c *C) {
	var format = "I am a log message: %v %d"
	var expected = fmt.Sprintf(format, "text", 123)
	var logger = s.newDetailLogger()
	logger.Infof(format, "text", 123)
	c.Assert(s.out.String(), Equals, expected)
}

func (s *LoggerSuite) TestDetailDetailf(c *C) {
	var format = "I am a log message: %v %d"
	var expected = fmt.Sprintf(format, "text", 123)
	var logger = s.newDetailLogger()
	logger.Detailf(format, "text", 123)
	c.Assert(s.out.String(), Equals, expected)
}

