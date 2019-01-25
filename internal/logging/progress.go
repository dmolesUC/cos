package logging

import (
	"math"
	"time"
)

const nsPerSecondFloat64 = float64(time.Second)

type ProgressReporter struct {
	totalBytes int64
	expectedBytes int64
	progress chan int64
}

// ------------------------------
// Factory method

func NewProgressReporter(expectedBytes int64) *ProgressReporter {
	return &ProgressReporter{
		expectedBytes: expectedBytes,
		progress: make(chan int64),
	}
}

// ------------------------------
// Exported methods

func (r *ProgressReporter) LogTo(logger Logger, interval time.Duration) {
	go monitorProgress(r.progress, r.expectedBytes, logger, interval)
}

func (r *ProgressReporter) TotalBytes() int64 {
	return r.totalBytes
}

// ------------------------------
// Unexported methods

func (r *ProgressReporter) updateTotal(additionalBytes int) {
	r.totalBytes += int64(additionalBytes)

	progress := r.progress
	if progress != nil {
		select {
		case progress <- r.totalBytes:
			// progress reported
		default:
			// wait for channel
		}
		if r.totalBytes >= r.expectedBytes {
			close(progress)
			r.progress = nil
		}
	}
}

func monitorProgress(progress chan int64, expectedBytes int64, logger Logger, interval time.Duration) {
	var currentBytes int64
	nsStart := time.Now().UnixNano()
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case latestBytes, ok := <-progress:
			if ok {
				currentBytes = latestBytes
			} else {
				return
			}
		case _ = <-ticker.C:
			logProgress(logger, nsStart, currentBytes, expectedBytes)
		}
	}
}

func logProgress(logger Logger, nsStart, currentBytes, expectedBytes int64) {
	nsElapsed := time.Now().UnixNano() - nsStart
	estBytesPerSecond := int64(math.Round((nsPerSecondFloat64 * float64(currentBytes)) / (float64(nsElapsed))))

	bytesRemaining := expectedBytes - currentBytes
	nsRemaining := int64(math.Round(nsPerSecondFloat64 * float64(bytesRemaining) / float64(estBytesPerSecond)))

	logger.Infof(
		"%v of %v (%v/s; %v elapsed, %v remaining)\n",
		FormatBytes(currentBytes),
		FormatBytes(expectedBytes),
		FormatBytes(estBytesPerSecond),
		FormatNanos(nsElapsed),
		FormatNanos(nsRemaining),
	)
}
