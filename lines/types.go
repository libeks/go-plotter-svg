package lines

import (
	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/primitives"
)

// implemented by LineSegment, Path, CircleArc, Circle
type LineLike interface {
	XML(color, width string) xmlwriter.Elem
	String() string
	IsEmpty() bool
	Len() float64
	Start() primitives.Point
	End() primitives.Point
	ControlLineXML(color, width string) xmlwriter.Elem
	OffsetLeft(distance float64) LineLike
}

// implemented by LineChunk, QuadraticBezierChunk, CubicBezierChunk, LineGapChunk
type PathChunk interface {
	XMLChunk() string
	Length(start primitives.Point) float64
	Endpoint() primitives.Point
	ControlLines() string
	Startpoint() primitives.Point
	OffsetLeft(distance float64) PathChunk
}
