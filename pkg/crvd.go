package pkg

import (
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/objects"
)

const (
	DefaultContentLengthBytes = 128
	DefaultRandomSeed = 1
)

type Crvd struct {
	Object        objects.Object
	ContentLength int64
	RandomSeed    int64
}

func NewDefaultCrvd(key, endpoint, region, bucket string) (*Crvd, error) {
	return NewCrvd(key, endpoint, region, bucket, DefaultContentLengthBytes, DefaultRandomSeed)
}

func NewCrvd(key, endpoint, region, bucket string, contentLength, randomSeed int64) (*Crvd, error) {
	if key == "" {
		key = fmt.Sprintf("cos-crvd-%d.bin", time.Now().Unix())
	}
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint URL must be specified")
	}

	bucketUrl, err := objects.ValidAbsURL(bucket)
	if err != nil {
		return nil, err
	}

	obj, err := objects.NewObjectBuilder().
		WithEndpointStr(endpoint).
		WithRegion(region).
		WithProtocolUri(bucketUrl).
		WithKey(key). // TODO: fix builder so we can set this first
		Build()
	if err != nil {
		return nil, err
	}

	var crvd = Crvd{
		Object:        obj,
		ContentLength: contentLength,
		RandomSeed:    randomSeed,
	}
	return &crvd, nil
}

func (c *Crvd) CreateRetrieveVerifyDelete() error {
	err := c.CreateRetrieveVerify()
	err2 := c.Object.Delete()
	if err == nil {
		return err2
	}
	return err
}

func (c *Crvd) CreateRetrieveVerify() error {
	obj := c.Object
	logger := logging.DefaultLogger()
	protocolUriStr := objects.ProtocolUriStr(obj)

	contentLength := c.ContentLength
	logger.Detailf("Creating object at %v\n", protocolUriStr)
	expectedDigest, err := c.create()
	if err != nil {
		return err
	}
	logger.Detailf("Created %v (%d bytes)\n", protocolUriStr, contentLength)
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
		logger.Detailf("Verified %v (%d bytes, SHA-256 digest %x)\n", protocolUriStr, contentLength, actualDigest)
	}
	return err
}

func (c *Crvd) newBody() io.Reader {
	random := rand.New(rand.NewSource(c.RandomSeed))
	return io.LimitReader(random, c.ContentLength)
}

func (c *Crvd) create() ([] byte, error) {
	obj := c.Object
	logger := logging.DefaultLogger()

	digest := sha256.New()
	tr := io.TeeReader(c.newBody(), digest)

	contentLength := c.ContentLength
	in := logging.NewProgressReader(tr, contentLength)
	in.LogTo(logger, time.Second)

	err := obj.Create(in, contentLength)
	if err != nil {
		return nil, err
	}
	logger.Detailf("%v to %v\n", logging.FormatBytes(in.TotalBytes()), objects.ProtocolUriStr(obj))
	return digest.Sum(nil), err
}
