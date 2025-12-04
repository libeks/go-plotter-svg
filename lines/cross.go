package lines

import (
	"math"

	"github.com/libeks/go-plotter-svg/primitives"
)

func Cross(center primitives.Point, size float64) []LineLike {
	upperLeft := center.Add(primitives.Vector{X: -math.Sqrt2, Y: -math.Sqrt2}.Mult(size))
	lowerRight := center.Add(primitives.Vector{X: math.Sqrt2, Y: math.Sqrt2}.Mult(size))
	upperRight := center.Add(primitives.Vector{X: math.Sqrt2, Y: -math.Sqrt2}.Mult(size))
	lowerLeft := center.Add(primitives.Vector{X: -math.Sqrt2, Y: math.Sqrt2}.Mult(size))
	return []LineLike{
		NewPath(upperLeft).AddPathChunk(
			LineChunk{Start: upperLeft, End: lowerRight}),
		NewPath(upperRight).AddPathChunk(
			LineChunk{Start: upperRight, End: lowerLeft}),
	}
}
