package scenes

import (
	"time"
)

// return the time spent moving the pen up and down for this many segments
func upDownEstimate(n int) time.Duration {
	oneUpAndDownEstimate := time.Millisecond * 400
	return oneUpAndDownEstimate * time.Duration(n)
}
