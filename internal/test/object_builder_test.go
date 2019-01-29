package test

import (
	"net/url"

	. "gopkg.in/check.v1"

	"github.com/dmolesUC3/cos/internal/logging"
	. "github.com/dmolesUC3/cos/internal/objects"
)

type ObjectSuite struct {
	logger logging.Logger
}

var _ = Suite(&ObjectSuite{})

func (s *ObjectSuite) SetUpSuite(c *C) {
	s.logger = logging.NewLogger(logging.Default)
}

func (s *ObjectSuite) TestParsingHttpObjectURL(c *C) {
	inputURLStr := "https://s3-cn-north-1.amazonaws.com/example.org/foo/bar/baz.qux"
	// expectedRegion := "cn-north-1"
	expectedKey := "/foo/bar/baz.qux"
	expectedBucket := "example.org"
	expectedEndpoint, _ := url.Parse("https://s3-cn-north-1.amazonaws.com/")

	b := NewObjectBuilder().WithObjectURLStr(inputURLStr)
	o, err := b.Build(s.logger)
	c.Assert(err, IsNil)
	c.Assert(*o.Key(), Equals, expectedKey)
	c.Assert(*o.Bucket(), Equals, expectedBucket)
	c.Assert(*o.Endpoint(), Equals, *expectedEndpoint)
}

func (s *ObjectSuite) TestParsingHttpObjectURLEmptyEndpoint(c *C) {
	inputURLStr := "https://s3-cn-north-1.amazonaws.com/example.org/foo/bar/baz.qux"
	// expectedRegion := "cn-north-1"
	expectedKey := "/foo/bar/baz.qux"
	expectedBucket := "example.org"
	expectedEndpoint, _ := url.Parse("https://s3-cn-north-1.amazonaws.com/")

	b := NewObjectBuilder().
		WithObjectURLStr(inputURLStr).
		WithEndpointStr("")
	o, err := b.Build(s.logger)
	c.Assert(err, IsNil)
	//c.Assert(*o.Region(), Equals, expectedRegion)
	c.Assert(*o.Key(), Equals, expectedKey)
	c.Assert(*o.Bucket(), Equals, expectedBucket)
	c.Assert(*o.Endpoint(), Equals, *expectedEndpoint)
}

func (s *ObjectSuite) TestParsingS3ObjectURL(c *C) {
	inputURLStr := "s3://example.org/foo/bar/baz.qux"
	// expectedRegion := "cn-north-1"
	expectedKey := "/foo/bar/baz.qux"
	expectedBucket := "example.org"
	expectedEndpoint, _ := url.Parse("https://s3-cn-north-1.amazonaws.com/")

	b := NewObjectBuilder().
		WithObjectURLStr(inputURLStr).
		WithEndpoint(expectedEndpoint)

	o, err := b.Build(s.logger)
	c.Assert(err, IsNil)
	//c.Assert(*o.Region(), Equals, expectedRegion)
	c.Assert(*o.Key(), Equals, expectedKey)
	c.Assert(*o.Bucket(), Equals, expectedBucket)
	c.Assert(*o.Endpoint(), Equals, *expectedEndpoint)
}

func (s *ObjectSuite) TestValidationFailureNoEndpoint(c *C) {
	expectedKey := "/foo/bar/baz.qux"
	expectedBucket := "example.org"

	b := NewObjectBuilder().
		WithKey(expectedKey).
		WithBucket(expectedBucket)

	_, err := b.Build(s.logger)
	c.Assert(err, ErrorMatches, ".*endpoint.*")
}

func (s *ObjectSuite) TestValidationFailureNoBucket(c *C) {
	expectedKey := "/foo/bar/baz.qux"
	expectedEndpoint, _ := url.Parse("https://s3-cn-north-1.amazonaws.com/")

	b := NewObjectBuilder().
		WithKey(expectedKey).
		WithEndpoint(expectedEndpoint)

	_, err := b.Build(s.logger)
	c.Assert(err, ErrorMatches, ".*bucket.*")
}

func (s *ObjectSuite) TestValidationFailureNoKey(c *C) {
	expectedBucket := "example.org"
	expectedEndpoint, _ := url.Parse("https://s3-cn-north-1.amazonaws.com/")

	b := NewObjectBuilder().
		WithBucket(expectedBucket).
		WithEndpoint(expectedEndpoint)

	_, err := b.Build(s.logger)
	c.Assert(err, ErrorMatches, ".*key.*")
}
