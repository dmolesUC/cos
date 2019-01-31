package pkg

import (
	"fmt"
	"strings"

	"github.com/minimaxir/big-list-of-naughty-strings/naughtystrings"

	"github.com/dmolesUC3/cos/internal/logging"
	. "github.com/dmolesUC3/cos/internal/objects"
)

type Keys struct {
	Endpoint Target
	Source   KeySource
}

func NewKeys(target Target, source KeySource) Keys {
	return Keys{Endpoint: target, Source: source}
}

type KeyFailure struct {
	SourceName string
	Index      int
	Key        string
	Error      error
}

func (k *Keys) CheckAll(startIndex int, endIndex int) ([]KeyFailure, error) {
	source := k.Source
	sourceKeys := source.Keys()
	count := source.Count()

	var failures []KeyFailure
	for index := startIndex; index < endIndex; index ++ {
		key := sourceKeys[index]
		f, err := k.Check(source.Name(), index, count, key)
		if err != nil {
			return nil, err
		}
		if f != nil {
			failures = append(failures, *f)
		}
	}
	return failures, nil
}

func (k *Keys) Check(sourceName string, index, count int, key string) (*KeyFailure, error) {
	logger := logging.DefaultLogger()
	crvd, err := NewDefaultCrvd(k.Endpoint, key)
	if err != nil {
		return nil, err
	}
	logger.Detailf("%d of %d from %v\n", 1 + index, count, sourceName)
	err = crvd.CreateRetrieveVerifyDelete()
	if err == nil {
		return nil, nil
	}
	if strings.Contains(fmt.Sprintf("%v", err), "no such host") {
		return nil, err
	}

	msg := fmt.Sprintf("%#v (%d of %d from %v) failed: %v",
		key,
		1+index,
		count,
		sourceName,
		strings.Replace(err.Error(), "\n", "\\n", -1),
	)
	fmt.Println(msg)
	logger.Detail(msg)
	return &KeyFailure{sourceName, index, key, err}, nil
}

type KeySource interface {
	Name() string
	Keys() []string
	Count() int
}

func NaughtyStrings() KeySource {
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
		for _, k := range naughtystrings.Unencoded() {
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
