package cos

import (
	"fmt"
	. "net/url"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

// ------------------------------------------------------------
// NewFromHttpUrl

type NewFromHttpUrl struct{}

var _ = Suite(&NewFromHttpUrl{})

func (s *NewFromHttpUrl) TestHttpsUrl(c *C) {
	inputUrlStr := "https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedObjectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointUrlStr := "https://s3-us-west-2.amazonaws.com/"

	objUrl, err := Parse(inputUrlStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpUrl(objUrl)
	c.Assert(err, IsNil)

	objUrlStr := objLoc.ObjectUrl.String()
	c.Assert(objUrlStr, Equals, expectedObjectUrlStr)

	endpointUrlStr := objLoc.EndpointUrl.String()
	c.Assert(endpointUrlStr, Equals, expectedEndpointUrlStr)
}

func (s *NewFromHttpUrl) TestHttpUrl(c *C) {
	inputUrlStr := "http://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedObjectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointUrlStr := "http://s3-us-west-2.amazonaws.com/"

	objUrl, err := Parse(inputUrlStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpUrl(objUrl)
	c.Assert(err, IsNil)

	objUrlStr := objLoc.ObjectUrl.String()
	c.Assert(objUrlStr, Equals, expectedObjectUrlStr)

	endpointUrlStr := objLoc.EndpointUrl.String()
	c.Assert(endpointUrlStr, Equals, expectedEndpointUrlStr)
}

func (s *NewFromHttpUrl) TestEscapedCharsInObjectPath(c *C) {
	inputUrlStr := "http://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive%201.svg"
	expectedObjUrlStr := "s3://www.dmoles.net/images/fa/archive%201.svg"
	expectedEndpointUrlStr := "http://s3-us-west-2.amazonaws.com/"

	objUrl, err := Parse(inputUrlStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpUrl(objUrl)
	c.Assert(err, IsNil)

	objUrlStr := objLoc.ObjectUrl.String()
	c.Assert(objUrlStr, Equals, expectedObjUrlStr)

	endpointUrlStr := objLoc.EndpointUrl.String()
	c.Assert(endpointUrlStr, Equals, expectedEndpointUrlStr)
}

func (s *NewFromHttpUrl) TestNonHttpUrl(c *C) {
	objUrl, err := Parse("s3://www.dmoles.net/images/fa/archive.svg")
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpUrl(objUrl)
	c.Assert(objLoc, IsNil)
	c.Assert(err, ErrorMatches, fmt.Sprintf(".*%v.*", objUrl))
}

func (s *NewFromHttpUrl) TestInvalidCharsInBucketName(c *C) {
	inputUrlStr := "https://example.org/{}/images/fa/archive.svg"
	expectedObjUrlStr := "s3://{}/images/fa/archive.svg"

	objUrl, err := Parse(inputUrlStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpUrl(objUrl)
	c.Assert(objLoc, IsNil)
	c.Assert(err, ErrorMatches, fmt.Sprintf(".*%v.*", expectedObjUrlStr))
}