package logging

import (
	"io"
)

type ProgressWriter struct {
	*ProgressReporter
	out               io.Writer
}

// ------------------------------
// Factory method

func NewProgressWriter(out io.Writer, expectedBytes int64) *ProgressWriter {
	return &ProgressWriter{
		NewProgressReporter(expectedBytes), out,
	}
}

// ------------------------------
// Writer implementation

func (w *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = w.out.Write(p)
	w.updateTotal(n)
	return
}
