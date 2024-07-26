package lines

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/primitives"
)

const (
	lengthEstimateAccuracy = 0.1
)

type LengthEstimator interface {
	At(float64) primitives.Point
	BBox() primitives.BBox
}

type QuadraticBezierChunk struct {
	Start primitives.Point
	P1    primitives.Point
	End   primitives.Point
}

func (c QuadraticBezierChunk) XMLChunk() string {
	return fmt.Sprintf("Q %.1f %.1f, %.1f %.1f", c.P1.X, c.P1.Y, c.End.X, c.End.Y)
}

// t is in [0,1]
func (c QuadraticBezierChunk) At(t float64) primitives.Point {
	a1 := LineStartEndChunk{Start: c.Start, End: c.P1}.At(t)
	a2 := LineStartEndChunk{Start: c.P1, End: c.End}.At(t)
	return LineStartEndChunk{Start: a1, End: a2}.At(t)
}

func (c QuadraticBezierChunk) BBox() primitives.BBox {
	return primitives.BBoxAroundPoints(c.Start, c.P1, c.End)
}

func (c QuadraticBezierChunk) Length(start primitives.Point) float64 {
	return estimateLength(c, lengthEstimateAccuracy)
}

func (c QuadraticBezierChunk) Bisect(t float64) (QuadraticBezierChunk, QuadraticBezierChunk) {
	p0 := c.Start
	p1 := c.End
	a1 := LineStartEndChunk{Start: c.Start, End: c.P1}.At(t)
	a2 := LineStartEndChunk{Start: c.P1, End: c.End}.At(t)
	b1 := LineStartEndChunk{Start: a1, End: a2}.At(t)
	return QuadraticBezierChunk{Start: p0, P1: a1, End: b1}, QuadraticBezierChunk{Start: b1, P1: a2, End: p1}
}

func (c QuadraticBezierChunk) Endpoint() primitives.Point {
	return c.End
}

func (c QuadraticBezierChunk) Guides() string {
	return fmt.Sprintf("M %.1f %.1f L %.1f %.1f L %.1f %.1f", c.Start.X, c.Start.Y, c.P1.X, c.P1.Y, c.End.X, c.End.Y)
}

type CubicBezierChunk struct {
	Start primitives.Point
	P1    primitives.Point
	P2    primitives.Point
	End   primitives.Point
}

func (c CubicBezierChunk) String() string {
	return fmt.Sprintf("Cubic Bezier with pts (%s %s %s %s)", c.Start, c.P1, c.P2, c.End)
}

func (c CubicBezierChunk) XMLChunk() string {
	return fmt.Sprintf("C %.1f %.1f, %.1f %.1f, %.1f %.1f", c.P1.X, c.P1.Y, c.P2.X, c.P2.Y, c.End.X, c.End.Y)
}

// t is in [0,1]
func (c CubicBezierChunk) At(t float64) primitives.Point {
	a1 := LineStartEndChunk{Start: c.Start, End: c.P1}.At(t)
	a2 := LineStartEndChunk{Start: c.P1, End: c.P2}.At(t)
	a3 := LineStartEndChunk{Start: c.P2, End: c.End}.At(t)
	b1 := LineStartEndChunk{Start: a1, End: a2}.At(t)
	b2 := LineStartEndChunk{Start: a2, End: a3}.At(t)
	return LineStartEndChunk{Start: b1, End: b2}.At(t)
}

func (c CubicBezierChunk) Bisect(t float64) (CubicBezierChunk, CubicBezierChunk) {
	p0 := c.Start
	p1 := c.End
	a1 := LineStartEndChunk{Start: c.Start, End: c.P1}.At(t)
	a2 := LineStartEndChunk{Start: c.P1, End: c.P2}.At(t)
	a3 := LineStartEndChunk{Start: c.P2, End: c.End}.At(t)
	b1 := LineStartEndChunk{Start: a1, End: a2}.At(t)
	b2 := LineStartEndChunk{Start: a2, End: a3}.At(t)
	cc := LineStartEndChunk{Start: b1, End: b2}.At(t)
	return CubicBezierChunk{Start: p0, P1: a1, P2: b1, End: cc}, CubicBezierChunk{Start: cc, P1: b2, P2: a3, End: p1}
}

func (c CubicBezierChunk) BBox() primitives.BBox {
	return primitives.BBoxAroundPoints(c.Start, c.P1, c.P2, c.End)
}

func (c CubicBezierChunk) Length(start primitives.Point) float64 {
	return estimateLength(c, lengthEstimateAccuracy)
}

func (c CubicBezierChunk) Guides() string {
	return fmt.Sprintf("M %.1f %.1f L %.1f %.1f L %.1f %.1f L %.1f %.1f", c.Start.X, c.Start.Y, c.P1.X, c.P1.Y, c.P2.X, c.P2.Y, c.End.X, c.End.Y)
}

func (c CubicBezierChunk) Endpoint() primitives.Point {
	return c.End
}
