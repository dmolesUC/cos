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

type Crvd struct {
	Object objects.Object
}

func (c *Crvd) CreateRetrieveValidate(body io.Reader, contentLength int64) (digest []byte, err error) {
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
			obj.Logger().Infof("validating %v (expected digest: %x)\n", objects.ProtocolUriStr(obj), digest)
			check := Check{Object: obj, Expected: digest, Algorithm: "sha256"}
			return check.CalcDigest()
		}
	}
	return digest, err
}

func (c *Crvd) CreateRetrieveValidateDelete(body io.Reader, contentLength int64) (digest []byte, err error) {
	digest, err = c.CreateRetrieveValidate(body, contentLength)
	if err == nil {
		err = c.Object.Delete()
	}
	return digest, err
}

func (c *Crvd) create(body io.Reader, contentLength int64) ([] byte, error) {
	// TODO: extract most of this to a struct & clean up
	obj := c.Object
	logger := obj.Logger()

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
		logger.Infof("uploading %d bytes to %v\n", contentLength, objects.ProtocolUriStr(obj))
		lr := io.LimitReader(tr, contentLength)
		errs <- obj.StreamUp(lr)
		logger.Detail("upload goroutine complete")
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
			bytesRead, err2 := streamer.StreamDown(logger, func(bytes []byte) error {
				_, err := hash.Write(bytes)
				return err
			})
			if bytesRead != contentLength {
				logger.Detailf("expected %d bytes, got %d\n", contentLength, bytesRead)
			}
			err = err2
		}
		errs <- err
		logger.Detail("hashing goroutine complete")
	}()

	logger.Detail("waiting for goroutines to complete")
	streamWg.Wait()
	close(errs)

	var allErrs []string
	for err := range errs {
		if err != nil {
			allErrs = append(allErrs, err.Error())
		}
	}
	if len(allErrs) > 0 {
		return nil, fmt.Errorf("error(s) creating object: %v", strings.Join(allErrs, ", "))
	}

	digest := hash.Sum(nil)
	return digest, nil
}
