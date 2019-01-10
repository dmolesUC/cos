package objects

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/protocols"
)

const protocolSwift = "swift"
const protocolS3 = "s3"

// An ObjectBuilder builds an Object
type ObjectBuilder struct {
	region      string
	protocol    string
	endpoint    *url.URL
	bucket      string
	key         string
	objURLStr   string
	endpointStr string
}

// NewObjectBuilder Returns a new empty ObjectBuilder
func NewObjectBuilder() ObjectBuilder {
	return ObjectBuilder{}
}

// WithRegion sets the region, or clears it if the specified region is blank
func (b ObjectBuilder) WithRegion(region string) ObjectBuilder {
	b.region = region
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
	b.endpointStr = endpointStr
	return b
}

// WithBucket sets the bucket, or clears it if the specified region is blank
func (b ObjectBuilder) WithBucket(bucket string) ObjectBuilder {
	b.bucket = bucket
	return b
}

// WithKey sets the key, or clears it if the specified region is blank
func (b ObjectBuilder) WithKey(key string) ObjectBuilder {
	b.key = key
	return b
}

// WithObjectURLStr sets the object URL as a string, or clears it if the
// specified object URL is blank
func (b ObjectBuilder) WithObjectURLStr(objURLStr string) ObjectBuilder {
	b.objURLStr = objURLStr
	return b
}

// Build builds a new Object from the state of this ObjectBuilder
func (b ObjectBuilder) Build(logger logging.Logger) (Object, error) {
	builder, err := b.parsingObjURLStr(logger)
	if err != nil {
		return nil, err
	}
	builder, err = builder.parsingEndpointStr()
	if err != nil {
		return nil, err
	}
	builder = builder.ensureRegion(logger)
	if err = builder.checkRequiredFields(); err != nil {
		return nil, err
	}
	if b.protocol == protocolS3 {
		if b.region == "" {
			return nil, fmt.Errorf("unable to determine region for S3 object")
		}
		return &S3Object{
			region:   builder.region,
			endpoint: builder.endpoint,
			bucket:   builder.bucket,
			key:      builder.key,
			logger:   logger,
		}, nil
	}

	swiftAPIUser := os.Getenv("SWIFT_API_USER")
	if swiftAPIUser == "" {
		return nil, fmt.Errorf("missing environment variable $SWIFT_API_USER")
	}
	swiftAPIKey := os.Getenv("SWIFT_API_KEY")
	if swiftAPIKey == "" {
		return nil, fmt.Errorf("missing environment variable $SWIFT_API_KEY")
	}

	return &SwiftObject{
		container:  builder.bucket,
		objectName: builder.key,
		logger:     logger,
		cnxParams: protocols.SwiftConnectionParams{
			UserName: swiftAPIUser,
			APIKey:   swiftAPIKey,
			AuthURL:  b.endpoint,
		},
	}, nil
}

func (b ObjectBuilder) checkRequiredFields() error {
	var missing []string
	if b.protocol == "" {
		missing = append(missing, "protocol")
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

func (b ObjectBuilder) withProtocolURI(protocolURI *url.URL, logger logging.Logger) (ObjectBuilder) {
	b.protocol = protocolURI.Scheme
	b.bucket = protocolURI.Host
	b.key = protocolURI.Path
	logger.Detailf("Parsed protocol '%v', bucket '%v' and key '%v' from URL '%v'\n", b.protocol, b.bucket, b.key, protocolURI)
	return b
}

func (b ObjectBuilder) parsingObjURLStr(logger logging.Logger) (ObjectBuilder, error) {
	if b.objURLStr == "" {
		return b, nil
	}
	objURL, err := ValidAbsURL(b.objURLStr)
	if err != nil {
		return b, err
	}

	scheme := objURL.Scheme
	if scheme == protocolSwift || scheme == protocolS3 {
		return b.withProtocolURI(objURL, logger), nil
	}

	var s3Uri *url.URL
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
		return b.withProtocolURI(s3Uri, logger), nil
	}
	return b, fmt.Errorf("object URL '%v' is not of form s3://<bucket>/<key>, swift://<container>//<name>, or http(s)://<S3 endpoint>/<bucket>/<key>", objURL)
}

func (b ObjectBuilder) parsingEndpointStr() (ObjectBuilder, error) {
	if b.endpointStr == "" {
		return b, nil
	}
	endpoint, err := ValidAbsURL(b.endpointStr)
	b.endpoint = endpoint
	return b, err
}

func (b ObjectBuilder) ensureRegion(logger logging.Logger) ObjectBuilder {
	// TODO: can some of this move into s3_utils?
	if b.region == "" {
		region, err := protocols.RegionFromEndpoint(b.endpoint)
		if err == nil {
			logger.Detailf("Found AWS region in endpoint URL %v: %v\n", b.endpoint, *region)
			b.region = *region
		} else {
			logger.Detailf("No AWS region found in endpoint URL '%v' (%v); using default region %v\n", b.endpoint, err, protocols.DefaultAwsRegion)
			regionStr := protocols.DefaultAwsRegion
			b.region = regionStr
		}
	} else {
		logger.Detailf("Using specified AWS region: %v\n", b.region)
	}
	return b
}
