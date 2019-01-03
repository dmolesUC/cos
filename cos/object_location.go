package cos

import (
	"errors"
	"fmt"
	. "net/url"
)

// ------------------------------------------------------------
// Exported types

type ObjectLocation struct {
	ObjectUrl   URL
	EndpointUrl URL
}

func (ol *ObjectLocation) String() string {
	return fmt.Sprintf("%v @ %v", ol.ObjectUrl.String(), ol.EndpointUrl.String())
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
		return nil, err
	}

	if objUrl.Scheme == "s3" {
		if paramLen < 2 {
			return nil, errors.New(fmt.Sprintf("s3 object URL '%v' requires an endpoint URL", objUrl))
		}
		endpointUrlStr := *params[1]
		endpointUrl, err := validUrl(endpointUrlStr)
		if err != nil {
			return nil, err
		}
		return NewObjectLocationFromObjectAndEndpointUrls(objUrl, endpointUrl)
	}

	return NewObjectLocationFromHttpUrl(objUrl)
}

func NewObjectLocationFromObjectAndEndpointUrls(objUrl *URL, endpointUrl *URL) (*ObjectLocation, error) {
	if objUrl.Scheme != "s3" {
		return nil, errors.New(fmt.Sprintf("object URL '%v' relative to endpoint '%v' must be S3", objUrl, endpointUrl))
	}
	if endpointUrl.Scheme == "http" || endpointUrl.Scheme == "https" {
		objLoc := ObjectLocation{
			ObjectUrl:   *objUrl,
			EndpointUrl: *endpointUrl,
		}
		return &objLoc, nil
	} else {
		return nil, errors.New(fmt.Sprintf("endpoint URL '%v' must be HTTP or HTTPS", endpointUrl))
	}
}

func NewObjectLocationFromHttpUrl(objUrl *URL) (*ObjectLocation, error) {
	objUrlScheme := objUrl.Scheme
	if objUrl.Scheme != "http" && objUrl.Scheme != "https" {
		return nil, errors.New(fmt.Sprintf("absolute object URL '%v' must be HTTP or HTTPS", objUrl))
	}

	endpointUrl, err := toEndpointUrl(objUrlScheme, objUrl.Host)
	if err != nil {
		return nil, err
	}

	s3ObjUrl, err := toS3ObjUrl(objUrl.Path)
	if err != nil {
		return nil, err
	}
	objLoc := ObjectLocation{
		ObjectUrl:   *s3ObjUrl,
		EndpointUrl: *endpointUrl,
	}
	return &objLoc, nil
}

// ------------------------------------------------------------
// Unexported functions

func toS3ObjUrl(path string) (*URL, error) {
	s3ObjUrlStr := fmt.Sprintf("s3:/%v", path)
	return Parse(s3ObjUrlStr)
}

func toEndpointUrl(scheme string, host string) (*URL, error) {
	endpointUrlStr := fmt.Sprintf("%v://%v/", scheme, host)
	return Parse(endpointUrlStr)
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
