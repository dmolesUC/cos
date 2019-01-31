package objects

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/ncw/swift"
)

const (
	defaultRetries = 3
)

// ------------------------------------------------------------
// SwiftTarget type

type SwiftTarget struct {
	UserName string
	APIKey   string
	AuthURL  *url.URL
	Container string
	cnx      *swift.Connection
}

// ------------------------------
// Factory method

func NewSwiftEndpoint(endpointUrl *url.URL, container string) (*SwiftTarget, error) {
	swiftAPIUser := os.Getenv("SWIFT_API_USER")
	if swiftAPIUser == "" {
		return nil, errors.New("missing environment variable $SWIFT_API_USER")
	}
	swiftAPIKey := os.Getenv("SWIFT_API_KEY")
	if swiftAPIKey == "" {
		return nil, errors.New("missing environment variable $SWIFT_API_KEY")
	}
	return &SwiftTarget{UserName: swiftAPIUser, APIKey: swiftAPIKey, AuthURL: endpointUrl, Container: container}, nil
}

// ------------------------------
// Target implementation

func (e *SwiftTarget) Object(key string) Object {
	return &SwiftObject{e, e.Container, key}
}

func (e *SwiftTarget) Pretty() string {
	var apiKeyStr string
	if e.APIKey == "" {
		apiKeyStr = "<not set>"
	} else {
		apiKeyStr = "<hidden>"
	}

	var authURLStr string
	if e.AuthURL == nil {
		authURLStr = "<nil>"
	} else {
		authURLStr = e.AuthURL.String()
	}

	return fmt.Sprintf("SwiftTarget { Username: %#v, APIKey: %v, AuthURL: %#v, Container: %#v }",
		e.UserName, apiKeyStr, authURLStr, e.Container)
}

func (e *SwiftTarget) String() string {
	return e.Pretty()
}

// ------------------------------
// Miscellaneous methods

func (e *SwiftTarget) Connection() (*swift.Connection, error) {
	if e.cnx == nil {
		authUrl := e.AuthURL
		if authUrl == nil {
			return nil, fmt.Errorf("authUrl not set in SwiftTarget: %v", e)
		}
		authUrlStr := authUrl.String()
		e.cnx = &swift.Connection{
			UserName: e.UserName,
			ApiKey:   e.APIKey,
			AuthUrl:  authUrlStr,
			Retries:  defaultRetries,
		}
	}
	return e.cnx, nil
}
