package lines

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/primitives"
)

type QuadraticBezierChunk struct {
	Start primitives.Point
	P1    primitives.Point
	End   primitives.Point
}

func (c QuadraticBezierChunk) XMLChunk() string {
	return fmt.Sprintf("Q %.1f %.1f, %.1f %.1f", c.P1.X, c.P1.Y, c.End.X, c.End.Y)
}

func (c QuadraticBezierChunk) Length(start primitives.Point) float64 {
	// TODO: estimate curve length
	return start.Subtract(c.End).Len()
}

func (c QuadraticBezierChunk) Endpoint() primitives.Point {
	return c.End
}

type CubicBezierChunk struct {
	Start primitives.Point
	P1    primitives.Point
	P2    primitives.Point
	End   primitives.Point
}

func (c CubicBezierChunk) XMLChunk() string {
	return fmt.Sprintf("C %.1f %.1f, %.1f %.1f, %.1f %.1f", c.P1.X, c.P1.Y, c.P2.X, c.P2.Y, c.End.X, c.End.Y)
}

func (c CubicBezierChunk) Length(start primitives.Point) float64 {
	// TODO: estimate curve length
	return start.Subtract(c.End).Len()
}

func (c CubicBezierChunk) Endpoint() primitives.Point {
	return c.End
}
