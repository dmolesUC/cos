package suite

import (
	"fmt"
	"strings"

	"github.com/dmolesUC3/cos/internal/logging"
	. "github.com/dmolesUC3/cos/pkg"

	"github.com/dmolesUC3/cos/internal/objects"
)

type seqCase struct {
	caseImpl
	allSeqs []string
}

func NewSeqCase(prefix string, seqName string, seqs []string) Case {
	c := seqCase{allSeqs: seqs}
	c.name = fmt.Sprintf("%v%v (%d sequences)", prefix, seqName, len(seqs))
	c.exec = c.doExec
	return &c
}

func (u *seqCase) doExec(target objects.Target) (ok bool, detail string) {
	invalidSeqsForKey := findInvalidSeqsForKeyIn(u.allSeqs, target)
	numInvalid := len(invalidSeqsForKey)
	if numInvalid == 0 {
		return true, ""
	}
	return false, fmt.Sprintf("%d invalid sequences: %#v", numInvalid, toMessage(invalidSeqsForKey))
}

func toMessage(invalidSeqs []string) string {
	var sb strings.Builder
	for i, s := range invalidSeqs {
		next := fmt.Sprintf("%v %X\n", s, []rune(s))
		if 1 + i < len(invalidSeqs) {
			next += ", "
		}
		// TODO: use emoji.DisplayWidth()
		if sb.Len() + len(next) > 60 {
			return sb.String()
		}
		sb.WriteString(next)
	}
	return sb.String()
}

func findInvalidSeqsForKeyIn(seqs []string, target objects.Target) []string {
	if len(seqs) == 0 {
		return nil
	}
	if lenTotal(seqs) < keyMaxBytes {
		filename := strings.Join(seqs, "")
		crvd := NewCrvd(target, filename, DefaultContentLengthBytes, DefaultRandomSeed)
		err := crvd.CreateRetrieveVerifyDelete()
		if err == nil {
			return nil
		} else {
			logging.DefaultLogger().Tracef("error creating %#v: %v\n", filename, err)
		}
		if len(seqs) == 1 {
			return seqs
		}
	}
	// Either:
	// 1. we have too many characters to test in a single key, so we split it, or
	// 2. we have one or more invalid sequences somewhere in this list, so we binary search for them
	s1, s2 := splitStrings(seqs)
	result1 := findInvalidSeqsForKeyIn(s1, target)
	result2 := findInvalidSeqsForKeyIn(s2, target)
	return append(result1, result2...)
}

func lenTotal(seqs []string) int {
	total := 0
	for _, s := range seqs {
		total += len(s)
	}
	return total
}

func splitStrings(s []string) (left, right []string) {
	left = s[0 : len(s)/2]
	right = s[len(s[0:len(s)/2]):]
	return left, right
}
