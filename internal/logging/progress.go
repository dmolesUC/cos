package logging

import (
	"math"
	"sync/atomic"
	"time"
)

const nsPerSecondFloat64 = float64(time.Second)

type ProgressReporter struct {
	totalBytes    *int64
	expectedBytes int64
}

// ------------------------------
// Factory method

func NewProgressReporter(expectedBytes int64) *ProgressReporter {
	zero := int64(0)
	return &ProgressReporter{
		totalBytes: &zero,
		expectedBytes: expectedBytes,
	}
}

// ------------------------------
// Exported methods

func (r *ProgressReporter) LogTo(logger Logger, interval time.Duration) {
	go r.monitorProgress(logger, interval)
	// go monitorProgress(r.progress, r.expectedBytes, logger, interval)
}

func (r *ProgressReporter) TotalBytes() int64 {
	return atomic.LoadInt64(r.totalBytes)
}

func (r *ProgressReporter) ExpectedBytes() int64 {
	return r.expectedBytes
}

// ------------------------------
// Unexported methods

func (r *ProgressReporter) updateTotal(additionalBytes int) {
	atomic.AddInt64(r.totalBytes, int64(additionalBytes))
}

func (r *ProgressReporter) monitorProgress(logger Logger, interval time.Duration) {
	expectedBytes := r.expectedBytes
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	nsStart := time.Now().UnixNano()
	for {
		select {
		case _ = <-ticker.C:
			currentBytes := r.TotalBytes()
			logProgress(logger, nsStart, currentBytes, expectedBytes)
			if currentBytes >= expectedBytes {
				return
			}
		default:
			// wait for next tick
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
