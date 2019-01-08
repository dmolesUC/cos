package internal

import (
	"fmt"
	"net/url"
)

// ------------------------------------------------------------
// Exported types

// The Object type represents an object in cloud storage.
type Object interface {
	Region() *string
	Endpoint() *url.URL
	Bucket() *string
	Key() *string
}

// An ObjectBuilder builds an Object
type ObjectBuilder struct {
	region *string
	endpoint *url.URL
	bucket *string
	key *string
	objURLStr *string
	endpointStr *string
}

// NewObjectBuilder Returns a new empty ObjectBuilder
func NewObjectBuilder() ObjectBuilder {
	return ObjectBuilder{}
}

// WithRegion sets the region, or clears it if the specified region is blank
func (b ObjectBuilder) WithRegion(region string) ObjectBuilder {
	if region == "" {
		b.region = nil
	} else {
		b.region = &region
	}
	return b
}

// WithEndpoint sets the endpoint as a URL
func (b ObjectBuilder) WithEndpoint(endpoint *url.URL) ObjectBuilder {
	b.endpoint = endpoint
	return b
}

// WithEndpointStr sets the endpoint as a string, or clears it if the
// specified endpoint is blank
func (b ObjectBuilder) WithEndpointStr(endpointStr string) ObjectBuilder {
	if endpointStr == "" {
		b.endpointStr = nil
	} else {
		b.endpointStr = &endpointStr
	}
	return b
}

// WithBucket sets the bucket, or clears it if the specified region is blank
func (b ObjectBuilder) WithBucket(bucket string) ObjectBuilder {
	if bucket == "" {
		b.bucket = nil
	} else {
		b.bucket = &bucket
	}
	return b
}

// WithKey sets the key, or clears it if the specified region is blank
func (b ObjectBuilder) WithKey(key string) ObjectBuilder {
	if key == "" {
		b.key = nil
	} else {
		b.key = &key
	}
	return b
}

// WithObjectURLStr sets the object URL as a string, or clears it if the
// specified object URL is blank
func (b ObjectBuilder) WithObjectURLStr(objURLStr string) ObjectBuilder {
	if objURLStr == "" {
		b.objURLStr = nil
	} else {
		b.objURLStr = &objURLStr
	}
	return b
}

// Build builds a new Object from the state of this ObjectBuilder
func (b ObjectBuilder) Build(logger Logger) (Object, error) {
	b, err := b.parsingObjURLStr(logger)
	if err == nil {
		b, err = b.parsingEndpointStr()
	}
	if err == nil {
		b = b.ensureRegion(logger)
	}
	// TODO: validate everything
	return object{b.region, b.endpoint, b.bucket, b.key}, err
}

func (b ObjectBuilder) withS3Uri(s3Uri *url.URL, logger Logger) (ObjectBuilder, error) {
	if s3Uri.Scheme != "s3" {
		return b, fmt.Errorf("not an S3 URL: %v", s3Uri)
	}
	b.bucket = &s3Uri.Host
	b.key = &s3Uri.Path
	logger.Detailf("Parsed bucket '%v' and key '%v' from S3 URL '%v'", b.bucket, b.key, s3Uri)
	return b, nil
}

func (b ObjectBuilder) parsingObjURLStr(logger Logger) (ObjectBuilder, error) {
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

func (b ObjectBuilder) parsingEndpointStr() (ObjectBuilder, error) {
	if b.endpointStr == nil {
		return b, nil
	}
	endpoint, err := ValidAbsURL(*b.endpointStr)
	b.endpoint = endpoint
	return b, err
}

func (b ObjectBuilder) ensureRegion(logger Logger) ObjectBuilder {
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
	region *string
	endpoint *url.URL
	bucket *string
	key *string
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