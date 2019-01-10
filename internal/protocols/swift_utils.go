package protocols

import (
	"fmt"
	"net/url"
)

// The SwiftConnectionParams struct encapsulates the authentication parameters
// needed to connect to an Openstack Swift server.
type SwiftConnectionParams struct {
	UserName string
	APIKey   string
	AuthURL  *url.URL
}

func (p SwiftConnectionParams) String() string {
	var authURLStr string
	if p.AuthURL == nil {
		authURLStr = "<nil>"
	} else {
		authURLStr = p.AuthURL.String()
	}

	return fmt.Sprintf(
		"SwiftConnectionParams { Username: %v, APIKey: %v, AuthURL: %v }",
		p.UserName, p.APIKey, authURLStr,
	)
}
