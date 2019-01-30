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
func (b ObjectBuilder) Build() (Object, error) {
	builder, err := b.parsingObjURLStr()
	if err != nil {
		return nil, err
	}
	builder, err = builder.parsingEndpointStr()
	if err != nil {
		return nil, err
	}
	if err = builder.checkRequiredFields(); err != nil {
		return nil, err
	}
	if builder.protocol == protocolS3 {
		builder.region = protocols.EnsureS3Region(builder.region, builder.endpoint)
		if builder.region == "" {
			return nil, fmt.Errorf("unable to determine region for S3 object")
		}
		return &S3Object{
			region:   builder.region,
			endpoint: builder.endpoint,
			bucket:   builder.bucket,
			key:      builder.key,
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

	objName := builder.key
	if strings.HasPrefix(objName, "/") {
		objName = objName[1:]
	}
	return &SwiftObject{
		container:  builder.bucket,
		objectName: objName,
		cnxParams: protocols.SwiftConnectionParams{
			UserName: swiftAPIUser,
			APIKey:   swiftAPIKey,
			AuthURL:  builder.endpoint,
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

func (b ObjectBuilder) WithProtocolUri(protocolURI *url.URL) ObjectBuilder {
	b.protocol = protocolURI.Scheme
	b.bucket = protocolURI.Host
	b.key = protocolURI.Path
	logging.DefaultLogger().Tracef("Parsed protocol '%v', bucket '%v' and key '%v' from URL '%v'\n", b.protocol, b.bucket, b.key, protocolURI)
	return b
}

func (b ObjectBuilder) parsingObjURLStr() (ObjectBuilder, error) {
	if b.objURLStr == "" {
		return b, nil
	}
	objURL, err := ValidAbsURL(b.objURLStr)
	if err != nil {
		return b, err
	}

	scheme := objURL.Scheme
	if scheme == protocolSwift || scheme == protocolS3 {
		return b.WithProtocolUri(objURL), nil
	}

	var s3Uri *url.URL
	if scheme == "http" || scheme == "https" {
		logger := logging.DefaultLogger()
		endpointStr := fmt.Sprintf("%v://%v/", scheme, objURL.Host)
		logger.Tracef("Extracted endpoint URL '%v' from object URL '%v'", endpointStr, objURL)
		b.endpoint, err = url.Parse(endpointStr)
		if err != nil {
			return b, err
		}
		s3Uri, err = url.Parse("s3:/" + objURL.Path)
		logger.Tracef("Constructed S3 URI '%v' from object URL '%v'", s3Uri, objURL)
		if err != nil {
			return b, err
		}
		return b.WithProtocolUri(s3Uri), nil
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
