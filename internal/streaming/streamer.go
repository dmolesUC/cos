package streaming

import (
	"fmt"
	"io"
	"time"

	"github.com/dmolesUC3/cos/internal/logging"
)

type Streamer struct {
	RangeSize     int64
	ContentLength int64
	FuncFillRange *func(byteRange *ByteRange) (bytesRead int64, err error)
	rangeCount *int64
}

func NewStreamer(rangeSize int64, contentLength int64, fillRange *func(*ByteRange) (int64, error)) (*Streamer, error) {
	if fillRange == nil {
		return nil, fmt.Errorf("stream operation's FuncFillRange cannot be nil")
	}
	s := Streamer{RangeSize: rangeSize, ContentLength: contentLength, FuncFillRange: fillRange}
	return &s, nil
}

func (s *Streamer) FillRange(byteRange *ByteRange) (bytesRead int64, err error) {
	if s.FuncFillRange == nil {
		return 0, fmt.Errorf("stream operation's FuncFillRange is nil")
	}
	fillRange := *s.FuncFillRange
	return fillRange(byteRange)
}

func (s *Streamer) RangeCount() int64 {
	var rangeCount int64
	if s.rangeCount == nil {
		rangeCount = RangeCount(s.RangeSize, s.ContentLength)
		s.rangeCount = &rangeCount
	} else {
		rangeCount = *s.rangeCount
	}
	return rangeCount
}

func (s *Streamer) NewByteRange(rangeIndex int64) *ByteRange {
	rangeSize := s.RangeSize
	contentLength := s.ContentLength
	byteRange := NewByteRange(rangeIndex, rangeSize, contentLength)
	return &byteRange
}

func (s *Streamer) StreamDown(logger logging.Logger, handleBytes func([]byte) error) (int64, error) {
	totalBytes := int64(0)
	nsStart := time.Now().UnixNano()
	nsLastUpdate := nsStart
	rangeCount := s.RangeCount()
	for rangeIndex := int64(0); rangeIndex < rangeCount; rangeIndex++ {
		byteRange := s.NewByteRange(rangeIndex)
		bytesRead, err := s.FillRange(byteRange)
		eof := err == io.EOF

		totalBytes = totalBytes + bytesRead
		nsNow := time.Now().UnixNano()
		if nsNow-nsLastUpdate > nsPerSecond || rangeIndex+1 >= s.RangeCount() || eof {
			nsLastUpdate = nsNow
			progress := Progress{
				NsElapsed:     nsNow - nsStart,
				TotalBytes:    totalBytes,
				ContentLength: s.ContentLength,
			}
			progress.InfoTo(logger)
		}

		if err != nil {
			return totalBytes, err
		}
		if bytesRead != byteRange.ExpectedBytes {
			return totalBytes, fmt.Errorf("read %d bytes, expected %v", bytesRead, byteRange)
		}

		err = handleBytes(byteRange.Buffer[0:bytesRead])
		if err != nil {
			return totalBytes, err
		}
		if eof {
			break
		}
	}
	return totalBytes, nil
}
