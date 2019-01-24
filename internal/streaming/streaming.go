package streaming

import (
	"fmt"
	"io"
)

func NextRange(currentTotal int64, maxRangeSize int64, contentLength int64) (start, end int64, size int) {
	start = currentTotal
	end = currentTotal + maxRangeSize
	if end > contentLength {
		end = contentLength - 1
	}
	size = int((end + 1) - currentTotal)
	return start, end, size
}

// ReadExactly reads exactly the number of bytes to fill the specified buffer,
// otherwise returning an error.
func ReadExactly(in io.Reader, buffer []byte) (err error) {
	bytesRead, err := io.ReadFull(in, buffer)
	if err == nil {
		expected := len(buffer)
		if bytesRead != expected {
			err = fmt.Errorf("expected to read %d bytes, got %d", expected, bytesRead)
		}
	}
	return
}

// WriteExactly writes exactly the number of bytes found in the specified buffer,
// otherwise returning an error.
func WriteExactly(out io.Writer, data []byte) (err error) {
	bytesWritten, err := out.Write(data)
	if err == nil {
		expected := len(data)
		if bytesWritten != expected {
			err = fmt.Errorf("expected to write %d bytes, got %d", expected, bytesWritten)
		}
	}
	return
}
