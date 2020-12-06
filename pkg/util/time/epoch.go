package time

import (
	"time"
)

// Epoch returns epoch as seconds with UTC
func Epoch() int64 {
	return time.Now().UTC().Unix()
}