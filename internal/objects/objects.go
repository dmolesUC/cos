package objects

import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/streaming"
)

// The Object type represents the location of an object in cloud storage.
type Object interface {
	Protocol() string
	Endpoint() *url.URL
	Bucket() *string
	Key() *string

	Refresh()

	ContentLength() (int64, error)

	Create(body io.Reader, length int64) (err error)
	DownloadRange(startInclusive, endInclusive int64, buffer []byte) (int64, error)
	Delete() (err error)
}

func ProtocolUriStr(obj Object) string {
	uriStr := fmt.Sprintf("%v://%v/%v", obj.Protocol(), logging.PrettyStrP(obj.Bucket()), logging.PrettyStrP(obj.Key()))
	return fmt.Sprintf("%#v", uriStr)
}

func Download(obj Object, rangeSize int64, out io.Writer) (totalRead int64, err error) {
	// this will 404 if the object doesn't exist
	contentLength, err := obj.ContentLength()
	if err != nil {
		return 0, err
	}
	logger := logging.DefaultLogger()

	target := logging.NewProgressWriter(out, contentLength)
	target.LogTo(logger, time.Second)

	expectedBytes := target.ExpectedBytes()

	for ; totalRead < expectedBytes; {
		start, end, size := streaming.NextRange(totalRead, rangeSize, expectedBytes)
		buffer := make([]byte, size)
		bytesRead, err := obj.DownloadRange(start, end, buffer)
		if err != nil {
			break
		}
		err = streaming.WriteExactly(target, buffer)
		if err != nil {
			break
		}
		totalRead += bytesRead
	}
	logger.Detailf("%v from %v\n", logging.FormatBytes(totalRead), ProtocolUriStr(obj))
	return totalRead, err
}

// CalcDigest calculates the digest of the object using the specified algorithm
// (md5 or sha256), using ranged downloads of the specified size.
func CalcDigest(obj Object, downloadRangeSize int64, algorithm string) ([] byte, error) {
	h := newHash(algorithm)
	_, err := Download(obj, downloadRangeSize, h)
	if err != nil {
		return nil, err
	}
	digest := h.Sum(nil)
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
	} else if algorithm == "md5" {
		return md5.New()
	}
	log.Fatalf("unsupported digest algorithm: '%v'\n", algorithm)
	return nil
}
