package internal

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ------------------------------------------------------------
// Exported types

// The ObjectLocation type represents the location of an object in cloud storage.
type ObjectLocation interface {
	Region() *string
	Endpoint() *url.URL
	Bucket() *string
	Key() *string
	Session() (*session.Session, error)
	GetObject() (*s3.GetObjectOutput, error)
	DownloadTo(w io.WriterAt) (int64, error)
}

// An ObjectLocationBuilder builds an ObjectLocation
type ObjectLocationBuilder struct {
	region         string
	endpoint       *url.URL
	bucket         string
	key            string
	objURLStr      string
	endpointStr    string
}

// NewObjectLocationBuilder Returns a new empty ObjectLocationBuilder
func NewObjectLocationBuilder() ObjectLocationBuilder {
	return ObjectLocationBuilder{}
}

// WithRegion sets the region, or clears it if the specified region is blank
func (b ObjectLocationBuilder) WithRegion(region string) ObjectLocationBuilder {
	b.region = region
	return b
}

// WithEndpoint sets the endpoint as a URL
func (b ObjectLocationBuilder) WithEndpoint(endpoint *url.URL) ObjectLocationBuilder {
	b.endpoint = endpoint
	return b
}

// WithEndpointStr sets the endpoint as a string, or clears it if the
// specified endpoint is blank
func (b ObjectLocationBuilder) WithEndpointStr(endpointStr string) ObjectLocationBuilder {
	b.endpointStr = endpointStr
	return b
}

// WithBucket sets the bucket, or clears it if the specified region is blank
func (b ObjectLocationBuilder) WithBucket(bucket string) ObjectLocationBuilder {
	b.bucket = bucket
	return b
}

// WithKey sets the key, or clears it if the specified region is blank
func (b ObjectLocationBuilder) WithKey(key string) ObjectLocationBuilder {
	b.key = key
	return b
}

// WithObjectURLStr sets the object URL as a string, or clears it if the
// specified object URL is blank
func (b ObjectLocationBuilder) WithObjectURLStr(objURLStr string) ObjectLocationBuilder {
	b.objURLStr = objURLStr
	return b
}

// Build builds a new ObjectLocation from the state of this ObjectLocationBuilder
func (b ObjectLocationBuilder) Build(logger Logger) (ObjectLocation, error) {
	builder, err := b.parsingObjURLStr(logger)
	if err != nil {
		return objLoc{}, err
	}
	builder, err = builder.parsingEndpointStr()
	if err != nil {
		return objLoc{}, err
	}
	builder = builder.ensureRegion(logger)
	if err = builder.validate(); err != nil {
		return objLoc{}, err
	}
	ol := objLoc{
		region:         builder.region,
		endpoint:       builder.endpoint,
		bucket:         builder.bucket,
		key:            builder.key,
		verboseLogging: logger.Verbose(),
	}
	return ol, nil
}

func (b ObjectLocationBuilder) validate() error {
	var missing []string
	if b.region == "" {
		missing = append(missing, "region")
	}
	if b.endpoint == nil {
		missing = append(missing, "endpoint")
	}
	if b.bucket == "" {
		missing = append(missing, "bucket")
	}
	if b.key == "" {
		missing = append(missing, "key")
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("missing fields: %v", strings.Join(missing, ", "))
}

func (b ObjectLocationBuilder) withS3Uri(s3Uri *url.URL, logger Logger) (ObjectLocationBuilder, error) {
	if s3Uri.Scheme != "s3" {
		return b, fmt.Errorf("not an S3 URL: %v", s3Uri)
	}
	b.bucket = s3Uri.Host
	b.key = s3Uri.Path
	logger.Detailf("Parsed bucket '%v' and key '%v' from S3 URL '%v'\n", b.bucket, b.key, s3Uri)
	return b, nil
}

func (b ObjectLocationBuilder) parsingObjURLStr(logger Logger) (ObjectLocationBuilder, error) {
	if b.objURLStr == "" {
		return b, nil
	}
	objURL, err := ValidAbsURL(b.objURLStr)
	if err != nil {
		return b, err
	}

	var s3Uri *url.URL
	scheme := objURL.Scheme
	if scheme == "http" || scheme == "https" {
		endpointStr := fmt.Sprintf("%v://%v/", scheme, objURL.Host)
		logger.Detailf("Extracted endpoint URL '%v' from object URL '%v'", endpointStr, objURL)
		b.endpoint, err = url.Parse(endpointStr)
		if err != nil {
			return b, err
		}
		s3Uri, err = url.Parse("s3:/" + objURL.Path)
		logger.Detailf("Constructed S3 URI '%v' from object URL '%v'", s3Uri, objURL)
		if err != nil {
			return b, err
		}
	} else if scheme == "s3" {
		s3Uri = objURL
	} else {
		return b, fmt.Errorf("object URL '%v' is not of form s3://<bucket>/<key> or http(s)://<endpoint>/<bucket>/<key>", objURL)
	}
	return b.withS3Uri(s3Uri, logger)
}

func (b ObjectLocationBuilder) parsingEndpointStr() (ObjectLocationBuilder, error) {
	if b.endpointStr == "" {
		return b, nil
	}
	endpoint, err := ValidAbsURL(b.endpointStr)
	b.endpoint = endpoint
	return b, err
}

func (b ObjectLocationBuilder) ensureRegion(logger Logger) ObjectLocationBuilder {
	if b.region == "" {
		region, err := RegionFromEndpoint(b.endpoint)
		if err == nil {
			logger.Detailf("Found AWS region in endpoint URL %v: %v\n", b.endpoint, *region)
			b.region = *region
		} else {
			logger.Detailf("No AWS region found in endpoint URL '%v' (%v); using default region %v\n", b.endpoint, err, DefaultAwsRegion)
			regionStr := DefaultAwsRegion
			b.region = regionStr
		}
	} else {
		logger.Detailf("Using specified AWS region: %v\n", b.region)
	}
	return b
}

// ------------------------------------------------------------
// Unexported implementation

type objLoc struct {
	region         string
	endpoint       *url.URL
	bucket         string
	key            string
	verboseLogging bool
	awsSession     *session.Session
}

func (ol objLoc) Region() *string {
	if ol.region == "" {
		return nil
	}
	return &ol.region
}

func (ol objLoc) Endpoint() *url.URL {
	return ol.endpoint
}

func (ol objLoc) Bucket() *string {
	if ol.bucket == "" {
		return nil
	}
	return &ol.bucket
}

func (ol objLoc) Key() *string {
	if ol.key == "" {
		return nil
	}
	return &ol.key
}

func (ol objLoc) Session() (*session.Session, error) {
	var err error
	if ol.awsSession == nil {
		endpointStr := ol.endpoint.String()
		ol.awsSession, err = InitSession(&endpointStr, ol.Region(), ol.verboseLogging)
	}
	return ol.awsSession, err
}

func (ol objLoc) GetObject() (*s3.GetObjectOutput, error) {
	awsSession, err := ol.Session()
	if err != nil {
		return nil, err
	}
	s3Svc := s3.New(awsSession)
	return s3Svc.GetObject(ol.toGetObjectInput())
}

func (ol objLoc) DownloadTo(w io.WriterAt) (int64, error) {
	awsSession, err := ol.Session()
	if err != nil {
		return 0, err
	}
	downloader := s3manager.NewDownloader(awsSession)
	return downloader.Download(w, ol.toGetObjectInput())
}

// ------------------------------------------------------------
// Unexported functions

func (ol objLoc) toGetObjectInput() *s3.GetObjectInput {
	goInput := s3.GetObjectInput{
		Bucket: ol.Bucket(),
		Key: ol.Key(),
	}
	return &goInput
}
