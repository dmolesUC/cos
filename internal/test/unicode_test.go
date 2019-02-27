package test

import (
	"golang.org/x/text/unicode/rangetable"
	. "gopkg.in/check.v1"

	"github.com/dmolesUC3/cos/internal/suite"
)

type UnicodeSuite struct {
}

var _ = Suite(&UnicodeSuite{})

func (s *UnicodeSuite) TestNonCharacter(c *C) {
	var count = 0
	rangetable.Visit(suite.NonCharacter, func(rune) {
		count += 1
	})
	// should be exactly 66 noncharacters, per
	// https://www.unicode.org/faq/private_use.html#noncharacters
	c.Assert(count, Equals, 66)
}