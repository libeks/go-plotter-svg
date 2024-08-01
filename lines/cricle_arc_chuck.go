package lines

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/primitives"
)

type CircleArcChunk struct {
	// TODO: refactor to not require all of these fields
	Radius      float64
	Center      primitives.Point
	Start       primitives.Point // could be replaced with startT
	End         primitives.Point // could be replaced with endT
	IsLong      bool             // could be removed altogether
	IsClockwise bool             // could be removed altogether
}

func (c CircleArcChunk) String() string {
	return fmt.Sprintf("CircleArcChunk: center %s, radius %.1f with start %s, end %s", c.Center, c.Radius, c.Start, c.End)
}

func (c CircleArcChunk) PathXML() string {
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

func (c CircleArcChunk) Startpoint() primitives.Point {
	return c.Start
}

func (c CircleArcChunk) ControlLines() string {
	return fmt.Sprintf("L %.1f %.1f", c.End.X, c.End.Y)
}

func (c CircleArcChunk) OffsetLeft(distance float64) PathChunk {
	fmt.Printf("offsetting left %s\n", c)
	// left is counterClockwise
	if c.IsClockwise {
		distance *= -1
	}
	newRadius := c.Radius + distance
	radiusRatio := newRadius / c.Radius
	sv := c.Start.Subtract(c.Center).Mult(radiusRatio)
	ev := c.End.Subtract(c.Center).Mult(radiusRatio)
	return CircleArcChunk{
		Radius:      newRadius,
		Center:      c.Center,
		Start:       c.Center.Add(sv),
		End:         c.Center.Add(ev), // need to know the center here...
		IsLong:      c.IsLong,
		IsClockwise: c.IsClockwise,
	}
}
