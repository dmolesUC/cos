package objects

import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"net/url"
)

// The Object type represents the location of an object in cloud storage.
type Object interface {
	Endpoint() *url.URL
	Bucket() *string
	Key() *string
	StreamDown(rangeSize int64, handleBytes func([]byte) error) (int64, error)
}

// CalcDigest calculates the digest of the object using the specified algorithm
// (md5 or sha256), using ranged downloads of the specified size.
func CalcDigest(obj Object, downloadRangeSize int64, algorithm string) ([] byte, error) {
	hash := newHash(algorithm)
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

// newHash returns a new hash of the specified algorithm ("sha256" or "md5")
func newHash(algorithm string) hash.Hash {
	if algorithm == "sha256" {
		return sha256.New()
	}
	return md5.New()
}