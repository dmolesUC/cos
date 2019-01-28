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
	fatals []string
}

func (s *LoggerSuite) newTerseLogger() Logger {
	return NewLoggerTo(false, s.out)
}

func (s *LoggerSuite) newVerboseLogger() Logger {
	return NewLoggerTo(true, s.out)
}

var _ = Suite(&LoggerSuite{})

func (s *LoggerSuite) SetUpTest(c *C) {
	s.out = &strings.Builder{}
	s.fatals = nil
}

func (s *LoggerSuite) TearDownTest(c *C) {
	s.out = nil
	s.fatals = nil
}

func (s *LoggerSuite) fatal(v ...interface{}) {
	s.fatals = append(s.fatals, fmt.Sprint(v...))
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

func (s *LoggerSuite) TestTerseInfoFail(c *C) {
	s.out = FailWriter{}
	var logger = s.newTerseLogger()
	logger.TrapFatal(s.fatal)
	logger.Info("I am a log message")
	c.Assert(len(s.fatals), Equals, 1)
}

func (s *LoggerSuite) TestTerseInfofFail(c *C) {
	s.out = FailWriter{}
	var logger = s.newTerseLogger()
	logger.TrapFatal(s.fatal)
	var format = "I am a log message: %v %d"
	logger.Infof(format, "text", 123)
	c.Assert(len(s.fatals), Equals, 1)
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

func (s *LoggerSuite) TestVerboseFlag(c *C) {
	flags := []bool { true, false }
	for _, verbose := range flags {
		logger := NewLogger(verbose)
		c.Assert(logger.Verbose(), Equals, verbose)

		var strExpected string
		if verbose {
			strExpected = "verbose"
		} else {
			strExpected = "terse"
		}
		c.Assert(logger.String(), Equals, strExpected)
	}
}


