package keys

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/dmolesUC3/cos/internal/logging"
	ns "github.com/minimaxir/big-list-of-naughty-strings/naughtystrings"
)

const (
	DefaultKeyListName = "Default"
)

func KnownKeyLists() []KeyList {
	var names []string
	for name := range knownListsByName {
		names = append(names, name)
	}
	sort.Strings(names)

	var lists []KeyList
	for _, name := range names {
		lists = append(lists, knownListsByName[name])
	}
	return lists
}

func KeyListForName(name string) (KeyList, error) {
	if s, ok := knownListsByName[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("no such source: %#v", name)
}

func KeyListForFile(path string) (KeyList, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	infile, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := infile.Close(); err != nil {
			logging.DefaultLogger().Infof("error closing file %v: %v", absPath, err.Error())
		}
	}()

	var keys []string
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		keys = append(keys, scanner.Text())
	}
	return NewKeyList(path, absPath, keys), nil
}

// ------------------------------------------------------------
// init()

const defaultFilter = "^$|^\\.$|^\\.{2,}"
const backslash = "\\\\"
const doubleBackslash = backslash + backslash

var knownListsByName map[string]KeyList

func init() {
	knownListsByName = map[string]KeyList{}

	naughtyStringSource := NewKeyList(
		"naughty-strings",
		"Big List of Naughty Strings (https://github.com/minimaxir/big-list-of-naughty-strings)",
		ns.Unencoded(),
	)
	addKeyList(naughtyStringSource)

	defaultSource := NewKeyList(
		DefaultKeyListName,
		fmt.Sprintf("default source (as naughty-strings, filtering out %#v, %#v, and leading %#v)", "", ".", ".."),
		filterKeys(naughtyStringSource.Keys(), defaultFilter),
	)
	addKeyList(defaultSource)

	disallowBackslash := NewKeyList(
		"disallow-backslash",
		"as default source, disallowing backlash",
		filterKeys(defaultSource.Keys(), backslash),
	)
	addKeyList(disallowBackslash)

	disallowDoubleBackslash := NewKeyList(
		"disallow-double-backslash",
		"as default source, disallowing double backlash",
		filterKeys(defaultSource.Keys(), doubleBackslash),
	)
	addKeyList(disallowDoubleBackslash)

	misc := NewKeyList(
		"misc",
		"miscellenous potential problems, incl. path elements & unicode blocks",
		MiscKeys(),
	)
	addKeyList(misc)
}

func filterKeys(keys []string, re string) []string {
	regexpP := regexp.MustCompile(re)
	var filtered []string
	for _, k := range keys {
		if !regexpP.MatchString(k) {
			filtered = append(filtered, k)
		}
	}
	return filtered
}

func addKeyList(list KeyList) {
	name := list.Name()
	if _, ok := knownListsByName[name]; ok {
		panic(fmt.Sprintf("list %#v already exists", name))
	}
	knownListsByName[name] = list
}

