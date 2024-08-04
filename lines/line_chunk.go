package lines

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/primitives"
)

type LineChunk struct {
	Start primitives.Point
	End   primitives.Point
}

func (c LineChunk) String() string {
	return fmt.Sprintf("LineChunk %s %s", c.Start, c.End)
}

func (c LineChunk) PathXML() string {
	return fmt.Sprintf("L %.1f %.1f", c.End.X, c.End.Y)
}

func (c LineChunk) Length() float64 {
	return c.Start.Subtract(c.End).Len()
}

func (l LineChunk) ControlLines() string {
	return fmt.Sprintf("L %.1f %.1f", l.End.X, l.End.Y)
}

func (c LineChunk) Endpoint() primitives.Point {
	return c.End
}

func (l LineChunk) At(t float64) primitives.Point {
	return l.Start.Add(l.End.Subtract(l.Start).Mult(t))
}

func (l LineChunk) Startpoint() primitives.Point {
	return l.Start
}

func (c LineChunk) OffsetLeft(distance float64) PathChunk {
	v := c.End.Subtract(c.Start).Perp().Unit().Mult(distance)
	return LineChunk{Start: c.Start.Add(v), End: c.End.Add(v)}
}

func (c LineChunk) Reverse() PathChunk {
	return LineChunk{
		Start: c.End,
		End:   c.Start,
	}
}

func (c LineChunk) Bisect(t float64) (PathChunk, PathChunk) {
	midpoint := c.At(t)
	return LineChunk{
			Start: c.Start,
			End:   midpoint,
		}, LineChunk{
			Start: midpoint,
			End:   c.End,
		}
}
