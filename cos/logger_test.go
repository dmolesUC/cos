package cos

import (
	"strings"
	"fmt"

	. "gopkg.in/check.v1"
)

// ------------------------------------------------------------
// Fixture

type LoggerSuite struct {
	out *strings.Builder
}

func (s *LoggerSuite) newTerseLogger() terseLogger {
	logger := NewLogger(false).(terseLogger)
	logger.out = s.out
	return logger
}

func (s *LoggerSuite) newVerboseLogger() verboseLogger {
	logger := NewLogger(true).(verboseLogger)
	logger.out = s.out
	return logger
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

func (s *LoggerSuite) TestTerseInfo(c *C) {
	var msg = "I am a log message"
	var logger = s.newTerseLogger()
	logger.Info(msg)
	c.Assert(s.out.String(), Equals, msg+"\n")
}

func (s *LoggerSuite) TestTerseDetail(c *C) {
	var msg = "I am a log message"
	var logger = s.newTerseLogger()
	logger.Detail(msg)
	c.Assert(s.out.String(), Equals, "")
}

func (s *LoggerSuite) TestTerseInfof(c *C) {
	var format = "I am a log message: %v %d"
	var expected = fmt.Sprintf(format, "text", 123)
	var logger = s.newTerseLogger()
	logger.Infof(format, "text", 123)
	c.Assert(s.out.String(), Equals, expected)
}

func (s *LoggerSuite) TestTerseDetailf(c *C) {
	var format = "I am a log message: %v %d"
	var logger = s.newTerseLogger()
	logger.Detailf(format, "text", 123)
	c.Assert(s.out.String(), Equals, "")
}

func (s *LoggerSuite) TestVerboseInfo(c *C) {
	var msg = "I am a log message"
	var logger = s.newVerboseLogger()
	logger.Info(msg)
	c.Assert(s.out.String(), Equals, msg+"\n")
}

func (s *LoggerSuite) TestVerboseDetail(c *C) {
	var msg = "I am a log message"
	var logger = s.newVerboseLogger()
	logger.Detail(msg)
	c.Assert(s.out.String(), Equals, msg+"\n")
}

func (s *LoggerSuite) TestVerboseInfof(c *C) {
	var format = "I am a log message: %v %d"
	var expected = fmt.Sprintf(format, "text", 123)
	var logger = s.newVerboseLogger()
	logger.Infof(format, "text", 123)
	c.Assert(s.out.String(), Equals, expected)
}

func (s *LoggerSuite) TestVerboseDetailf(c *C) {
	var format = "I am a log message: %v %d"
	var expected = fmt.Sprintf(format, "text", 123)
	var logger = s.newVerboseLogger()
	logger.Detailf(format, "text", 123)
	c.Assert(s.out.String(), Equals, expected)
}
