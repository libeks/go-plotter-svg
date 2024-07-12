package objects

import (
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

type Object interface {
	Inside(p primitives.Point) bool
	IntersectTs(line lines.Line) []float64     // return the t-values of the line intersecting with this object
	IntersectCircleTs(circle Circle) []float64 // return the angle-t values of the circle intersecting with this object
}
