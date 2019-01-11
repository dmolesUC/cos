package logging

import (
	"fmt"
	"time"
)

func DetailProgress(logger Logger, totalBytes int64, contentLength int64, estKps float64, nsElapsed int64, nsRemaining int64) {
	elapsedStr := formatNanos(nsElapsed)
	remainingStr := formatNanos(nsRemaining)
	logger.Detailf(
		"read %d of %d bytes (%0.f KiB/s; %v elapsed, %v remaining)\n",
		totalBytes, contentLength, estKps, elapsedStr, remainingStr,
	)
}

func formatNanos(ns int64) string {
	hours := ns / int64(time.Hour)
	remainder := ns % int64(time.Hour)
	minutes := remainder / int64(time.Minute)
	remainder = ns % int64(time.Minute)
	seconds := remainder / int64(time.Second)
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
