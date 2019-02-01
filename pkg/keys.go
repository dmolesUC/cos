package pkg

import (
	"fmt"
	"strings"

	"github.com/dmolesUC3/cos/internal/keys"
	"github.com/dmolesUC3/cos/internal/logging"
	. "github.com/dmolesUC3/cos/internal/objects"
)

type Keys struct {
	Endpoint Target
	KeyList  keys.KeyList
}

func NewKeys(target Target, keyList keys.KeyList) Keys {
	return Keys{Endpoint: target, KeyList: keyList}
}

type KeyFailure struct {
	SourceName string
	Index      int
	Key        string
	Error      error
}

func (k *Keys) CheckAll(startIndex int, endIndex int) ([]KeyFailure, error) {
	keyList := k.KeyList
	listKeys := keyList.Keys()
	count := keyList.Count()

	var failures []KeyFailure
	for index := startIndex; index < endIndex; index ++ {
		key := listKeys[index]
		f, err := k.Check(keyList.Name(), index, count, key)
		if err != nil {
			return nil, err
		}
		if f != nil {
			failures = append(failures, *f)
		}
	}
	return failures, nil
}

func (k *Keys) Check(listName string, index, count int, key string) (*KeyFailure, error) {
	logger := logging.DefaultLogger()
	crvd, err := NewDefaultCrvd(k.Endpoint, key)
	if err != nil {
		return nil, err
	}
	logger.Detailf("%d of %d from %v\n", 1 + index, count, listName)
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
		listName,
		strings.Replace(err.Error(), "\n", "\\n", -1),
	)
	fmt.Println(msg)
	logger.Detail(msg)
	return &KeyFailure{listName, index, key, err}, nil
}

