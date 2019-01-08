package internal

import (
	"errors"
	"fmt"
	. "net/url"
)

// ------------------------------------------------------------
// Exported types

type ObjectLocation struct {
	S3Uri    URL
	Endpoint URL
}

func (ol ObjectLocation) String() string {
	return fmt.Sprintf("%#v@%#v", ol.S3Uri.String(), ol.Endpoint.String())
}

func (ol ObjectLocation) Bucket() string {
	return ol.S3Uri.Host
}

func (ol ObjectLocation) Key() string {
	return ol.S3Uri.Path
}

// ------------------------------------------------------------
// Exported functions

func NewObjectLocationFromStrings(params ...*string) (*ObjectLocation, error) {
	paramLen := len(params)
	if paramLen < 1 {
		return nil, errors.New("no object or endpoint URL provided")
	}
	if paramLen > 2 {
		return nil, errors.New(fmt.Sprintf("too many params: expected [object-url <endpoint-url>], got %v", params))
	}

	objUrlStr := *params[0]
	objUrl, err := validUrl(objUrlStr)
	if err != nil {
		err = fmt.Errorf("error parsing object URL: %v", err)
		return nil, err
	}

	if objUrl.Scheme == "s3" {
		if paramLen > 1 {
			endpointStr := *params[1]
			if endpointStr != "" {
				endpoint, err := validUrl(endpointStr)
				if err != nil {
					err = fmt.Errorf("error parsing endpoint URL: %v", err)
					return nil, err
				}
				return NewObjectLocationFromS3UriAndEndpoint(objUrl, endpoint)
			}
		}
		return nil, fmt.Errorf("s3 object URL '%v' requires an endpoint URL (e.g. '%v')", objUrl, DefaultS3EndpointUrl)
	}

	return NewObjectLocationFromHttpUrl(objUrl)
}

func NewObjectLocationFromS3UriAndEndpoint(s3Uri *URL, endpoint *URL) (*ObjectLocation, error) {
	if s3Uri.Scheme != "s3" {
		return nil, errors.New(fmt.Sprintf("object URL '%v' relative to endpoint '%v' must be S3", s3Uri, endpoint))
	}
	if endpoint.Scheme == "http" || endpoint.Scheme == "https" {
		objLoc := ObjectLocation{
			S3Uri:    *s3Uri,
			Endpoint: *endpoint,
		}
		return &objLoc, nil
	} else {
		return nil, errors.New(fmt.Sprintf("endpoint URL '%v' must be HTTP or HTTPS", endpoint))
	}
}

func NewObjectLocationFromHttpUrl(s3Uri *URL) (*ObjectLocation, error) {
	s3UriScheme := s3Uri.Scheme
	if s3Uri.Scheme != "http" && s3Uri.Scheme != "https" {
		return nil, errors.New(fmt.Sprintf("absolute object URL '%v' must be HTTP or HTTPS", s3Uri))
	}

	endpoint, err := toEndpoint(s3UriScheme, s3Uri.Host)
	if err != nil {
		return nil, err
	}

	s3S3Uri, err := toS3S3Uri(s3Uri.Path)
	if err != nil {
		return nil, err
	}
	objLoc := ObjectLocation{
		S3Uri:    *s3S3Uri,
		Endpoint: *endpoint,
	}
	return &objLoc, nil
}

// ------------------------------------------------------------
// Unexported functions

func toS3S3Uri(path string) (*URL, error) {
	s3S3UriStr := fmt.Sprintf("s3:/%v", path)
	return Parse(s3S3UriStr)
}

func toEndpoint(scheme string, host string) (*URL, error) {
	endpointStr := fmt.Sprintf("%v://%v/", scheme, host)
	return Parse(endpointStr)
}

func validUrl(urlStr string) (*URL, error) {
	url, err := Parse(urlStr)
	if err != nil {
		return url, err
	}
	if !url.IsAbs() {
		msg := fmt.Sprintf("URL '%v' must have a scheme", urlStr)
		return nil, errors.New(msg)
	}
	return url, nil
}
