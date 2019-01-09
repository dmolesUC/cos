package internal

import (
	"fmt"
	"net/url"
	"strings"
)

// ------------------------------------------------------------
// Exported types

// The ObjectLocation type represents the location of an object in cloud storage.
// TODO: don't use string pointers just b/c Amazon likes them
type ObjectLocation interface {
	Region() *string
	Endpoint() *url.URL
	Bucket() *string
	Key() *string
}

func EndpointP(o ObjectLocation) *string {
	endpointStr := o.Endpoint().String()
	return &endpointStr
}

// An ObjectLocationBuilder builds an ObjectLocation
type ObjectLocationBuilder struct {
	region      *string
	endpoint    *url.URL
	bucket      *string
	key         *string
	objURLStr   *string
	endpointStr *string
}

// NewObjectLocationBuilder Returns a new empty ObjectLocationBuilder
func NewObjectLocationBuilder() ObjectLocationBuilder {
	return ObjectLocationBuilder{}
}

// WithRegion sets the region, or clears it if the specified region is blank
func (b ObjectLocationBuilder) WithRegion(region string) ObjectLocationBuilder {
	if region == "" {
		b.region = nil
	} else {
		b.region = &region
	}
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
	if endpointStr == "" {
		b.endpointStr = nil
	} else {
		b.endpointStr = &endpointStr
	}
	return b
}

// WithBucket sets the bucket, or clears it if the specified region is blank
func (b ObjectLocationBuilder) WithBucket(bucket string) ObjectLocationBuilder {
	if bucket == "" {
		b.bucket = nil
	} else {
		b.bucket = &bucket
	}
	return b
}

// WithKey sets the key, or clears it if the specified region is blank
func (b ObjectLocationBuilder) WithKey(key string) ObjectLocationBuilder {
	if key == "" {
		b.key = nil
	} else {
		b.key = &key
	}
	return b
}

// WithObjectURLStr sets the object URL as a string, or clears it if the
// specified object URL is blank
func (b ObjectLocationBuilder) WithObjectURLStr(objURLStr string) ObjectLocationBuilder {
	if objURLStr == "" {
		b.objURLStr = nil
	} else {
		b.objURLStr = &objURLStr
	}
	return b
}

// Build builds a new ObjectLocation from the state of this ObjectLocationBuilder
func (b ObjectLocationBuilder) Build(logger Logger) (ObjectLocation, error) {
	b, err := b.parsingObjURLStr(logger)
	if err != nil {
		return object{}, err
	}
	b, err = b.parsingEndpointStr()
	if err != nil {
		return object{}, err
	}
	b = b.ensureRegion(logger)
	if err = b.validate(); err != nil {
		return object{}, err
	}
	return object{b.region, b.endpoint, b.bucket, b.key}, nil
}

func (b ObjectLocationBuilder) validate() error {
	var missing []string
	if b.region == nil {
		missing = append(missing, "region")
	}
	if b.endpoint == nil {
		missing = append(missing, "endpoint")
	}
	if b.bucket == nil {
		missing = append(missing, "bucket")
	}
	if b.key == nil {
		missing = append(missing, "key")
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("missing object fields: %v", strings.Join(missing, ", "))
}

func (b ObjectLocationBuilder) withS3Uri(s3Uri *url.URL, logger Logger) (ObjectLocationBuilder, error) {
	if s3Uri.Scheme != "s3" {
		return b, fmt.Errorf("not an S3 URL: %v", s3Uri)
	}
	b.bucket = &s3Uri.Host
	b.key = &s3Uri.Path
	logger.Detailf("Parsed bucket '%v' and key '%v' from S3 URL '%v'\n", *b.bucket, *b.key, s3Uri)
	return b, nil
}

func (b ObjectLocationBuilder) parsingObjURLStr(logger Logger) (ObjectLocationBuilder, error) {
	if b.objURLStr == nil {
		return b, nil
	}
	objURL, err := ValidAbsURL(*b.objURLStr)
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
	if b.endpointStr == nil {
		return b, nil
	}
	endpoint, err := ValidAbsURL(*b.endpointStr)
	b.endpoint = endpoint
	return b, err
}

func (b ObjectLocationBuilder) ensureRegion(logger Logger) ObjectLocationBuilder {
	if b.region == nil {
		region, err := RegionFromEndpoint(b.endpoint)
		if err == nil {
			logger.Detailf("Found AWS region in endpoint URL %v: %v\n", b.endpoint, region)
			b.region = region
		} else {
			logger.Detailf("No AWS region found in endpoint URL '%v'; using default region %v\n", b.endpoint, DefaultAwsRegion)
			regionStr := DefaultAwsRegion
			b.region = &regionStr
		}
	} else {
		logger.Detailf("Using specified AWS region: %v\n", b.region)
	}
	return b
}

// ------------------------------------------------------------
// Unexported implementation

type object struct {
	region   *string
	endpoint *url.URL
	bucket   *string
	key      *string
}

func (b object) Region() *string {
	return b.region
}

func (b object) Endpoint() *url.URL {
	return b.endpoint
}

func (b object) Bucket() *string {
	return b.bucket
}

func (b object) Key() *string {
	return b.key
}

