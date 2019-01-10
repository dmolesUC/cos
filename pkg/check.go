package pkg

import (
	"bytes"
	"fmt"

	"github.com/dmolesUC3/cos/internal"
)

// DefaultRangeSize is the default range size for ranged downloads
const DefaultRangeSize = int64(1024 * 1024 * 5)

// The Check struct represents a fixity check operation
type Check struct {
	Logger    internal.Logger
	Object    internal.Object
	Expected  []byte
	Algorithm string
	Region    string
}

// CalcDigest gets the digest, returning an error if the object cannot be retrieved or,
// when an expected digest is provided, if the calculated digest does not match.
func (c Check) CalcDigest() ([]byte, error) {
	digest, err := c.Object.CalcDigest(DefaultRangeSize, c.Algorithm)
	if err != nil {
		return nil, err
	}
	if len(c.Expected) > 0 {
		if !bytes.Equal(c.Expected, digest) {
			err = fmt.Errorf("digest mismatch: expected: %x, actual: %x", c.Expected, digest)
		}
	}
	return digest, err
}

