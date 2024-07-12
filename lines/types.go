package lines

import (
	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/shabbyrobe/xmlwriter"
)

// implemented by LineSegment, Path, CircleArc, Circle
type LineLike interface {
	XML(color, width string) xmlwriter.Elem
	String() string
	IsEmpty() bool
	Len() float64
	Start() primitives.Point
	End() primitives.Point
}

type PathChunk interface {
	XMLChunk() string
	Length(start primitives.Point) float64
	Endpoint() primitives.Point
}
