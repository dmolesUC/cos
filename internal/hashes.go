package internal

import (
	"crypto/md5"
	"crypto/sha256"
	"hash"
)

// NewHash returns a new hash of the specified algorithm ("sha256" or "md5")
func NewHash(algorithm string) hash.Hash {
	if algorithm == "sha256" {
		return sha256.New()
	}
	return md5.New()
}
