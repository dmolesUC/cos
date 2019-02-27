package suite

import (
	"fmt"
	"unicode"

	"github.com/dmolesUC3/cos/internal/objects"

	"github.com/dmolesUC3/cos/internal/logging"
	. "github.com/dmolesUC3/cos/pkg"
)

const (
	maxRunesToReport = 10
)

type rangeCase struct {
	caseImpl
	allRunes []rune
}

func NewRangeTableCase(prefix string, rangeName string, rt *unicode.RangeTable) Case {
	allRunes := rangeTableToRunes(rt)
	c := rangeCase{allRunes: allRunes}
	c.name = fmt.Sprintf("%v%v (%d characters)", prefix, rangeName, len(allRunes))
	c.exec = c.doExec
	return &c
}

func (u *rangeCase) doExec(target objects.Target) (ok bool, detail string) {
	invalidRunesForKey := findInvalidRunesForKeyIn(u.allRunes, target)
	numInvalid := len(invalidRunesForKey)
	if numInvalid == 0 {
		return true, ""
	}
	var invalidRunesStr string
	if numInvalid < maxRunesToReport {
		invalidRunesStr = string(invalidRunesForKey)
	} else {
		invalidRunesStr = string(invalidRunesForKey[0:maxRunesToReport]) + "…"
	}
	return false, fmt.Sprintf("%d invalid characters: %#v", numInvalid, invalidRunesStr)
}

// TODO: parallelize this?
func findInvalidRunesForKeyIn(keyRunes []rune, target objects.Target) []rune {
	if len(keyRunes) == 0 {
		return nil
	}
	if len(keyRunes) < keyMaxBytes {
		filename := string(keyRunes)
		crvd := NewCrvd(target, filename, DefaultContentLengthBytes, DefaultRandomSeed)
		err := crvd.CreateRetrieveVerifyDelete()
		if err == nil {
			return nil
		} else {
			logging.DefaultLogger().Tracef("error creating %#v: %v\n", filename, err)
		}
		if len(keyRunes) == 1 {
			return keyRunes
		}
	}
	// Either:
	// 1. we have too many characters to test in a single key, so we split it, or
	// 2. we have one or more invalid key characters somewhere in this string, so we binary search for them
	kr1, kr2 := splitRunes(keyRunes)
	result1 := findInvalidRunesForKeyIn(kr1, target)
	result2 := findInvalidRunesForKeyIn(kr2, target)
	return append(result1, result2...)
}

func splitRunes(r []rune) (left, right []rune) {
	left = r[0 : len(r)/2]
	right = r[len(r[0:len(r)/2]):]
	return left, right
}

// TODO: use rangetable.Visit()
func rangeTableToRunes(rt *unicode.RangeTable) []rune {
	var runes []rune
	for _, r16 := range rt.R16 {
		runes = append(runes, range16ToRunes(r16)...)
	}
	for _, r32 := range rt.R32 {
		runes = append(runes, range32ToRunes(r32)...)
	}
	return runes
}

func range16ToRunes(r16 unicode.Range16) []rune {
	var runes []rune
	for cp := r16.Lo; cp <= r16.Hi; cp += r16.Stride {
		runes = append(runes, rune(cp))
	}
	return runes
}

func range32ToRunes(r32 unicode.Range32) []rune {
	var runes []rune
	for cp := r32.Lo; cp <= r32.Hi; cp += r32.Stride {
		runes = append(runes, rune(cp))
	}
	return runes
}

func isEmpty(rt *unicode.RangeTable) bool {
	return len(rt.R16) == 0 && len(rt.R32) == 0
}
