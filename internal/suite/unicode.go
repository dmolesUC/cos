package suite

import (
	"fmt"
	"sort"
	"unicode"

	"github.com/dmolesUC3/cos/internal/logging"

	. "github.com/dmolesUC3/cos/pkg"

	"github.com/dmolesUC3/cos/internal/objects"

	"github.com/dmolesUC3/emoji"

	emojidata "github.com/dmolesUC3/emoji/data"
)

const (
	keyMaxBytes      = 1024
	maxRunesToReport = 10
)

func AllUnicodeCases() []Case {
	var cases []Case
	cases = append(cases, UnicodePropertiesCases()...)
	cases = append(cases, UnicodeScriptsCases()...)
	cases = append(cases, UnicodeEmojiCases()...)
	return cases
}

func UnicodePropertiesCases() []Case {
	return toCases("Unicode properties: ", unicode.Properties)
}

func UnicodeScriptsCases() []Case {
	return toCases("Unicode scripts: ", unicode.Scripts)
}

func UnicodeEmojiCases() []Case {
	var tables = map[string]*unicode.RangeTable{}
	for _, prop := range emojidata.AllProperties {
		rt := emoji.Latest.RangeTable(prop)
		if isEmpty(rt) {
			continue
		}
		tables[prop.String()] = rt
	}
	return toCases("Unicode emoji: ", tables)
}

// TODO: emoji sequences

func toCases(prefix string, tables map[string]*unicode.RangeTable) []Case {
	var rangeNames []string
	for rangeName := range tables {
		rangeNames = append(rangeNames, rangeName)
	}
	sort.Strings(rangeNames)

	var cases []Case
	for _, rangeName := range rangeNames {
		rt := tables[rangeName]
		// Bad things happen if we try to cast these to runes
		if rt == unicode.Noncharacter_Code_Point {
			continue
		}
		uc := newUnicodeCase(prefix+rangeName, rangeTableToRunes(rt))
		cases = append(cases, uc)
	}
	return cases
}

type unicodeCase struct {
	caseImpl
	allRunes []rune
}

func newUnicodeCase(rangeName string, allRunes []rune) Case {
	c := unicodeCase{allRunes: allRunes}
	c.name = fmt.Sprintf("%v (%d characters)", rangeName, len(allRunes))
	c.exec = c.doExec
	return &c
}

func (u *unicodeCase) doExec(target objects.Target) (ok bool, detail string) {
	invalidRunesForKey := findInvalidRunesForKeyIn(u.allRunes, target)
	numInvalidRunes := len(invalidRunesForKey)
	if numInvalidRunes == 0 {
		return true, ""
	}
	var invalidRunesStr string
	if numInvalidRunes < maxRunesToReport {
		invalidRunesStr = string(invalidRunesForKey)
	} else {
		invalidRunesStr = string(invalidRunesForKey[0:maxRunesToReport]) + "…"
	}
	return false, fmt.Sprintf("%d invalid characters: %#v", numInvalidRunes, invalidRunesStr)
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
		runes := []rune(keyRunes)
		if len(runes) == 1 {
			return runes
		}
	}
	// Either:
	// 1. we have too many characters to test in a single key, so we split it, or
	// 2. we have one or more invalid key characters somewhere in this string, so we binary search for them
	kr1, kr2 := split(keyRunes)
	result1 := findInvalidRunesForKeyIn(kr1, target)
	result2 := findInvalidRunesForKeyIn(kr2, target)
	return append(result1, result2...)
}

func split(s []rune) (left, right []rune) {
	r := []rune(s)
	left = r[0 : len(r)/2]
	right = r[len(r[0:len(r)/2]):]
	return left, right
}

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