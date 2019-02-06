package keys

import (
	"fmt"
	"math/rand"

	"golang.org/x/tools/container/intsets"
)

type KeyList interface {
	Name() string
	Desc() string
	Keys() []string
	Count() int
}

func NewKeyList(name string, desc string, keys []string) KeyList {
	return &MemoryKeyList{name, desc, keys}
}

func SamplingKeyList(origList KeyList, sampleSize int) (KeyList, error) {
	if sampleSize > origList.Count() {
		return nil, fmt.Errorf("sample size %d must be <= original list count %d", sampleSize, origList.Count())
	}
	if sampleSize == origList.Count() {
		return origList, nil
	}
	var origIndices intsets.Sparse
	for origIndices.Len() < sampleSize {
		origIndex := rand.Intn(sampleSize)
		origIndices.Insert(origIndex)
	}
	origKeys := origList.Keys()
	sampleKeys := make([]string, sampleSize)
	for sampleIndex := 0; sampleIndex < sampleSize; sampleIndex++ {
		var origIndex int
		if origIndices.TakeMin(&origIndex) {
			sampleKeys[sampleIndex] = origKeys[origIndex]
		} else {
			// should never happen
			return nil, fmt.Errorf("unable to take key for sample index %d", sampleIndex)
		}
	}
	return NewKeyList(
		fmt.Sprintf("%d from %v", sampleSize, origList.Name()),
		fmt.Sprintf("%d keys sampled from %v", sampleSize, origList.Desc()),
		sampleKeys,
	), nil
}