package internal

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"unsafe"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	DefaultAwsRegion     = "us-west-2"
	DefaultS3EndpointURL = "https://s3-us-west-2.amazonaws.com"
	awsRegionRegexpStr = "https?://s3-([^.]+)\\.amazonaws\\.com"
)
var awsRegionRegexp = regexp.MustCompile(awsRegionRegexpStr)

func RegionFromEndpoint(endpoint *url.URL) (*string, error) {
	if endpoint == nil {
		return nil, fmt.Errorf("can't extract region from nil endpoint")
	}
	matches := awsRegionRegexp.FindStringSubmatch(endpoint.String())
	if len(matches) == 2 {
		regionStr := matches[1]
		return &regionStr, nil
	}
	return nil, fmt.Errorf("no AWS region found in endpoint URL %v", endpoint)
}

func InitSession(endpointP *string, regionStrP *string, verbose bool) (*session.Session, error) {
	s3Config := aws.Config{
		Endpoint:                      endpointP,
		Region:                        regionStrP,
		S3ForcePathStyle:              aws.Bool(true),
		CredentialsChainVerboseErrors: aws.Bool(verbose),
	}
	s3Opts := session.Options{
		Config:            s3Config,
		SharedConfigState: session.SharedConfigEnable,
	}
	sess, err := session.NewSessionWithOptions(s3Opts)
	if err != nil {
		return nil, err
	}
	return validateCredentials(sess)
}

// TODO: https://github.com/aws/aws-sdk-go/issues/2392
func validateCredentials(sess *session.Session) (*session.Session, error) {
	providerVal := reflect.ValueOf(*sess.Config.Credentials).FieldByName("provider").Elem()
	if providerVal.Type() == reflect.TypeOf((*credentials.ChainProvider)(nil)) {
		chainProvider := (*credentials.ChainProvider)(unsafe.Pointer(providerVal.Pointer()))
		providers := chainProvider.Providers
		if len(providers) > 0 {
			err := reflect.ValueOf(providers[0]).Elem().FieldByName("Err")
			if err.IsValid() {
				if e2, ok := err.Interface().(error); ok {
					return nil, e2
				}
			}
		}
	}
	return sess, nil
}

func ValidAbsURL(urlStr string) (*url.URL, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return u, err
	}
	if !u.IsAbs() {
		msg := fmt.Sprintf("URL '%v' must have a scheme", urlStr)
		return nil, errors.New(msg)
	}
	return u, nil
}
