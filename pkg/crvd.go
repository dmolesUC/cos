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

func (c *Crvd) CreateRetrieveVerify(body io.Reader, contentLength int64) error {
	obj := c.Object
	logger := obj.Logger()
	protocolUriStr := objects.ProtocolUriStr(obj)

	logger.Infof("Creating object at %v\n", protocolUriStr)
	expectedDigest, err := c.create(body, contentLength)
	if err != nil {
		return err
	}
	logger.Infof("Created %v (%d bytes)\n", protocolUriStr, contentLength)
	logger.Tracef("Calculated digest on upload: %x\n", expectedDigest)

	obj.Refresh()
	var actualLength int64
	actualLength, err = obj.ContentLength()
	if err != nil {
		return fmt.Errorf("unable to determine content-length after upload: %v", err)
	}

	if actualLength != contentLength {
		return fmt.Errorf("content-length mismatch: expected: %d, actual: %d", contentLength, actualLength)
	}
	logger.Detailf("Uploaded %d bytes\n", contentLength)
	logger.Detailf("Verifying %v (expected digest: %x)\n", protocolUriStr, expectedDigest)
	check := Check{Object: obj, Expected: expectedDigest, Algorithm: "sha256"}
	actualDigest, err := check.CalcDigest()
	if err == nil {
		logger.Infof("Verified %v (%d bytes, SHA-256 digest %x)\n", protocolUriStr, contentLength, actualDigest)
	}
	return err
}

func (c *Crvd) CreateRetrieveVerifyDelete(body io.Reader, contentLength int64) error {
	err := c.CreateRetrieveVerify(body, contentLength)
	if err == nil {
		obj := c.Object
		obj.Refresh()
		err = obj.Delete()
	}
	return err
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
	logger.Detailf("%v to %v\n", logging.FormatBytes(in.TotalBytes()), objects.ProtocolUriStr(obj))
	return digest.Sum(nil), err
}
