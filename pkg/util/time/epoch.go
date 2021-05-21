package time

import (
	"time"
)

// Epoch returns epoch as seconds with UTC tz
func Epoch() int64 {
	return time.Now().UTC().Unix()
}

// EpochMillis returns epoch as millis with UTC tz
func EpochMillis() int64 {
	return time.Now().UTC().UnixNano() / 1_000_000
}