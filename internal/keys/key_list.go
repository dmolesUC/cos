package keys

import (
	"fmt"
	"regexp"
	"sort"

	ns "github.com/minimaxir/big-list-of-naughty-strings/naughtystrings"
)

const (
	DefaultKeyListName = "Default"
)

// ------------------------------------------------------------
// KeyList

type KeyList interface {
	Name() string
	Desc() string
	Keys() []string
	Count() int
}

func DefaultKeyList() KeyList {
	return listsByName[DefaultKeyListName]
}

func AllKeyLists() []KeyList {
	var names []string
	for name := range listsByName {
		names = append(names, name)
	}
	sort.Strings(names)

	var lists []KeyList
	for _, name := range names {
		lists = append(lists, listsByName[name])
	}
	return lists
}

func KeyListForName(name string) (KeyList, error) {
	if s, ok := listsByName[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("no such source: %#v", name)
}

func newKeyList(name string, desc string, keys []string) KeyList {
	if _, ok := listsByName[name]; ok {
		panic(fmt.Sprintf("source %#v already exists", name))
	}
	source := &keyList{name, desc, keys}
	listsByName[name] = source
	return source
}

// ------------------------------------------------------------
// init()

const defaultFilter = "^$|^\\.$|^\\.{2,}"
const backslash = "\\\\"

var listsByName map[string]KeyList

func init() {
	listsByName = map[string]KeyList{}

	naughtyStringSource := newKeyList(
		"naughty-strings",
		"Big List of Naughty Strings (https://github.com/minimaxir/big-list-of-naughty-strings)",
		ns.Unencoded(),
	)
	defaultSource := newKeyList(
		DefaultKeyListName,
		fmt.Sprintf("Default source (naughty-strings, filtering out /%v/)", defaultFilter),
		FilterKeys(naughtyStringSource.Keys(), defaultFilter),
	)
	disallowBackslash := newKeyList(
		"disallow-backslash",
		"default source, disallowing backlash",
		FilterKeys(defaultSource.Keys(), backslash),
	)
	_ = disallowBackslash
}

func FilterKeys(keys []string, re string) []string {
	regexpP := regexp.MustCompile(re)
	var filtered []string
	for _, k := range keys {
		if !regexpP.MatchString(k) {
			filtered = append(filtered, k)
		}
	}
	return filtered
}
