package test

import (
	"strings"

	"golang.org/x/text/unicode/rangetable"
	. "gopkg.in/check.v1"

	"github.com/dmolesUC3/cos/internal/logging"

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

func (s *UnicodeSuite) TestUTF8InvalidSequences(c *C) {
	const badChar = rune(0xfffd)
	for i, bb := range suite.UTF8InvalidSequences {
		bytesStr := logging.FormatByteArray(bb)
		asString := string(bb)
		c.Check(strings.ContainsRune(asString, badChar), Equals, true,
			Commentf("%d: %v: expected %#x (%#v), got %#v", i, bytesStr, badChar, string(badChar), asString))
	}
}
