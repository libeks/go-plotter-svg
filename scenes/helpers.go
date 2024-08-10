package scenes

import (
	"time"

	"github.com/libeks/go-plotter-svg/lines"
)

func segmentsToLineLikes(segments []lines.LineSegment) []lines.LineLike {
	linelikes := make([]lines.LineLike, len(segments))
	for i, seg := range segments {
		linelikes[i] = seg
	}
	return linelikes
}

// return the time spent moving the pen up and down for this many segments
func upDownEstimate(n int) time.Duration {
	oneUpAndDownEstimate := time.Millisecond * 400
	return oneUpAndDownEstimate * time.Duration(n)
}
