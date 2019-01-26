package pkg

import (
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/objects"
)

type Crvd struct {
	Object objects.Object
}

func (c *Crvd) CreateRetrieveVerify(body io.Reader, contentLength int64) (digest []byte, err error) {
	obj := c.Object
	logger := obj.Logger()
	logger.Infof("creating object at %v\n", objects.ProtocolUriStr(obj))
	digest, err = c.create(body, contentLength)
	if err == nil {
		logger.Detailf("calculated digest on upload: %x\n", digest)

		obj.Refresh()
		var actualLength int64
		actualLength, err = obj.ContentLength()
		if err == nil {
			if actualLength != contentLength {
				return digest, fmt.Errorf("content-length mismatch: expected: %d, actual: %d", contentLength, actualLength)
			}
			logger.Detailf("uploaded %d bytes\n", contentLength)
			logger.Infof("verifying %v (expected digest: %x)\n", objects.ProtocolUriStr(obj), digest)
			check := Check{Object: obj, Expected: digest, Algorithm: "sha256"}
			return check.CalcDigest()
		} else {
			err = fmt.Errorf("unable to determine content-length after upload: %v", err)
		}
	}
	return digest, err
}

func (c *Crvd) CreateRetrieveVerifyDelete(body io.Reader, contentLength int64) (digest []byte, err error) {
	digest, err = c.CreateRetrieveVerify(body, contentLength)
	if err == nil {
		obj := c.Object
		logger := obj.Logger()

		obj.Refresh()
		logger.Detailf("verified %v\n", objects.ProtocolUriStr(obj))
		err = obj.Delete()
	}
	return digest, err
}

func (c *Crvd) create(body io.Reader, contentLength int64) ([] byte, error) {
	obj := c.Object
	logger := obj.Logger()

	digest := sha256.New()
	tr := io.TeeReader(body, digest)

	in := logging.NewProgressReader(tr, contentLength)
	in.LogTo(logger, time.Second)

	err := obj.Create(in, contentLength)
	if err != nil {
		return nil, err
	}
	logger.Infof("%v to %v\n", logging.FormatBytes(in.TotalBytes()), objects.ProtocolUriStr(obj))
	return digest.Sum(nil), err
}
