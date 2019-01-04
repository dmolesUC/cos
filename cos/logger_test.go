package cos

import (
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

// ------------------------------------------------------------
// Fixture

func TestLogger(t *testing.T) { TestingT(t) }

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
