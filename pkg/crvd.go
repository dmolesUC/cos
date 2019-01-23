package pkg

import (
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/dmolesUC3/cos/internal/objects"
)

type Crvd struct {
	Object objects.Object
}

func (c *Crvd) CreateRetrieveVerify(body io.Reader, contentLength int64) (digest []byte, err error) {
	digest, err = c.create(body, contentLength)
	if err == nil {
		obj := c.Object
		obj.Logger().Detailf("calculated digest on upload: %x\n", digest)

		obj.Reset()
		var actualLength int64
		actualLength, err = obj.ContentLength()
		if err == nil {
			if actualLength != contentLength {
				return digest, fmt.Errorf("content-length mismatch: expected: %d, actual: %d", contentLength, actualLength)
			}
			obj.Logger().Detailf("uploaded %d bytes\n", contentLength)
			obj.Logger().Infof("verifying %v (expected digest: %x)\n", objects.ProtocolUriStr(obj), digest)
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
		obj.Reset()
		obj.Logger().Detailf("verified %v\n", objects.ProtocolUriStr(obj))
		err = obj.Delete()
	}
	return digest, err
}

func (c *Crvd) create(body io.Reader, contentLength int64) ([] byte, error) {
	hash := sha256.New()
	err := c.Object.StreamUp(io.TeeReader(body, hash), contentLength)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), err
}
