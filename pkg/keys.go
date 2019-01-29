package pkg

import (
	ns "github.com/minimaxir/big-list-of-naughty-strings/naughtystrings"

	"github.com/dmolesUC3/cos/internal/logging"
)

type Keys struct {
	endpoint string
	region   string
	bucket   string
	logger   logging.Logger
}

func NewKeys(endpoint, region, bucket string, logger logging.Logger) Keys {
	return Keys{endpoint, region, bucket, logger}
}

type KeyFailure struct {
	Source string
	Index  int
	Key    string
	Error  error
}

func (k *Keys) TotalKeys() int {
	return len(ns.Unencoded())
}

func (k *Keys) CheckAll() ([]KeyFailure, error) {
	source := ns.Unencoded()
	sourceName := "minimaxir/big-list-of-naughty-strings"
	var failures []KeyFailure
	for index, key := range source {
		f, err := k.Check(sourceName, index, key)
		if err != nil {
			return nil, err
		}
		if f != nil {
			failures = append(failures, *f)
		}
	}
	return failures, nil
}

func (k *Keys) Check(source string, index int, key string) (*KeyFailure, error) {
	endpoint := k.endpoint
	region := k.region
	bucket := k.bucket
	logger := k.logger
	crvd, err := NewDefaultCrvd(key, endpoint, region, bucket, logger)
	if err != nil {
		return nil, err
	}
	err = crvd.CreateRetrieveVerifyDelete()
	if err != nil {
		k.logger.Infof("%#v (%d of %v): %v\n", key, index, source, err)
		return &KeyFailure{source, index, key, err}, nil
	}
	return nil, err
}
