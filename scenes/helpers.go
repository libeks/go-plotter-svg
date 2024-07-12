package scenes

import (
	"github.com/libeks/go-plotter-svg/lines"
)

func segmentsToLineLikes(segments []lines.LineSegment) []lines.LineLike {
	linelikes := make([]lines.LineLike, len(segments))
	for i, seg := range segments {
		linelikes[i] = seg
	}
	return linelikes
}
