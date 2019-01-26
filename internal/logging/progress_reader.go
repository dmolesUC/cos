package logging

import (
	"io"
)

type ProgressReader struct {
	*ProgressReporter
	in io.Reader
}

// ------------------------------
// Factory method

func NewProgressReader(in io.Reader, expectedBytes int64) *ProgressReader {
	return &ProgressReader{
		NewProgressReporter(expectedBytes), in,
	}
}

// ------------------------------
// Reader implementation

func (r *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = r.in.Read(p)
	r.updateTotal(n)
	return
}
