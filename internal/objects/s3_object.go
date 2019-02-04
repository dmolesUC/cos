package objects

import (
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/dmolesUC3/cos/internal/logging"
)

// ------------------------------------------------------------
// S3Object type

type S3Object struct {
	Endpoint *S3Target
	Bucket   string
	Key      string
}

// ------------------------------
// Object implementation

func (obj *S3Object) Pretty() string {
	return fmt.Sprintf("s3://%v/%v", obj.Bucket, obj.Key)
}

func (obj *S3Object) String() string {
	return obj.Pretty()
}

func (obj *S3Object) GetEndpoint() Target {
	return obj.Endpoint
}

func (obj *S3Object) ContentLength() (length int64, err error) {
	h, err := obj.Head()
	if err != nil {
		return 0, err
	}
	lengthP := h.ContentLength
	if lengthP == nil {
		logger := logging.DefaultLogger()
		logger.Tracef("s3.HeadObject() returned nil content-length; trying GetObject()\n")
		o, err := obj.Get()
		if o != nil {
			defer func() {
				if o.Body == nil {
					return
				}
				if err := o.Body.Close(); err != nil {
					logger.Tracef("error closing object body: %v\n", err.Error())
				}
			}()
			lengthP = o.ContentLength
			if lengthP == nil {
				return 0, fmt.Errorf("s3.GetObject() returned nil content-length")
			}
		}
		if err != nil {
			return 0, err
		}
	}
	return *lengthP, nil
}

// SupportsRanges returns true if the object supports ranged downloads,
// false otherwise
func (obj *S3Object) SupportsRanges() bool {
	h, err := obj.Head()
	if err == nil {
		logger := logging.DefaultLogger()
		acceptRanges := h.AcceptRanges
		if acceptRanges != nil {
			actual := *acceptRanges
			if "bytes" == actual {
				return true
			}
			logger.Tracef("range request not supported; expected accept-ranges: 'bytes' but was '%v'\n", actual)
		} else {
			logger.Trace("range request not supported; expected accept-ranges: 'bytes' but was no accept-ranges header found")
		}
	}
	return false
}

func (obj *S3Object) DownloadRange(startInclusive, endInclusive int64, buffer []byte) (n int64, err error) {
	if !obj.SupportsRanges() {
		logging.DefaultLogger().Tracef("object %v may not support ranged downloads; trying anyway\n", obj)
	}

	awsSession, err := obj.Endpoint.Session()
	if err != nil {
		return 0, err
	}

	out := aws.NewWriteAtBuffer(buffer)
	rangeStr := fmt.Sprintf("bytes=%d-%d", startInclusive, endInclusive)
	downloader := s3manager.NewDownloader(awsSession)
	return downloader.Download(out, &s3.GetObjectInput{
		Bucket: &obj.Bucket,
		Key:    &obj.Key,
		Range:  &rangeStr,
	})
}

func (obj *S3Object) Create(body io.Reader, length int64) (err error) {
	awsSession, err := obj.Endpoint.Session()
	if err != nil {
		return err
	}
	logging.DefaultLogger().Detailf("Uploading %d bytes to %v\n", length, obj)

	uploader := s3manager.NewUploader(awsSession)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &obj.Bucket,
		Key:    &obj.Key,
		Body:   body,
	})
	if err == nil {
		logging.DefaultLogger().Detailf("Uploaded %d bytes to %v\n", length, result.Location)
	}
	return err
}

func (obj *S3Object) Delete() (err error) {
	protocolUriStr := obj
	awsSession, err := obj.Endpoint.Session()
	if err != nil {
		return err
	}
	logger := logging.DefaultLogger()
	logger.Detailf("Deleting %v\n", protocolUriStr)
	_, err = s3.New(awsSession).DeleteObject(&s3.DeleteObjectInput{
		Bucket: &obj.Bucket,
		Key:    &obj.Key,
	})
	if err == nil {
		logger.Detailf("Deleted %v\n", protocolUriStr)
	} else {
		logger.Detailf("Deleting %v failed: %v", protocolUriStr, err)
	}
	return err
}

// ------------------------------
// Miscellaneous methods

func (obj *S3Object) Head() (h *s3.HeadObjectOutput, err error) {
	s3Svc, err := obj.Endpoint.S3()
	if err != nil {
		return nil, err
	}

	h, err = s3Svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &obj.Bucket,
		Key:    &obj.Key,
	})
	if h != nil {
		return h, nil
	} else {
		return nil, errors.New("s3.HeadObject() returned nil")
	}
}

func (obj *S3Object) Get() (h *s3.GetObjectOutput, err error) {
	s3Svc, err := obj.Endpoint.S3()
	if err != nil {
		return nil, err
	}

	h, err = s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: &obj.Bucket,
		Key:    &obj.Key,
	})
	if h != nil {
		return h, nil
	} else {
		return nil, errors.New("s3.GetObject() returned nil")
	}
}

