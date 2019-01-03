package cos

import (
	"fmt"
	. "net/url"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

// ------------------------------------------------------------
// NewFromObjectAndEndpointUrls

type NewFromObjectAndEndpointUrls struct{}

var _ = Suite(&NewFromObjectAndEndpointUrls{})

func (s *NewFromObjectAndEndpointUrls) TestHttpsEndpointUrl(c *C) {
	objectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	objUrl, err := Parse(objectUrlStr)
	c.Assert(err, IsNil) // just to be sure

	endpointUrlStr := "https://s3-us-west-2.amazonaws.com/"
	endpointUrl, err := Parse(endpointUrlStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromObjectAndEndpointUrls(objUrl, endpointUrl)
	c.Assert(err, IsNil)

	c.Assert(&objLoc.ObjectUrl, DeepEquals, objUrl)
	c.Assert(&objLoc.EndpointUrl, DeepEquals,endpointUrl)
}

func (s *NewFromObjectAndEndpointUrls) TestHttpEndpointUrl(c *C) {
	objectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	objUrl, err := Parse(objectUrlStr)
	c.Assert(err, IsNil) // just to be sure

	endpointUrlStr := "http://s3-us-west-2.amazonaws.com/"
	endpointUrl, err := Parse(endpointUrlStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromObjectAndEndpointUrls(objUrl, endpointUrl)
	c.Assert(err, IsNil)

	c.Assert(&objLoc.ObjectUrl, DeepEquals, objUrl)
	c.Assert(&objLoc.EndpointUrl, DeepEquals,endpointUrl)
}

func (s *NewFromObjectAndEndpointUrls) TestBadEndpointUrl(c *C) {
	objectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	objUrl, err := Parse(objectUrlStr)
	c.Assert(err, IsNil) // just to be sure

	endpointUrlStr := "s3://us-west-2.amazonaws.com/"
	endpointUrl, err := Parse(endpointUrlStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromObjectAndEndpointUrls(objUrl, endpointUrl)
	c.Assert(objLoc, IsNil)
	c.Assert(err, ErrorMatches, fmt.Sprintf(".*%v.*", endpointUrlStr))
}

func (s *NewFromObjectAndEndpointUrls) TestBadObjectUrl(c *C) {
	objectUrlStr := "http://www.dmoles.net/images/fa/archive.svg"
	objUrl, err := Parse(objectUrlStr)
	c.Assert(err, IsNil) // just to be sure

	endpointUrlStr := "https://s3-us-west-2.amazonaws.com/"
	endpointUrl, err := Parse(endpointUrlStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromObjectAndEndpointUrls(objUrl, endpointUrl)
	c.Assert(objLoc, IsNil)
	c.Assert(err, ErrorMatches, fmt.Sprintf(".*%v.*", objectUrlStr))
}

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

// ------------------------------------------------------------
// NewFromStrings

type NewFromStrings struct{}

var _ = Suite(&NewFromStrings{})

func (s *NewFromStrings) TestHttpsObjectUrl(c *C) {
	inputUrlStr := "https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedObjectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointUrlStr := "https://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&inputUrlStr)
	c.Assert(err, IsNil)

	c.Assert(objLoc.ObjectUrl.String(), Equals, expectedObjectUrlStr)
	c.Assert(objLoc.EndpointUrl.String(), Equals, expectedEndpointUrlStr)
}

func (s *NewFromStrings) TestHttpObjectUrl(c *C) {
	inputUrlStr := "http://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedObjectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointUrlStr := "http://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&inputUrlStr)
	c.Assert(err, IsNil)

	c.Assert(objLoc.ObjectUrl.String(), Equals, expectedObjectUrlStr)
	c.Assert(objLoc.EndpointUrl.String(), Equals, expectedEndpointUrlStr)
}

func (s *NewFromStrings) TestHttpsObjectUrlWithNilEndpoint(c *C) {
	inputUrlStr := "https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedObjectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointUrlStr := "https://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&inputUrlStr, nil)
	c.Assert(err, IsNil)

	c.Assert(objLoc.ObjectUrl.String(), Equals, expectedObjectUrlStr)
	c.Assert(objLoc.EndpointUrl.String(), Equals, expectedEndpointUrlStr)
}


func (s *NewFromStrings) TestHttpsEndpointUrl(c *C) {
	objectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	endpointUrlStr := "https://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&objectUrlStr, &endpointUrlStr)
	c.Assert(err, IsNil)

	c.Assert(objLoc.ObjectUrl.String(), Equals, objectUrlStr)
	c.Assert(objLoc.EndpointUrl.String(), Equals,endpointUrlStr)
}

func (s *NewFromStrings) TestHttpEndpointUrl(c *C) {
	objectUrlStr := "s3://www.dmoles.net/images/fa/archive.svg"
	endpointUrlStr := "http://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&objectUrlStr, &endpointUrlStr)
	c.Assert(err, IsNil)

	c.Assert(objLoc.ObjectUrl.String(), Equals, objectUrlStr)
	c.Assert(objLoc.EndpointUrl.String(), Equals,endpointUrlStr)
}
