package lines

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/primitives"
)

type LineChunk struct {
	Start primitives.Point
	End   primitives.Point
}

func (c LineChunk) XMLChunk() string {
	return fmt.Sprintf("L %.1f %.1f", c.End.X, c.End.Y)
}

func (c LineChunk) Length(start primitives.Point) float64 {
	return start.Subtract(c.End).Len()
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
