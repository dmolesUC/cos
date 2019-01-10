package pkg

import "net/url"

// The Object type represents the location of an object in cloud storage.
type Object interface {
	// TODO: push Region() down to S3Object
	Region() *string
	Endpoint() *url.URL
	Bucket() *string
	Key() *string
	CalcDigest(downloadRangeSize int64, algorithm string) ([] byte, error)
	StreamDown(chunkSize int64, handleChunk func([]byte) error) (int64, error)
}

