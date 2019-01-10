package objects

import (
	"fmt"
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
	logger     logging.Logger
	awsSession *session.Session
	goOutput   *s3.GetObjectOutput
}

func (obj *S3Object) String() string {
	var endpointStr string
	if obj.endpoint == nil {
		endpointStr = "<nil>"
	} else {
		endpointStr = obj.Endpoint().String()
	}
	var sessionStr string
	if obj.awsSession == nil {
		sessionStr = "<nil>"
	} else {
		sessionStr = "(initialized)"
	}
	return fmt.Sprintf(
		"{region: %v, endpoint: %v, bucket: %v, key: %v, logger: %v, awsSession: %v}",
		obj.region, endpointStr, obj.bucket, obj.key, obj.logger, sessionStr,
	)
}

// Region returns the AWS region of the object
func (obj *S3Object) Region() *string {
	if obj.region == "" {
		return nil
	}
	return &obj.region
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
	goOutput, err := obj.getObject()
	if err != nil {
		obj.logger.Detailf("error determining content-length: %v", err)
		return 0, err
	}
	return *goOutput.ContentLength, nil
}

// SupportsRanges returns true if the object supports ranged downloads,
// false otherwise
func (obj *S3Object) SupportsRanges() bool {
	goOutput, err := obj.getObject()
	if err == nil {
		acceptRanges := goOutput.AcceptRanges
		if acceptRanges != nil {
			actual := *acceptRanges
			if "bytes" == actual {
				return true
			}
			obj.logger.Detailf("range request not supported; expected accept-ranges: 'bytes' but was '%v'\n", actual)
		} else {
			obj.logger.Detail("range request not supported; expected accept-ranges: 'bytes' but was no accept-ranges header found")
		}
	}
	return false
}

// StreamDown streams the object down in ranged requests of the specified size, passing
// each byte range retrieved to the specified handler function, in sequence.
func (obj *S3Object) StreamDown(rangeSize int64, handleBytes func([]byte) error) (int64, error) {
	// TODO: parallel downloads, serial handling? cf. https://coderwall.com/p/uz2noa/fast-parallel-downloads-in-golang-with-accept-ranges-and-goroutines
	if !obj.SupportsRanges() {
		return 0, fmt.Errorf("object %v does not support ranged downloads", obj)
	}
	logger := obj.logger

	awsSession, err := obj.session()
	if err != nil {
		return 0, err
	}
	downloader := s3manager.NewDownloader(awsSession)

	contentLength, err := obj.ContentLength()
	if err != nil {
		return 0, err
	}

	totalBytes := int64(0)
	rangeCount := (contentLength + rangeSize - 1) / rangeSize
	for rangeIndex := int64(0); rangeIndex < rangeCount; rangeIndex++ {
		// byte ranges are 0-indexed and inclusive
		startInclusive := rangeIndex * rangeSize
		var endInclusive int64
		if (rangeIndex + 1) < rangeCount {
			endInclusive = startInclusive + rangeSize - 1
		} else {
			endInclusive = contentLength - 1
		}
		rangeStr := fmt.Sprintf("bytes=%d-%d", startInclusive, endInclusive)

		expectedBytes := (endInclusive + 1) - startInclusive
		logger.Detailf("range %d of %d: retrieving %d bytes (%d - %d)\n", rangeIndex, rangeCount, expectedBytes, startInclusive, endInclusive)
		goInput := s3.GetObjectInput{
			Bucket: obj.Bucket(),
			Key: obj.Key(),
			Range: &rangeStr,
		}

		target := aws.NewWriteAtBuffer(make([]byte, expectedBytes))
		actualBytes, err := downloader.Download(target, &goInput)
		totalBytes = totalBytes + actualBytes
		if err != nil {
			return totalBytes, err
		}
		if actualBytes != expectedBytes {
			logger.Infof("range %d of %d: expected %d bytes (%d - %d), got %d\n", rangeIndex, rangeCount, expectedBytes, startInclusive, endInclusive, actualBytes)
		}
		byteRange := target.Bytes()
		err = handleBytes(byteRange)
		if err != nil {
			return totalBytes, err
		}
	}
	return totalBytes, nil
}

// ------------------------------------------------------------
// Unexported functions

func (obj *S3Object) session() (*session.Session, error) {
	var err error
	if obj.awsSession == nil {
		endpointStr := obj.endpoint.String()
		verboseLogging := obj.logger.Verbose()
		// TODO: move this all back to s3_utils
		obj.awsSession, err = protocols.InitS3Session(&endpointStr, obj.Region(), verboseLogging)
		isEC2, err := protocols.IsEC2()
		if err != nil {
			obj.logger.Detailf("Rrror trying to determine whether we're running in EC2 (assume we're not): %v", err)
			isEC2 = false
		}
		if isEC2 {
			obj.logger.Detailf("Running in EC2; allowing IAM role credentials\n")
		} else {
			// TODO: https://github.com/aws/aws-sdk-go/issues/2392
			obj.logger.Detailf("Not running in EC2; disallowing IAM role credentials\n")
			return protocols.ValidateCredentials(obj.awsSession)
		}
	}
	return obj.awsSession, err
}

func (obj *S3Object) toGetObjectInput() *s3.GetObjectInput {
	goInput := s3.GetObjectInput{
		Bucket: obj.Bucket(),
		Key:    obj.Key(),
	}
	return &goInput
}

func (obj *S3Object) getObject() (*s3.GetObjectOutput, error) {
	var err error
	if obj.goOutput == nil {
		awsSession, err := obj.session()
		if err == nil {
			s3Svc := s3.New(awsSession)
			obj.goOutput, err = s3Svc.GetObject(obj.toGetObjectInput())
		}
	}
	return obj.goOutput, err
}

