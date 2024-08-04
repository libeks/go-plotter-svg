package lines

import (
	"math"

	"github.com/libeks/go-plotter-svg/primitives"
)

func FullCircle(center primitives.Point, radius float64) LineLike {
	return NewPath(center.Add(primitives.UnitRight.Mult(radius))).AddPathChunk(CircleArcChunk(center, radius, 0, math.Pi, true)).AddPathChunk(CircleArcChunk(center, radius, math.Pi, 0, true))
}
