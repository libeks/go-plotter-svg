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
	Reverse() LineLike
	Bisect(t float64) (Path, Path)
}

// implemented by LineChunk, QuadraticBezierChunk, CubicBezierChunk, LineGapChunk
type PathChunk interface {
	PathXML() string // returns the XML of this chunk, assuming that the starting point is the endpoint of the
	// previous chunk

	Length() float64
	Startpoint() primitives.Point
	Endpoint() primitives.Point

	ControlLines() string
	OffsetLeft(distance float64) PathChunk
	Reverse() PathChunk
	Bisect(t float64) (PathChunk, PathChunk)
}
