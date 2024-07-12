package lines

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/primitives"
)

type LineGapChunk struct {
	Start        primitives.Point
	GapSizeRatio float64 // the relative ratio of the length of the line to keep empty in the middle
	End          primitives.Point
}

func (c LineGapChunk) XMLChunk() string {
	v := c.End.Subtract(c.Start)
	end1 := c.Start.Add(v.Mult(0.5 - c.GapSizeRatio/2))
	start2 := c.Start.Add(v.Mult(0.5 + c.GapSizeRatio/2))
	return fmt.Sprintf("L %.1f %.1f M %.1f %.1f L %.1f %.1f", end1.X, end1.Y, start2.X, start2.Y, c.End.X, c.End.Y)
}

func (c LineGapChunk) Length(start primitives.Point) float64 {
	return start.Subtract(c.End).Len()
}

func (c LineGapChunk) Endpoint() primitives.Point {
	return c.End
}
