package lines

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/maths"
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

func (c circleArcChunk) Length() float64 {
	angle := c.Angle()
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

func (c circleArcChunk) At(t float64) primitives.Point {
	var angle float64
	if c.isClockwise {
		angle = maths.Interpolate(c.endRad, c.startRad, t)
	} else {
		angle = maths.Interpolate(c.startRad, c.endRad+math.Pi*2, t)
	}
	return c.center.Add(primitives.UnitRight.RotateCCW(angle).Mult(c.radius))
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

func (c circleArcChunk) Translate(v primitives.Vector) PathChunk {
	c.center = c.center.Add(v)
	return c
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

func (c circleArcChunk) Bisect(t float64) (PathChunk, PathChunk) {
	var midT float64
	if c.isClockwise {
		midT = maths.Interpolate(c.startRad, c.endRad, t)
	} else {
		midT = maths.Interpolate(c.endRad, c.startRad, t)
	}
	return circleArcChunk{
			radius:      c.radius,
			center:      c.center,
			startRad:    c.startRad,
			endRad:      midT,
			isClockwise: c.isClockwise,
		}, circleArcChunk{
			radius:      c.radius,
			center:      c.center,
			startRad:    midT,
			endRad:      c.endRad,
			isClockwise: c.isClockwise,
		}
}
