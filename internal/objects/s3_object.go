package objects

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/protocols"
)

// S3Object is an S3 implementation of Object
type S3Object struct {
	region     string
	endpoint   *url.URL
	bucket     string
	key        string
	awsSession *session.Session
	head       *s3.HeadObjectOutput
}

func (obj *S3Object) Protocol() string {
	return protocolS3
}

func (obj *S3Object) Pretty() string {
	format := `S3Object { 
				region:  '%v' 
				endpoint: %v 
				bucket:  '%v' 
		        key :    '%v'
				session: '%v'
			}`
	format = logging.Untabify(format, " ")
	args := logging.Prettify(obj.region, obj.formatEndpoint(), obj.bucket, obj.key, obj.formatSession())
	return fmt.Sprintf(format, args...)
}

func (obj *S3Object) String() string {
	return fmt.Sprintf(
		"{region: %v, endpoint: %v, bucket: %v, key: %v, awsSession: %v}",
		obj.region, obj.formatEndpoint(), obj.bucket, obj.key, obj.formatSession(),
	)
}

func (obj *S3Object) Refresh() {
	obj.awsSession = nil
	obj.head = nil
}

// Endpoint returns the endpoint URL used to access the object
func (obj *S3Object) Endpoint() *url.URL {
	return obj.endpoint
}

// Bucket returns the object's storage bucket
func (obj *S3Object) Bucket() *string {
	if obj.bucket == "" {
		return nil
	}
	return &obj.bucket
}

// Key returns the object's storage key
func (obj *S3Object) Key() *string {
	if obj.key == "" {
		return nil
	}
	return &obj.key
}

// ContentLength gets the size of the object in bytes, or returns an
// error if the size cannot be determined.
func (obj *S3Object) ContentLength() (int64, error) {
	goOutput, err := obj.Head()
	if err != nil {
		logging.DefaultLogger().Tracef("error determining content-length: %v\n", err)
		return 0, err
	}
	if goOutput == nil {
		return 0, errors.New("no output returned by GetObject")
	}
	contentLength := goOutput.ContentLength
	if contentLength == nil {
		return 0, errors.New("no content-length returned by GetObject")
	}
	return *contentLength, nil
}

// SupportsRanges returns true if the object supports ranged downloads,
// false otherwise
func (obj *S3Object) SupportsRanges() bool {
	goOutput, err := obj.Head()
	if err == nil {
		logger := logging.DefaultLogger()
		acceptRanges := goOutput.AcceptRanges
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

func (obj *S3Object) DownloadRange(startInclusive, endInclusive int64, buffer []byte) (int64, error) {
	if !obj.SupportsRanges() {
		logging.DefaultLogger().Tracef("object %v may not support ranged downloads; trying anyway\n", obj)
	}
	rangeStr := fmt.Sprintf("bytes=%d-%d", startInclusive, endInclusive)
	goInput := s3.GetObjectInput{
		Bucket: obj.Bucket(),
		Key:    obj.Key(),
		Range:  &rangeStr,
	}

	awsSession, err := obj.sessionP()
	if err != nil {
		return 0, err
	}
	downloader := s3manager.NewDownloader(awsSession)
	target := aws.NewWriteAtBuffer(buffer)
	return downloader.Download(target, &goInput)
}

func (obj *S3Object) Create(body io.Reader, length int64) (err error) {
	awsSession, err := obj.sessionP()
	if err != nil {
		return err
	}
	logging.DefaultLogger().Detailf("Uploading %d bytes to %v\n", length, ProtocolUriStr(obj))

	// TODO: allow object to include an expected MD5
	uploader := s3manager.NewUploader(awsSession)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: obj.Bucket(),
		Key:    obj.Key(),
		Body:   body,
	})
	if err == nil {
		logging.DefaultLogger().Detailf("Uploaded %d bytes to %v\n", length, result.Location)
	}
	return err
}

func (obj *S3Object) Delete() (err error) {
	protocolUriStr := ProtocolUriStr(obj)
	logging.DefaultLogger().Tracef("Delete: getting session for %v\n", protocolUriStr)
	awsSession, err := obj.sessionP()
	if err != nil {
		return err
	}
	doInput := s3.DeleteObjectInput{
		Bucket: obj.Bucket(),
		Key:    obj.Key(),
	}
	logging.DefaultLogger().Detailf("Deleting %v\n", protocolUriStr)
	_, err = s3.New(awsSession).DeleteObject(&doInput)
	if err == nil {
		logging.DefaultLogger().Detailf("Deleted %v\n", protocolUriStr)
	} else {
		logging.DefaultLogger().Detailf("Deleting %v failed: %v", protocolUriStr, err)
	}
	obj.Refresh()
	return err
}

// ------------------------------------------------------------
// Unexported functions

// Region returns the AWS region of the object
func (obj *S3Object) regionP() *string {
	if obj.region == "" {
		return nil
	}
	return &obj.region
}

// TODO: cache sessions above the object level?
func (obj *S3Object) sessionP() (*session.Session, error) {
	var err error
	if obj.awsSession == nil {
		endpointStr := obj.endpoint.String()
		obj.awsSession, err = protocols.ValidS3Session(&endpointStr, obj.regionP())
	}
	return obj.awsSession, err
}

func (obj *S3Object) toHeadObjectInput() *s3.HeadObjectInput {
	goInput := s3.HeadObjectInput{
		Bucket: obj.Bucket(),
		Key:    obj.Key(),
	}
	return &goInput
}

func (obj *S3Object) Head() (*s3.HeadObjectOutput, error) {
	var err error
	if obj.head == nil {
		awsSession, err := obj.sessionP()
		if err == nil {
			s3Svc := s3.New(awsSession)
			head, err := s3Svc.HeadObject(obj.toHeadObjectInput())
			if err != nil {
				return nil, err
			}
			if head == nil {
				return nil, fmt.Errorf("nil *HeadObjectOutput returned by S3.HeadObject")
			}
			obj.head = head

		}
	}
	return obj.head, err
}

func (obj *S3Object) formatSession() string {
	var sessionStr string
	if obj.awsSession == nil {
		sessionStr = "<nil>"
	} else {
		sessionStr = "(initialized)"
	}
	return sessionStr
}

func (obj *S3Object) formatEndpoint() string {
	var endpointStr string
	if obj.endpoint == nil {
		endpointStr = "<nil>"
	} else {
		endpointStr = obj.Endpoint().String()
	}
	return endpointStr
}
