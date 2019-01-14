package logging

import (
	"strings"

	. "gopkg.in/check.v1"
)

// ------------------------------------------------------------
// Fixture

type PrettySuite struct {
	out StringableWriter
	logger Logger
}

var _ = Suite(&PrettySuite{})

func (s *PrettySuite) SetUpTest(c *C) {
	s.out = &strings.Builder{}
	logger := NewLogger(true).(verboseLogger)
	logger.out = s.out
	s.logger = logger
}

func (s *PrettySuite) TearDownTest(c *C) {
	s.out = nil
}

// ------------------------------------------------------------
// Tests

func (s *PrettySuite) TestPrettyInfo(c *C) {
	p := Prettifiable{"msg"}

	s.logger.Info(p)
	c.Assert(s.out.String(), Equals, p.Pretty() + "\n")
}

func (s *PrettySuite) TestPrettyInfof(c *C) {
	p := Prettifiable{"msg"}

	s.logger.Infof("Is %v pretty?", p)
	c.Assert(s.out.String(), Equals, "Is " + p.Pretty() + " pretty?")
}
