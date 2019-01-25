package logging

import (
	"io"
	"time"
)

type ProgressWriter struct {
	totalBytesWritten int64
	expectedBytes     int64
	out               io.Writer
	progress          chan int64
}

// ------------------------------
// Factory method

func NewProgressWriter(out io.Writer, expectedBytes int64) *ProgressWriter {
	return &ProgressWriter{out: out, expectedBytes: expectedBytes, progress: make(chan int64)}
}

// ------------------------------
// Accessors

func (w *ProgressWriter) TotalBytesWritten() int64 {
	return w.totalBytesWritten
}

// ------------------------------
// Public methods

func (w *ProgressWriter) LogTo(logger Logger, interval time.Duration) {
	go LogProgress(w.progress, w.expectedBytes, logger, interval)

	go func() {
		for {
			select {
			case w.progress <- w.totalBytesWritten:
				// progress reported
			default:
				// wait for channel
			}

			if w.totalBytesWritten >= w.expectedBytes {
				close(w.progress)
				break
			}
		}
	}()
}

// ------------------------------
// Writer implementation

func (w *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = w.out.Write(p)
	w.totalBytesWritten += int64(n)
	return
}
