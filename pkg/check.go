package pkg

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/dmolesUC3/cos/internal"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Chunk size for ranged downloads
const DefaultChunkSize = int64(1024 * 1024 * 5)

// The Check struct represents a fixity check operation
type Check struct {
	Logger    internal.Logger
	ObjLoc    internal.ObjectLocation
	Expected  []byte
	Algorithm string
	Region    string
}

// GetDigest gets the digest, returning an error if the object cannot be retrieved or,
// when an expected digest is provided, if the calculated digest does not match.
func (c Check) GetDigest() ([]byte, error) {
	//digest, err := c.hashFromTempFile()
	digest, err := c.hashFromRanges()
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

func (c Check) hashFromRanges() ([] byte, error) {
	ol := c.ObjLoc
	logger := c.Logger
	goOutput, err := ol.GetObject()
	if err != nil {
		return nil, err
	}
	contentLength := *goOutput.ContentLength
	logger.Detailf("Expected ContentLength: %d\n", contentLength)

	acceptRanges := goOutput.AcceptRanges
	if acceptRanges == nil || *acceptRanges != "bytes" {
		var actual string
		if acceptRanges == nil {
			actual = "<nil>"
		} else {
			actual = *acceptRanges
		}
		return nil, fmt.Errorf("range request not supported; expected accept-ranges: 'bytes' but was '%v'", actual)
	}

	awsSession, err := ol.Session()
	if err != nil {
		return nil, err
	}

	h := c.newHash()
	downloader := s3manager.NewDownloader(awsSession)

	// TODO: make this configurable
	chunkSize := DefaultChunkSize

	chunkCount := (contentLength + chunkSize - 1) / chunkSize
	for chunk := int64(0); chunk < chunkCount; chunk += 1 {
		// byte ranges are 0-indexed and inclusive
		startInclusive := chunk * chunkSize
		var endInclusive int64
		if chunk + 1 < chunkCount {
			endInclusive = startInclusive + chunkSize - 1
		} else {
			endInclusive = contentLength - 1
		}
		expectedBytes := (endInclusive + 1) - startInclusive
		logger.Detailf("chunk %d of %d: retrieving %d bytes (%d - %d)\n", chunk, chunkCount, expectedBytes, startInclusive, endInclusive)
		goInput := s3.GetObjectInput{
			Bucket: ol.Bucket(),
			Key: ol.Key(),
		}

		rangeStr := fmt.Sprintf("bytes=%d-%d", startInclusive, endInclusive)
		w := aws.NewWriteAtBuffer(make([]byte, expectedBytes))
		goInput.Range = &rangeStr
		bytesDownloaded, err := downloader.Download(w, &goInput)
		if err != nil {
			return nil, err
		}
		if bytesDownloaded != expectedBytes {
			logger.Infof("chunk %d of %d: expected %d bytes (%d - %d), got %d\n", chunk, chunkCount, expectedBytes, startInclusive, endInclusive, bytesDownloaded)
		}
		result := w.Bytes()
		h.Write(result)
	}
	digest := h.Sum(nil)
	return digest, nil
}

// Deprecated: use hashFromRanges
func (c Check) hashFromTempFile() ([]byte, error) {
	objLoc := c.ObjLoc
	logger := c.Logger

	// TODO: don't write to tempfile
	filename := path.Base(*objLoc.Key())
	outfile, err := ioutil.TempFile("", filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := os.Remove(outfile.Name())
		if err != nil {
			logger.Info(err)
		}
	}()
	logger.Detailf("Downloading to tempfile: %v\n", outfile.Name())

	bytesDownloaded, err := objLoc.DownloadTo(outfile)
	logger.Detailf("Downloaded %d bytes\n", bytesDownloaded)
	if err != nil {
		return nil, err
	}

	err = outfile.Close() // TODO is this necessary?
	if err != nil {
		return nil, err
	}

	infile, err := os.Open(outfile.Name())
	if err != nil {
		return nil, err
	}
	h := c.newHash()
	bytesHashed, err := io.Copy(h, infile)
	logger.Detailf("Hashed %d bytes\n", bytesHashed)
	if err != nil {
		return nil, err
	}
	digest := h.Sum(nil)
	return digest, nil
}

func (c Check) newHash() hash.Hash {
	if c.Algorithm == "sha256" {
		return sha256.New()
	}
	return md5.New()
}



