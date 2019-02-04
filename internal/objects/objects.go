package objects

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"net/url"
	"time"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/streaming"
)

// ------------------------------------------------------------
// Object type

type Object interface {
	GetEndpoint() Target

	Create(body io.Reader, length int64) (err error)
	ContentLength() (length int64, err error)
	DownloadRange(startInclusive, endInclusive int64, buffer []byte) (n int64, err error)
	Delete() (err error)

	Pretty() string
}

// ------------------------------
// Factory methods

func NewObject(objURL, endpointURL *url.URL, regionStr string) (Object, error) {
	protocol := objURL.Scheme
	bucket := objURL.Host
	key := objURL.Path

	bucketUrlStr := fmt.Sprintf("%v://%v", protocol, bucket)
	bucketURL, err := url.Parse(bucketUrlStr)
	if err != nil {
		return nil, err
	}

	target, err := NewTarget(endpointURL, bucketURL, regionStr)
	if err != nil {
		return nil, err
	}
	return target.Object(key), nil
}

// ------------------------------------------------------------
// Utility functions

// Download downloads the object in chunks of the specified rangeSize, writing
// the downloaded bytes to the specified io.Writer.
func Download(obj Object, rangeSize int64, out io.Writer) (n int64, err error) {
	// this will 404 if the object doesn't exist
	contentLength, err := obj.ContentLength()
	if err != nil {
		return 0, err
	}
	logger := logging.DefaultLogger()

	outWithProgress := logging.NewProgressWriter(out, contentLength)
	outWithProgress.LogTo(logger, time.Second)

	expectedBytes := outWithProgress.ExpectedBytes()

	for ; n < expectedBytes; {
		start, end, size := streaming.NextRange(n, rangeSize, expectedBytes)
		buffer := make([]byte, size)
		bytesRead, err := obj.DownloadRange(start, end, buffer)
		if err != nil {
			break
		}
		err = streaming.WriteExactly(outWithProgress, buffer)
		if err != nil {
			break
		}
		n += bytesRead
	}
	logger.Detailf("%v from %v\n", logging.FormatBytes(n), obj)
	return n, err
}

// CalcDigest calculates the digest of the object using the specified algorithm
// (md5 or sha256), using ranged downloads of the specified size.
func CalcDigest(obj Object, downloadRangeSize int64, algorithm string) ([] byte, error) {
	h, err := newHash(algorithm)
	if err != nil {
		return nil, err
	}
	_, err = Download(obj, downloadRangeSize, h)
	if err != nil {
		return nil, err
	}
	digest := h.Sum(nil)
	return digest, nil
}

// newHash returns a new hash of the specified algorithm ("sha256" or "md5")
func newHash(algorithm string) (hash.Hash, error) {
	if algorithm == "sha256" {
		return sha256.New(), nil
	} else if algorithm == "md5" {
		return md5.New(), nil
	}
	return nil, fmt.Errorf("unsupported digest algorithm: '%v'\n", algorithm)
}
