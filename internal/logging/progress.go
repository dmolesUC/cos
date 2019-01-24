package logging

import (
	"math"
	"sync"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

const nsPerSecondFloat64 = float64(time.Second)

// ------------------------------------------------------------
// Exported functions

func ReportProgress(expectedBytes int64, logger Logger, interval time.Duration) (progress chan int64) {
	nsStart := time.Now().UnixNano()
	ticker := time.NewTicker(interval)
	go func() {
		var currentBytes int64
		for {
			select {
			case latestBytes, ok := <-progress:
				if ok {
					currentBytes = latestBytes
				} else {
					return
				}
			case _ = <-ticker.C:
				ProgressInfoTo(logger, nsStart, currentBytes, expectedBytes)
			}
		}
	}()
	return
}

func ProgressInfoTo(logger Logger, nsStart, currentBytes, expectedBytes int64) {
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

// ------------------------------------------------------------
// Progress

// Deprecated: Use ProgressInfoTo()
type Progress struct {
	NsElapsed     int64
	TotalBytes    int64
	ContentLength int64
	estimatedBps  *int64
}

func (p *Progress) InfoTo(logger Logger) {
	logger.Infof(
		"%v of %v (%v/s; %v elapsed, %v remaining)\n",
		FormatBytes(p.TotalBytes),
		FormatBytes(p.ContentLength),
		FormatBytes(p.EstimatedBps()),
		FormatNanos(p.NsElapsed),
		FormatNanos(p.estNsRemaining()),
	)
}

func (p *Progress) EstimatedBps() int64 {
	var estBytesPerSecond int64
	if p.estimatedBps == nil {
		estBps := (nsPerSecondFloat64 * float64(p.TotalBytes)) / (float64(p.NsElapsed))
		estBytesPerSecond = int64(math.Round(estBps))
		p.estimatedBps = &estBytesPerSecond
	} else {
		estBytesPerSecond = *p.estimatedBps
	}
	return estBytesPerSecond
}

// ------------------------------------------------------------
// Unexported functions

func (p *Progress) fmtTotalBytes() string {
	totalBytes := uint64(p.TotalBytes)
	return bytefmt.ByteSize(totalBytes)
}

func (p *Progress) fmtContentLength() string {
	contentLength := uint64(p.ContentLength)
	return bytefmt.ByteSize(contentLength)
}

func (p *Progress) fmtEstBps() string {
	estBps := uint64(p.EstimatedBps())
	return bytefmt.ByteSize(estBps)
}

func (p *Progress) estNsRemaining() int64 {
	var nsRemaining float64
	estBytesPerSecond := p.EstimatedBps()
	bytesRemaining := p.ContentLength - p.TotalBytes
	nsRemaining = nsPerSecondFloat64 * float64(bytesRemaining) / float64(estBytesPerSecond)
	return int64(math.Round(nsRemaining))
}
