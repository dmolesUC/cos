package logging

import (
	"io"
	"time"
)

type ProgressReader struct {
	totalBytesRead int64
	expectedBytes  int64
	in             io.Reader
	progress       chan int64
}

// ------------------------------
// Factory method

func NewProgressReader(in io.Reader, expectedBytes int64) *ProgressReader {
	return &ProgressReader{in: in, expectedBytes: expectedBytes, progress: make(chan int64)}
}

// ------------------------------
// Accessors

func (r *ProgressReader) TotalBytesRead() int64 {
	return r.totalBytesRead
}

// ------------------------------
// Public methods

func (r *ProgressReader) LogTo(logger Logger, interval time.Duration) {
	go LogProgress(r.progress, r.expectedBytes, logger, interval)

	go func() {
		for {
			select {
			case r.progress <- r.totalBytesRead:
				// progress reported
			default:
				// wait for channel
			}

			if r.totalBytesRead >= r.expectedBytes {
				close(r.progress)
				break
			}
		}
	}()
}

// ------------------------------
// Reader implementation

func (r *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = r.in.Read(p)
	r.totalBytesRead += int64(n)
	return
}
