package streaming

import (
	"fmt"
)

type ByteRange struct {
	RangeIndex int64
	RangeCount int64
	StartInclusive int64
	EndInclusive int64
	ExpectedBytes int64
	Buffer []byte
}

func (r ByteRange) String() string {
	return fmt.Sprintf("%d bytes [%d, %d] (range %d of %d)", r.RangeIndex, r.RangeCount, r.ExpectedBytes, r.StartInclusive, r.EndInclusive)
}

func NewByteRange(rangeIndex int64, rangeSize int64, contentLength int64) ByteRange {
	rangeCount := RangeCount(rangeSize, contentLength)
	// byte ranges are 0-indexed and inclusive
	startInclusive := rangeIndex * rangeSize
	var endInclusive int64
	if (rangeIndex + 1) < rangeCount {
		endInclusive = startInclusive + rangeSize - 1
	} else {
		endInclusive = contentLength - 1
	}
	expectedBytes := (endInclusive + 1) - startInclusive
	byteRange := ByteRange{
		RangeIndex:     rangeIndex,
		RangeCount:     rangeCount,
		StartInclusive: startInclusive,
		EndInclusive:   endInclusive,
		ExpectedBytes:  expectedBytes,
		Buffer:         make([]byte, expectedBytes),
	}
	return byteRange
}

func RangeCount(rangeSize int64, contentLength int64) int64 {
	return (contentLength + rangeSize - 1) / rangeSize
}
