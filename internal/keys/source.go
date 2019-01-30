package keys

import (
	"strings"

	ns "github.com/minimaxir/big-list-of-naughty-strings/naughtystrings"
)

type Source interface {
	Name() string
	Keys() []string
	Count() int
}

func NaughtyStrings() Source {
	return &naughtyStrings{}
}

type naughtyStrings struct {
	keys []string
}

func (n *naughtyStrings) Name() string {
	return "minimaxir/big-list-of-naughty-strings"
}

func (n *naughtyStrings) Keys() []string {
	if len(n.keys) == 0 {
		for _, k := range ns.Unencoded() {
			if exclude(k) {
				continue
			}
			//fmt.Printf("including %d: %#v", index + 1, k)
			n.keys = append(n.keys, k)
		}
	}
	return n.keys
}

func (n *naughtyStrings) Count() int {
	return len(n.Keys())
}

func exclude(k string) bool {
	// silly
	if k == "" {
		return true
	}
	// known not to work with AWS, and dangerous
	if k == "." {
		return true
	}
	if strings.HasPrefix(k, "..") {
		return true
	}
	return false
}


