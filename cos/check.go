package cos

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

	. "github.com/dmolesUC3/cos/util"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Check struct {
	Logger    Logger
	ObjLoc    ObjectLocation
	Expected  []byte
	Algorithm string
	Region    string
}

func (c Check) GetDigest() ([]byte, error) {
	c.Logger.Detail("Initializing session")
	sess, err := InitSession(c.endpointP(), c.regionStrP(), c.Logger.Verbose())
	if err != nil {
		return nil, err
	}

	// TODO: don't write to tempfile
	outfile, err := ioutil.TempFile("", c.objFilename())
	if err != nil {
		return nil, err
	}
	defer func() {
		err := os.Remove(outfile.Name())
		if err != nil {
			c.Logger.Info(err)
		}
	}()
	c.Logger.Detailf("Downloading to tempfile: %v\n", outfile.Name())
	downloader := s3manager.NewDownloader(sess)
	bytesDownloaded, err := downloader.Download(outfile, &s3.GetObjectInput{
		Bucket: c.bucketP(),
		Key:    c.keyP(),
	})
	c.Logger.Detailf("Downloaded %d bytes\n", bytesDownloaded)
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
	c.Logger.Detailf("Hashed %d bytes\n", bytesHashed)
	if err != nil {
		return nil, err
	}
	digest := h.Sum(nil)

	if len(c.Expected) > 0 {
		if !bytes.Equal(c.Expected, digest) {
			err = fmt.Errorf("digest mismatch: expected: %x, actual: %x", c.Expected, digest)
		}
	}

	return digest, err
}

func (c Check) newHash() hash.Hash {
	if c.Algorithm == "sha256" {
		return sha256.New()
	}
	return md5.New()
}

func (c Check) regionStrP() *string {
	if c.Region != "" {
		c.Logger.Detailf("Using specified AWS region: %v\n", c.Region)
		return &c.Region
	}
	endpoint := c.endpointStr()
	regionStr := ExtractRegion(endpoint, c.Logger)
	return &regionStr
}

func (c Check) endpointStr() string {
	return c.ObjLoc.Endpoint.String()
}

func (c Check) endpointP() *string {
	endpointUrlStr := c.endpointStr()
	return &endpointUrlStr
}

func (c Check) objFilename() string {
	return path.Base(c.ObjLoc.Key())
}

func (c Check) bucketP() *string {
	bucket := c.ObjLoc.Bucket()
	return &bucket
}

func (c Check) keyP() *string {
	key := c.ObjLoc.Key()
	return &key
}
