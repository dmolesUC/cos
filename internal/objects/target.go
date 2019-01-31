package objects

import (
	"fmt"
	"net/url"
)

const (
	protocolSwift = "swift"
	protocolS3    = "s3"
)

// Target encapsulates a service URL and a bucket or container
type Target interface {
	Object(key string) Object
	Pretty() string
}

func NewTarget(endpointURL *url.URL, bucketURL *url.URL, region string) (Target, error) {
	protocol := bucketURL.Scheme
	bucket := bucketURL.Host

	if protocol == protocolSwift {
		return NewSwiftEndpoint(endpointURL, bucket)
	} else if protocol == protocolS3 {
		return NewS3Target(region, endpointURL, bucket), nil
	}
	return nil, fmt.Errorf("unsupported protocol: %#v", protocol)
}
