package lines

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/primitives"
)

type LineChunk struct {
	End primitives.Point
}

func (c LineChunk) XMLChunk() string {
	return fmt.Sprintf("L %.1f %.1f", c.End.X, c.End.Y)
}

func (c LineChunk) Length(start primitives.Point) float64 {
	return start.Subtract(c.End).Len()
}

func (c LineChunk) Endpoint() primitives.Point {
	return c.End
}
