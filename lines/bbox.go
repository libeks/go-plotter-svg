package lines

import (
	"github.com/libeks/go-plotter-svg/primitives"
)

// Lines returns a linelike representing this box.
func LinesFromBBox(b primitives.BBox) []LineLike {
	path := NewPath(b.NWCorner())
	// find the starting point - extreme point of box in direction perpendicular to

	path = path.AddPathChunk(LineChunk{Start: b.NWCorner(), End: b.NECorner()})
	path = path.AddPathChunk(LineChunk{Start: b.NECorner(), End: b.SECorner()})
	path = path.AddPathChunk(LineChunk{Start: b.SECorner(), End: b.SWCorner()})
	path = path.AddPathChunk(LineChunk{Start: b.SWCorner(), End: b.NWCorner()})

	return []LineLike{
		path,
	}
}
