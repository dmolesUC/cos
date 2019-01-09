package internal

import (
	"fmt"
	. "net/url"

	. "gopkg.in/check.v1"
)

// ------------------------------------------------------------
// NewFromS3UriAndEndpoint

type NewFromS3UriAndEndpoint struct{}

var _ = Suite(&NewFromS3UriAndEndpoint{})

func (s *NewFromS3UriAndEndpoint) TestHttpsEndpoint(c *C) {
	s3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	s3Uri, err := Parse(s3UriStr)
	c.Assert(err, IsNil) // just to be sure

	endpointStr := "https://s3-us-west-2.amazonaws.com/"
	endpoint, err := Parse(endpointStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromS3UriAndEndpoint(s3Uri, endpoint)
	c.Assert(err, IsNil)

	c.Assert(&objLoc.S3Uri, DeepEquals, s3Uri)
	c.Assert(&objLoc.Endpoint, DeepEquals,endpoint)
}

func (s *NewFromS3UriAndEndpoint) TestHttpEndpoint(c *C) {
	s3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	s3Uri, err := Parse(s3UriStr)
	c.Assert(err, IsNil) // just to be sure

	endpointStr := "http://s3-us-west-2.amazonaws.com/"
	endpoint, err := Parse(endpointStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromS3UriAndEndpoint(s3Uri, endpoint)
	c.Assert(err, IsNil)

	c.Assert(&objLoc.S3Uri, DeepEquals, s3Uri)
	c.Assert(&objLoc.Endpoint, DeepEquals,endpoint)
}

func (s *NewFromS3UriAndEndpoint) TestBadEndpoint(c *C) {
	s3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	s3Uri, err := Parse(s3UriStr)
	c.Assert(err, IsNil) // just to be sure

	endpointStr := "s3://us-west-2.amazonaws.com/"
	endpoint, err := Parse(endpointStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromS3UriAndEndpoint(s3Uri, endpoint)
	c.Assert(objLoc, IsNil)
	c.Assert(err, ErrorMatches, fmt.Sprintf(".*%v.*", endpointStr))
}

func (s *NewFromS3UriAndEndpoint) TestBadS3Uri(c *C) {
	s3UriStr := "http://www.dmoles.net/images/fa/archive.svg"
	s3Uri, err := Parse(s3UriStr)
	c.Assert(err, IsNil) // just to be sure

	endpointStr := "https://s3-us-west-2.amazonaws.com/"
	endpoint, err := Parse(endpointStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromS3UriAndEndpoint(s3Uri, endpoint)
	c.Assert(objLoc, IsNil)
	c.Assert(err, ErrorMatches, fmt.Sprintf(".*%v.*", s3UriStr))
}

// ------------------------------------------------------------
// NewFromHttpURL

type NewFromHttpURL struct{}

var _ = Suite(&NewFromHttpURL{})

func (s *NewFromHttpURL) TestHttpsURL(c *C) {
	inputURLStr := "https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedS3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointStr := "https://s3-us-west-2.amazonaws.com/"

	s3Uri, err := Parse(inputURLStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpURL(s3Uri)
	c.Assert(err, IsNil)

	s3UriStr := objLoc.S3Uri.String()
	c.Assert(s3UriStr, Equals, expectedS3UriStr)

	endpointStr := objLoc.Endpoint.String()
	c.Assert(endpointStr, Equals, expectedEndpointStr)
}

func (s *NewFromHttpURL) TestHttpURL(c *C) {
	inputURLStr := "http://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedS3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointStr := "http://s3-us-west-2.amazonaws.com/"

	s3Uri, err := Parse(inputURLStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpURL(s3Uri)
	c.Assert(err, IsNil)

	s3UriStr := objLoc.S3Uri.String()
	c.Assert(s3UriStr, Equals, expectedS3UriStr)

	endpointStr := objLoc.Endpoint.String()
	c.Assert(endpointStr, Equals, expectedEndpointStr)
}

func (s *NewFromHttpURL) TestEscapedCharsInObjectPath(c *C) {
	inputURLStr := "http://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive%201.svg"
	expectedS3UriStr := "s3://www.dmoles.net/images/fa/archive%201.svg"
	expectedEndpointStr := "http://s3-us-west-2.amazonaws.com/"

	s3Uri, err := Parse(inputURLStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpURL(s3Uri)
	c.Assert(err, IsNil)

	s3UriStr := objLoc.S3Uri.String()
	c.Assert(s3UriStr, Equals, expectedS3UriStr)

	endpointStr := objLoc.Endpoint.String()
	c.Assert(endpointStr, Equals, expectedEndpointStr)
}

func (s *NewFromHttpURL) TestNonHttpURL(c *C) {
	s3Uri, err := Parse("s3://www.dmoles.net/images/fa/archive.svg")
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpURL(s3Uri)
	c.Assert(objLoc, IsNil)
	c.Assert(err, ErrorMatches, fmt.Sprintf(".*%v.*", s3Uri))
}

func (s *NewFromHttpURL) TestInvalidCharsInBucketName(c *C) {
	inputURLStr := "https://example.org/{}/images/fa/archive.svg"
	expectedS3UriStr := "s3://{}/images/fa/archive.svg"

	s3Uri, err := Parse(inputURLStr)
	c.Assert(err, IsNil) // just to be sure

	objLoc, err := NewObjectLocationFromHttpURL(s3Uri)
	c.Assert(objLoc, IsNil)
	c.Assert(err, ErrorMatches, fmt.Sprintf(".*%v.*", expectedS3UriStr))
}

// ------------------------------------------------------------
// NewFromStrings

type NewFromStrings struct{}

var _ = Suite(&NewFromStrings{})

func (s *NewFromStrings) TestHttpsS3Uri(c *C) {
	inputURLStr := "https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedS3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointStr := "https://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&inputURLStr)
	c.Assert(err, IsNil)

	c.Assert(objLoc.S3Uri.String(), Equals, expectedS3UriStr)
	c.Assert(objLoc.Endpoint.String(), Equals, expectedEndpointStr)
}

func (s *NewFromStrings) TestHttpS3Uri(c *C) {
	inputURLStr := "http://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedS3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointStr := "http://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&inputURLStr)
	c.Assert(err, IsNil)

	c.Assert(objLoc.S3Uri.String(), Equals, expectedS3UriStr)
	c.Assert(objLoc.Endpoint.String(), Equals, expectedEndpointStr)
}

func (s *NewFromStrings) TestHttpsS3UriWithNilEndpoint(c *C) {
	inputURLStr := "https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg"
	expectedS3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	expectedEndpointStr := "https://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&inputURLStr, nil)
	c.Assert(err, IsNil)

	c.Assert(objLoc.S3Uri.String(), Equals, expectedS3UriStr)
	c.Assert(objLoc.Endpoint.String(), Equals, expectedEndpointStr)
}


func (s *NewFromStrings) TestHttpsEndpoint(c *C) {
	s3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	endpointStr := "https://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&s3UriStr, &endpointStr)
	c.Assert(err, IsNil)

	c.Assert(objLoc.S3Uri.String(), Equals, s3UriStr)
	c.Assert(objLoc.Endpoint.String(), Equals,endpointStr)
}

func (s *NewFromStrings) TestHttpEndpoint(c *C) {
	s3UriStr := "s3://www.dmoles.net/images/fa/archive.svg"
	endpointStr := "http://s3-us-west-2.amazonaws.com/"

	objLoc, err := NewObjectLocationFromStrings(&s3UriStr, &endpointStr)
	c.Assert(err, IsNil)

	c.Assert(objLoc.S3Uri.String(), Equals, s3UriStr)
	c.Assert(objLoc.Endpoint.String(), Equals,endpointStr)
}