package pkg

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/dmolesUC3/cos/internal/objects"
	"github.com/dmolesUC3/cos/internal/streaming"
)

type Crd struct {
	Object objects.Object
}

func (c *Crd) CreateRetrieve(body io.Reader, contentLength int64) (digest []byte, err error) {
	digest, err = c.create(body, contentLength)
	if err == nil {
		check := Check{Object: c.Object, Expected: digest, Algorithm: "sha-256"}
		return check.CalcDigest()
	}
	return digest, err
}

func (c *Crd) create(body io.Reader, contentLength int64) ([] byte, error) {
	// TODO: extract most of this to a struct & clean up
	obj := c.Object

	errs := make(chan error, 3)

	pr, pw := io.Pipe()
	tr := io.TeeReader(body, pw)

	var streamWg sync.WaitGroup
	streamWg.Add(2)

	go func() {
		defer streamWg.Done()
		defer func() {
			errs <- pw.Close()
		}()
		errs <- obj.StreamUp(tr)
	}()

	hash := sha256.New()
	go func() {
		defer streamWg.Done()
		fillRange := func(byteRange *streaming.ByteRange) (int64, error) {
			bytesRead, err := io.ReadFull(pr, byteRange.Buffer)
			return int64(bytesRead), err
		}
		streamer, err := streaming.NewStreamer(streaming.DefaultRangeSize, contentLength, &fillRange)
		if err == nil {
			logger := obj.Logger()
			bytesRead, err2 := streamer.StreamDown(logger, func(bytes []byte) error {
				_, err := hash.Write(bytes)
				return err
			})
			if bytesRead != contentLength {
				logger.Detailf("expected %d bytes, got %d", contentLength, bytesRead)
			}
			err = err2
		}
		if err != nil {
			errs <- err
		}
	}()
	streamWg.Wait()
	close(errs)

	var allErrs []string
	for err := range errs {
		allErrs = append(allErrs, err.Error())
	}
	if len(allErrs) > 0 {
		return nil, fmt.Errorf("error(s) creating object: %v", strings.Join(allErrs, ", "))
	}

	digest := hash.Sum(nil)
	return digest, nil
}
