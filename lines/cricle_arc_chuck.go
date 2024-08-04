package lines

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/primitives"
)

func CircleArcChunk(center primitives.Point, radius float64, startRad, endRad float64, isClockwise bool) circleArcChunk {
	return circleArcChunk{
		radius:      radius,
		center:      center,
		isClockwise: isClockwise,
		startRad:    startRad,
		endRad:      endRad,
	}
}

type circleArcChunk struct {
	radius      float64
	center      primitives.Point
	startRad    float64
	endRad      float64
	isClockwise bool
}

func (c circleArcChunk) String() string {
	return fmt.Sprintf("CircleArcChunk: center %s, radius %.1f with start %.1f, end %.1f", c.center, c.radius, c.startRad, c.endRad)
}

func (c circleArcChunk) IsLong() bool {
	angle := c.Angle()
	// fmt.Printf("angle is %.1f\n", angle/math.Pi)
	return angle > math.Pi || angle < -math.Pi
}

func angleMath(angle float64) float64 {
	if angle < 2*math.Pi && angle > 0 {
		return angle
	}
	if angle >= 2*math.Pi {
		return angle - 2*math.Pi
	}
	return angle + 2*math.Pi

}

func (c circleArcChunk) Angle() float64 {
	if c.isClockwise {
		return angleMath(c.endRad - c.startRad)
	}
	return angleMath(c.startRad - c.endRad)
}

func (c circleArcChunk) PathXML() string {
	long := 0
	if c.IsLong() {
		long = 1
	}
	clockwise := 0
	if c.isClockwise {
		clockwise = 1
	}
	endpoint := c.Endpoint()
	return fmt.Sprintf("A %.1f %.1f 0 %d %d %.1f %.1f", c.radius, c.radius, long, clockwise, endpoint.X, endpoint.Y)
}

func (c circleArcChunk) Length(start primitives.Point) float64 {
	angle := c.Angle()
	if c.IsLong() {
		return 2 * (math.Pi - angle) * c.radius
	}
	return 2 * angle * c.radius
}

func (c circleArcChunk) Endpoint() primitives.Point {
	return c.center.Add(primitives.UnitRight.RotateCCW(c.endRad).Mult(c.radius))
}

func (c circleArcChunk) Startpoint() primitives.Point {
	return c.center.Add(primitives.UnitRight.RotateCCW(c.startRad).Mult(c.radius))
}

func (c circleArcChunk) ControlLines() string {
	end := c.Endpoint()
	return fmt.Sprintf("L %.1f %.1f", end.X, end.Y)
}

func (c circleArcChunk) OffsetLeft(distance float64) PathChunk {
	fmt.Printf("offsetting left %s\n", c)
	// left is counterClockwise
	if c.isClockwise {
		distance *= -1
	}
	newRadius := c.radius + distance
	return circleArcChunk{
		radius:      newRadius,
		center:      c.center,
		isClockwise: c.isClockwise,
		startRad:    c.startRad,
		endRad:      c.endRad,
	}
}

func (c circleArcChunk) Reverse() PathChunk {
	return circleArcChunk{
		radius:      c.radius,
		center:      c.center,
		startRad:    c.endRad,
		endRad:      c.startRad,
		isClockwise: !c.isClockwise,
	}
}

type CircleArcChunkLegacy struct {
	// TODO: refactor to not require all of these fields
	Radius      float64
	Center      primitives.Point
	Start       primitives.Point // could be replaced with startT
	End         primitives.Point // could be replaced with endT
	IsLong      bool             // could be removed altogether
	IsClockwise bool             // could be removed altogether
}

func (c CircleArcChunkLegacy) String() string {
	return fmt.Sprintf("CircleArcChunk: center %s, radius %.1f with start %s, end %s", c.Center, c.Radius, c.Start, c.End)
}

func (c CircleArcChunkLegacy) PathXML() string {
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

func (c CircleArcChunkLegacy) Length(start primitives.Point) float64 {
	dv := start.Subtract(c.End)
	distance := dv.Len()
	angle := math.Asin(distance / (2 * c.Radius))
	if c.IsLong {
		return 2 * (math.Pi - angle) * c.Radius
	}
	return 2 * angle * c.Radius
}

func (c CircleArcChunkLegacy) Endpoint() primitives.Point {
	return c.End
}

func (c CircleArcChunkLegacy) Startpoint() primitives.Point {
	return c.Start
}

func (c CircleArcChunkLegacy) ControlLines() string {
	return fmt.Sprintf("L %.1f %.1f", c.End.X, c.End.Y)
}

func (c CircleArcChunkLegacy) OffsetLeft(distance float64) PathChunk {
	fmt.Printf("offsetting left %s\n", c)
	// left is counterClockwise
	if c.IsClockwise {
		distance *= -1
	}
	newRadius := c.Radius + distance
	radiusRatio := newRadius / c.Radius
	sv := c.Start.Subtract(c.Center).Mult(radiusRatio)
	ev := c.End.Subtract(c.Center).Mult(radiusRatio)
	return CircleArcChunkLegacy{
		Radius:      newRadius,
		Center:      c.Center,
		Start:       c.Center.Add(sv),
		End:         c.Center.Add(ev), // need to know the center here...
		IsLong:      c.IsLong,
		IsClockwise: c.IsClockwise,
	}
}

func (c CircleArcChunkLegacy) Reverse() PathChunk {
	return CircleArcChunkLegacy{
		Radius:      c.Radius,
		Center:      c.Center,
		Start:       c.End,
		End:         c.Start,
		IsLong:      c.IsLong,
		IsClockwise: !c.IsClockwise,
	}
}
