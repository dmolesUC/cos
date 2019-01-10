package objects

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/dmolesUC3/cos/internal/util"
)

// The Object type represents the location of an object in cloud storage.
type Object interface {
	Endpoint() *url.URL
	Bucket() *string
	Key() *string
	StreamDown(chunkSize int64, handleChunk func([]byte) error) (int64, error)
}

// CalcDigest calculates the digest of the object using the specified algorithm
// (md5 or sha256), using ranged downloads of the specified size.
func CalcDigest(obj Object, downloadRangeSize int64, algorithm string) ([] byte, error) {
	hash := util.NewHash(algorithm)
	_, err := obj.StreamDown(downloadRangeSize, func(bytes []byte) error {
		_, err := hash.Write(bytes)
		return err
	})
	if err != nil {
		return nil, err
	}
	digest := hash.Sum(nil)
	return digest, nil
}

// ValidAbsURL parses the specified URL string, returning an error if the
// URL cannot be parsed, or is not absolute (i.e., does not have a scheme)
func ValidAbsURL(urlStr string) (*url.URL, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return u, err
	}
	if !u.IsAbs() {
		msg := fmt.Sprintf("URL '%v' must have a scheme", urlStr)
		return nil, errors.New(msg)
	}
	return u, nil
}

