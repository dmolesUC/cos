package suite

import (
	"fmt"
	"strings"

	"github.com/dmolesUC3/emoji"

	"github.com/dmolesUC3/cos/internal/logging"
	. "github.com/dmolesUC3/cos/pkg"

	"github.com/dmolesUC3/cos/internal/objects"
)

type seqCase struct {
	caseImpl
	allSeqs []string
	linear bool
}

func NewBinarySearchSeqCase(prefix string, seqName string, seqs []string) Case {
	return NewSeqCase(prefix, seqName, seqs, false)
}

func NewSeqCase(prefix string, seqName string, seqs []string, linear bool) Case {
	c := seqCase{allSeqs: seqs, linear: linear}
	c.name = fmt.Sprintf("%v%v (%d sequences)", prefix, seqName, len(seqs))
	c.exec = c.doExec
	return &c
}

func (u *seqCase) doExec(target objects.Target) (ok bool, detail string) {
	var invalidSeqsForKey []string
	if u.linear {
		invalidSeqsForKey = listInvalidSeqsForKeyIn(u.allSeqs, target)
	} else {
		invalidSeqsForKey = findInvalidSeqsForKeyIn(u.allSeqs, target)
	}
	numInvalid := len(invalidSeqsForKey)
	if numInvalid == 0 {
		return true, ""
	}
	return false, fmt.Sprintf("%d invalid sequences: %#v", numInvalid, toMessage(invalidSeqsForKey))
}

func toMessage(invalidSeqs []string) string {
	msg := ""
	for i, s := range invalidSeqs {
		next := fmt.Sprintf("%v %v", s, logging.FormatStringBytes(s))
		if 1 + i < len(invalidSeqs) {
			next += ", "
		}
		msgNext := msg + next
		if emoji.DisplayWidth(msgNext) > 60 {
			return msg + "…"
		}
		msg = msgNext
	}
	return msg
}

func listInvalidSeqsForKeyIn(seqs []string, target objects.Target) []string {
	if len(seqs) == 0 {
		return nil
	}
	var invalid []string
	for _, seq := range seqs {
		if len(seq) > keyMaxBytes {
			panic("key too long: " + logging.FormatStringBytes(seq))
		}
		crvd := NewCrvd(target, seq, DefaultContentLengthBytes, DefaultRandomSeed)
		err := crvd.CreateRetrieveVerifyDelete()
		if err != nil {
			logging.DefaultLogger().Tracef("error creating %#v: %v\n", seq, err)
			invalid = append(invalid, seq)
		}
	}
	return invalid
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
