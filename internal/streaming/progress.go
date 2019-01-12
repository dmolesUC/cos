package streaming

import (
	"fmt"
	"time"

	"github.com/dmolesUC3/cos/internal/logging"
)

const nsPerSecond = int64(time.Second)
const nsPerMinute = int64(time.Minute)
const nsPerHour = int64(time.Hour)
const nsPerSecondFloat64 = float64(time.Second)

type Progress struct {
	NsElapsed     int64
	TotalBytes    int64
	ContentLength int64
	estimatedBytesPerSecond *float64
}

func (p *Progress) DetailTo(logger logging.Logger) {
	logger.Detailf(
		"read %d of %d bytes (%0.f KiB/s; %v elapsed, %v remaining)\n",
		p.TotalBytes, p.ContentLength, p.EstimatedKibPerSecond(), p.NsElapsedStr(), p.NsRemainingStr(),
	)
}

func (p *Progress) EstimatedKibPerSecond() float64 {
	estBytesPerSecond := p.EstimatedBytesPerSecond()
	estKibPerSecond := estBytesPerSecond / float64(1024)
	return estKibPerSecond
}

func (p *Progress) NsElapsedStr() string {
	return formatNanos(p.NsElapsed)
}

func (p *Progress) NsRemainingStr() string {
	return formatNanos(p.EstimatedNsRemaining())
}

func (p *Progress) EstimatedNsRemaining() int64 {
	estBytesPerSecond := p.EstimatedBytesPerSecond()
	bytesRemaining := p.ContentLength - p.TotalBytes
	nsRemaining := int64(float64(bytesRemaining) / estBytesPerSecond)
	return nsRemaining
}

func (p *Progress) EstimatedBytesPerSecond() float64 {
	var estBytesPerSecond float64
	if p.estimatedBytesPerSecond == nil {
		estBytesPerSecond = (nsPerSecondFloat64 * float64(p.TotalBytes)) / (float64(p.NsElapsed))
		p.estimatedBytesPerSecond = &estBytesPerSecond
	} else {
		estBytesPerSecond = *p.estimatedBytesPerSecond
	}
	return estBytesPerSecond
}

func formatNanos(ns int64) string {
	hours := ns / nsPerHour
	remainder := ns % nsPerHour
	minutes := remainder / nsPerMinute
	remainder = ns % nsPerMinute
	seconds := remainder / nsPerSecond
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
