package streaming

import (
	"math"
	"time"

	"code.cloudfoundry.org/bytefmt"

	"github.com/dmolesUC3/cos/internal/logging"
)

const nsPerSecond = int64(time.Second)
const nsPerSecondFloat64 = float64(time.Second)

type Progress struct {
	NsElapsed     int64
	TotalBytes    int64
	ContentLength int64
	estimatedBps  *int64
}

func (p *Progress) InfoTo(logger logging.Logger) {
	logger.Infof(
		"%v of %v (%v/s; %v elapsed, %v remaining)\n",
		logging.FormatBytes(p.TotalBytes),
		logging.FormatBytes(p.ContentLength),
		logging.FormatBytes(p.EstimatedBps()),
		logging.FormatNanos(p.NsElapsed),
		logging.FormatNanos(p.estNsRemaining()),
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

