package pkg

import (
	"fmt"
	"strings"

	"github.com/dmolesUC3/cos/internal/logging"

	"github.com/dmolesUC3/cos/internal/keys"
)

type Keys struct {
	endpoint string
	region   string
	bucket   string
}

func NewKeys(endpoint, region, bucket string) Keys {
	return Keys{endpoint, region, bucket}
}

type KeyFailure struct {
	Source string
	Index  int
	Key    string
	Error  error
}

func (k *Keys) CheckAll(source keys.Source, startIndex int, endIndex int) ([]KeyFailure, error) {
	var failures []KeyFailure
	sourceKeys := source.Keys()
	count := source.Count()
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
	crvd, err := NewDefaultCrvd(key, k.endpoint, k.region, k.bucket)
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
