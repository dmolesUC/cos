package internal

import (
	"reflect"
	"regexp"
	"unsafe"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	DefaultAwsRegion     = "us-west-2"
	DefaultS3EndpointUrl = "https://s3-us-west-2.amazonaws.com"
	awsRegionRegexpStr = "https?://s3-([^.]+)\\.amazonaws\\.com/"
)
var awsRegionRegexp = regexp.MustCompile(awsRegionRegexpStr)

func ExtractRegion(endpoint string, logger Logger) string {
	matches := awsRegionRegexp.FindStringSubmatch(endpoint)
	regionStr := DefaultAwsRegion
	if len(matches) == 2 {
		regionStr = matches[1]
		logger.Detailf("Found AWS region in endpoint URL %v: %v\n", endpoint, regionStr)
	} else {
		logger.Detailf("No AWS region found in endpoint URL %v; using default region %v\n", endpoint, DefaultAwsRegion)
	}
	return regionStr
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
