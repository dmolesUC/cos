package streaming

import "io"

type CountWriter struct {
	BytesWritten int64
	out          io.Writer
}

func (w *CountWriter) Write(p []byte) (n int, err error) {
	n, err = w.out.Write(p)
	w.BytesWritten += int64(n)
	return
}

func NewCountWriter(out io.Writer) *CountWriter {
	return &CountWriter{out: out}
}

type CountReader struct {
	BytesRead int64
	in        io.Reader
}

func (w *CountReader) Read(p []byte) (n int, err error) {
	n, err = w.in.Read(p)
	w.BytesRead += int64(n)
	return
}

func NewCountReader(in io.Reader) *CountReader {
	return &CountReader{in: in}
}
