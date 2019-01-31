package pkg

import (
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"time"

	. "github.com/dmolesUC3/cos/internal/objects"

	"github.com/dmolesUC3/cos/internal/logging"
)

const (
	DefaultContentLengthBytes = 128
	DefaultRandomSeed         = 1
)

type Crvd struct {
	Object        Object
	ContentLength int64
	RandomSeed    int64
}

func NewDefaultCrvd(target Target, key string) (*Crvd, error) {
	return NewCrvd(target, key, DefaultContentLengthBytes, DefaultRandomSeed)
}

func NewCrvd(target Target, key string, contentLength int64, randomSeed int64) (*Crvd, error) {
	if key == "" {
		key = fmt.Sprintf("cos-crvd-%d.bin", time.Now().Unix())
	}
	obj := target.Object(key)
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
	contentLength := c.ContentLength

	logger := logging.DefaultLogger()
	logger.Detailf("Creating object (%v) at %v\n", logging.FormatBytes(contentLength), obj)
	expectedDigest, err := c.create()
	if err != nil {
		return err
	}
	logger.Detailf("Created %v (%d bytes)\n", obj, contentLength)
	logger.Tracef("Calculated digest on upload: %x\n", expectedDigest)

	var actualLength int64
	actualLength, err = obj.ContentLength()
	if err != nil {
		return fmt.Errorf("unable to determine content-length after upload: %v", err)
	}

	if actualLength != contentLength {
		return fmt.Errorf("content-length mismatch: expected: %d, actual: %d", contentLength, actualLength)
	}
	logger.Detailf("Uploaded %d bytes\n", contentLength)
	logger.Detailf("Verifying %v (expected digest: %x)\n", obj, expectedDigest)
	check := Check{Object: obj, Expected: expectedDigest, Algorithm: "sha256"}
	actualDigest, err := check.VerifyDigest()
	if err == nil {
		logger.Detailf("Verified %v (%d bytes, SHA-256 digest %x)\n", obj, contentLength, actualDigest)
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
	logger.Detailf("%v to %v\n", logging.FormatBytes(in.TotalBytes()), obj)
	return digest.Sum(nil), err
}
