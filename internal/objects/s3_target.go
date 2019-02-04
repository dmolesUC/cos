package objects

import (
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ------------------------------------------------------------
// S3Target type

type S3Target struct {
	Region   string
	Endpoint string
	Bucket   string

	awsSession *session.Session
	s3Svc      *s3.S3
}

func NewS3Target(region string, endpointURL *url.URL, bucket string) *S3Target {
	return &S3Target{
		Region:   EnsureS3Region(region, endpointURL),
		Endpoint: endpointURL.String(),
		Bucket:   bucket,
	}
}

// ------------------------------
// Target implementation

func (e *S3Target) Object(key string) Object {
	return &S3Object{Endpoint: e, Key: key}
}

func (e *S3Target) Pretty() string {
	return fmt.Sprintf("S3Target{ Region: %#v, Endpoint: %#v, Bucket: %#v }", e.Region, e.Endpoint, e.Bucket)
}

func (e *S3Target) String() string {
	return e.Pretty()
}

// ------------------------------
// Miscellaneous methods

func (e *S3Target) Session() (*session.Session, error) {
	if e.awsSession == nil {
		awsSession, err := ValidS3Session(&e.Endpoint, &e.Region)
		if err != nil {
			return nil, err
		}
		e.awsSession = awsSession
	}
	return e.awsSession, nil
}

func (e *S3Target) S3() (*s3.S3, error) {
	if e.s3Svc == nil {
		awsSession, err := e.Session()
		if err != nil {
			return nil, err
		}
		e.s3Svc = s3.New(awsSession)
	}
	return e.s3Svc, nil
}
