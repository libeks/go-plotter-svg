package lines

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/primitives"
)

type CircleArcChunk struct {
	Radius      float64
	End         primitives.Point
	IsLong      bool
	IsClockwise bool
}

func (c CircleArcChunk) XMLChunk() string {
	long := 0
	if c.IsLong {
		long = 1
	}
	clockwise := 0
	if c.IsClockwise {
		clockwise = 1
	}
	return fmt.Sprintf("A %.1f %.1f 0 %d %d %.1f %.1f", c.Radius, c.Radius, long, clockwise, c.End.X, c.End.Y)
}

func (c CircleArcChunk) Length(start primitives.Point) float64 {
	dv := start.Subtract(c.End)
	distance := dv.Len()
	angle := math.Asin(distance / (2 * c.Radius))
	if c.IsLong {
		return 2 * (math.Pi - angle) * c.Radius
	}
	return 2 * angle * c.Radius
}

func (c CircleArcChunk) Endpoint() primitives.Point {
	return c.End
}
