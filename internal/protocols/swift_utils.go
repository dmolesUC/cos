package protocols

import (
	"fmt"
	"net/url"

	"github.com/dmolesUC3/cos/internal/logging"
)

// The SwiftConnectionParams struct encapsulates the authentication parameters
// needed to connect to an Openstack Swift server.
type SwiftConnectionParams struct {
	UserName string
	APIKey   string
	AuthURL  *url.URL
}

func (p SwiftConnectionParams) String() string {
	return fmt.Sprintf(
		"SwiftConnectionParams { Username: %v, APIKey: %v, AuthURL: %v }",
		p.UserName, p.apiKeyStr(), p.authUrlStr(),
	)
}

func (p SwiftConnectionParams) Pretty() string {
	format := `SwiftConnectionParams {
		Username: %v
		APIKey: %v
		AuthURL: %v
	}`
	format = logging.Untabify(format, "    ")
	return fmt.Sprintf(format, p.UserName, p.apiKeyStr(), p.authUrlStr())
}

func (p SwiftConnectionParams) apiKeyStr() string {
	var apiKeyStr string
	if p.APIKey == "" {
		apiKeyStr = "<not set>"
	} else {
		apiKeyStr = "<hidden>"
	}
	return apiKeyStr
}

func (p SwiftConnectionParams) authUrlStr() string {
	var authURLStr string
	if p.AuthURL == nil {
		authURLStr = "<nil>"
	} else {
		authURLStr = p.AuthURL.String()
	}
	return authURLStr
}
