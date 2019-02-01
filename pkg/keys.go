package pkg

import (
	"fmt"
	"io"
	"strings"

	. "github.com/dmolesUC3/cos/internal/keys"
	"github.com/dmolesUC3/cos/internal/logging"
	. "github.com/dmolesUC3/cos/internal/objects"
)

type Keys struct {
	Endpoint Target
	KeyList  KeyList
}

func NewKeys(target Target, keyList KeyList) Keys {
	return Keys{
		Endpoint: target,
		KeyList:  keyList,
	}
}

func (k *Keys) CheckAll(startIndex int, endIndex int, okOut io.Writer, badOut io.Writer, raw bool) ([]KeyResult, error) {
	if okOutC, ok := okOut.(io.WriteCloser); ok {
		//noinspection GoUnhandledErrorResult
		defer okOutC.Close()
	}
	if badOutC, bad := badOut.(io.WriteCloser); bad {
		//noinspection GoUnhandledErrorResult
		defer badOutC.Close()
	}

	logger := logging.DefaultLogger()

	var failures []KeyResult
	for index := startIndex; index < endIndex; index ++ {
		key := k.KeyList.Keys()[index]
		err := k.Check(key)
		if err != nil && strings.Contains(err.Error(), "no such host") {
			// network problem, or we ran out of file handles
			return nil, err
		}
		result := &KeyResult{
			List:  k.KeyList,
			Index: index,
			Key:   key,
			Error: err,
		}
		logger.Detailf(result.Pretty())

		if result.Success() {
			err = writeKey(okOut, key, raw)
			if err != nil {
				return nil, err
			}
		} else {
			failures = append(failures, *result)
			err = writeKey(badOut, key, raw)
			if err != nil {
				return nil, err
			}
		}
	}
	return failures, nil
}

func (k *Keys) Check(key string) (err error) {
	crvd, err := NewDefaultCrvd(k.Endpoint, key)
	if err != nil {
		return err
	}
	return crvd.CreateRetrieveVerifyDelete()
}

func writeKey(w io.Writer, key string, raw bool) (err error) {
	if w == nil {
		return
	}
	if raw {
		_, err = fmt.Fprintln(w, key)
	} else {
		_, err = fmt.Fprintf(w, "%#v\n", key)
	}
	return
}
