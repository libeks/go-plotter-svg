package objects

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/primitives"
)

type Polygon struct {
	Points []primitives.Point
}

func (p Polygon) String() string {
	return fmt.Sprintf("Polygon (%v)", p.Points)
}

func (p Polygon) Inside(pt primitives.Point) bool {
	// TODO: take point as input
	// compute the winding angle from the point
	totalAngle := 0.0
	for i, p1 := range p.Points {
		j := (i + 1) % len(p.Points)
		p2 := p.Points[j]

		p2Angle := p2.Subtract(pt).Atan()
		p1Angle := p1.Subtract(pt).Atan()
		angle := maths.AngleDifference(p2Angle, p1Angle)
		totalAngle += angle
	}
	if math.Abs(totalAngle-math.Pi) < 0.01 {
		return true
	}
	if math.Abs(totalAngle+math.Pi) < 0.01 {
		return true
	}
	if math.Abs(totalAngle) < 0.01 {
		return false
	}
	panic(fmt.Errorf("not sure what to do with winding angle %.3f", totalAngle))
}

func (p Polygon) EdgeLines() []lines.LineSegment {
	segments := []lines.LineSegment{}
	for i, p1 := range p.Points {
		j := (i + 1) % len(p.Points)
		p2 := p.Points[j]
		segments = append(segments, lines.LineSegment{p1, p2})
	}
	return segments
}

func (p Polygon) IntersectTs(line lines.Line) []float64 {
	ts := []float64{}
	for _, segment := range p.EdgeLines() {
		if t := line.IntersectLineSegmentT(segment); t != nil {
			ts = append(ts, *t)
		}
	}
	return ts
}

// should return circle t-values
func (p Polygon) IntersectCircleTs(circle Circle) []float64 {
	ts := []float64{}
	for _, segment := range p.EdgeLines() {
		t := circle.IntersectLineSegmentT(segment)
		ts = append(ts, t...)
	}
	return ts
}
