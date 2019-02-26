package suite

import (
	"sort"
	"unicode"

	"github.com/dmolesUC3/emoji"

	emojidata "github.com/dmolesUC3/emoji/data"
)

const (
	keyMaxBytes      = 1024
)

func AllUnicodeCases() []Case {
	var cases []Case
	cases = append(cases, UnicodePropertiesCases()...)
	cases = append(cases, UnicodeScriptsCases()...)
	cases = append(cases, UnicodeEmojiCases()...)
	return cases
}

func UnicodePropertiesCases() []Case {
	return rangeTablesToCases("Unicode properties: ", unicode.Properties)
}

func UnicodeScriptsCases() []Case {
	return rangeTablesToCases("Unicode scripts: ", unicode.Scripts)
}

func UnicodeEmojiCases() []Case {
	var cases []Case
	cases = append(cases, UnicodeEmojiPropertyCases()...)
	cases = append(cases, UnicodeEmojiSequenceCases()...)
	return cases
}

func UnicodeEmojiPropertyCases() []Case {
	var tables = map[string]*unicode.RangeTable{}
	for _, prop := range emojidata.AllProperties {
		rt := emoji.Latest.RangeTable(prop)
		if isEmpty(rt) {
			continue
		}
		tables[prop.String()] = rt
	}
	return rangeTablesToCases("Unicode emoji properties: ", tables)
}

func UnicodeEmojiSequenceCases() []Case {
	var sequences = map[string][]string{}
	for _, seqType := range emojidata.AllSeqTypes {
		seq := emoji.Latest.Sequences(seqType)
		if len(seq) == 0 {
			continue
		}
		sequences[seqType.String()] = seq
	}
	return sequencesToCases("Unicode emoji sequences: ", sequences)
}

func rangeTablesToCases(prefix string, tables map[string]*unicode.RangeTable) []Case {
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
		cases = append(cases, NewRangeTableCase(prefix, rangeName, rt))
	}
	return cases
}

func sequencesToCases(prefix string, sequences map[string][]string) []Case {
	var seqNames []string
	for seqName := range sequences {
		seqNames = append(seqNames, seqName)
	}
	sort.Strings(seqNames)

	var cases []Case
	for _, seqName := range seqNames {
		cases = append(cases, NewSeqCase(prefix, seqName, sequences[seqName]))
	}
	return cases
}

