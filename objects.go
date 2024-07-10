package main

import (
	"fmt"
	"math"
	"slices"

	"github.com/shabbyrobe/xmlwriter"
)

type Object interface {
	Inside(p Point) bool
	IntersectTs(line Line) []float64           // return the t-values of the line intersecting with this object
	IntersectCircleTs(circle Circle) []float64 // return the angle-t values of the circle intersecting with this object
}

func ClipLineToObject(line Line, obj Object) []LineSegment {
	ts := obj.IntersectTs(line)
	if len(ts) == 0 {
		return nil
	}
	if len(ts) == 1 {
		panic(fmt.Errorf("not sure what to do with only one intersection %v", ts))
	}
	slices.Sort(ts)
	segments := []LineSegment{}
	for i, t1 := range ts {
		if i == len(ts)-1 {
			break
		}
		t2 := ts[(i + 1)]
		midT := average(t1, t2)
		midPoint := line.At(midT)
		if obj.Inside(midPoint) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			seg := LineSegment{p1, p2}
			segments = append(segments, seg)
		}
	}
	return segments
}

func ClipLineSegmentToObject(ls LineSegment, obj Object) []LineSegment {
	line := ls.Line()
	ts := append(obj.IntersectTs(line), 0.0, 1.0)
	ts = filterToRange(ts, 0.0, 1.0)
	slices.Sort(ts)
	segments := []LineSegment{}
	for i, t1 := range ts {
		t2 := ts[(i+1)%len(ts)]
		midT := average(t1, t2)
		midPoint := line.At(midT)
		if obj.Inside(midPoint) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			segments = append(segments, LineSegment{
				p1, p2,
			})
		}
	}
	return segments
}

func ClipCircleToObject(c Circle, obj Object) []LineLike {
	ts := obj.IntersectCircleTs(c)
	if len(ts) == 0 {
		return []LineLike{c}
	}
	slices.Sort(ts)
	segments := []LineLike{}
	for i, t1 := range ts {
		t2 := ts[(i+1)%len(ts)]
		midT := average(t1, t2)
		midPoint := c.At(midT)
		if obj.Inside(midPoint) {
			segments = append(segments, CircleArc(c, t1, t2))
		}
	}
	return segments
}

func filterToRange(vals []float64, min, max float64) []float64 {
	rets := []float64{}
	for _, val := range vals {
		if val <= max && val >= min {
			rets = append(rets, val)
		}
	}
	return rets
}

func NewCompositeWithWithout(with []Object, without []Object) CompositeObject {
	return CompositeObject{}.With(with...).Without(without...)
}
func NewComposite() CompositeObject {
	return CompositeObject{}
}

type CompositeObject struct {
	positive []Object
	negative []Object
}

func (o CompositeObject) With(obj ...Object) CompositeObject {
	return CompositeObject{
		positive: append(o.positive, obj...),
		negative: o.negative,
	}
}

func (o CompositeObject) Without(obj ...Object) CompositeObject {
	return CompositeObject{
		positive: o.positive,
		negative: append(o.negative, obj...),
	}
}

func (o CompositeObject) Inside(p Point) bool {
	inside := false
	for _, pos := range o.positive {
		if pos.Inside(p) {
			inside = true
			break
		}
	}
	if !inside {
		return false
	}
	for _, neg := range o.negative {
		if neg.Inside(p) {
			return false
		}
	}
	return true
}

func (o CompositeObject) IntersectTs(line Line) []float64 {
	ts := []float64{}
	for _, obj := range append(o.positive, o.negative...) {
		ts = append(ts, obj.IntersectTs(line)...)
	}
	return ts
}

func (o CompositeObject) IntersectCircleTs(circle Circle) []float64 {
	ts := []float64{}
	for _, obj := range append(o.positive, o.negative...) {
		ts = append(ts, obj.IntersectCircleTs(circle)...)
	}
	return ts
}

func angleDifference(a2, a1 float64) float64 {
	diff := a2 - a1
	if diff < -math.Pi {
		diff = diff + math.Pi
	}
	if diff > math.Pi {
		diff = diff - math.Pi
	}
	return diff
}

type Polygon struct {
	points []Point
}

func (p Polygon) String() string {
	return fmt.Sprintf("Polygon (%v)", p.points)
}

func (p Polygon) Inside(pt Point) bool {
	// TODO: take point as input
	// compute the winding angle from the point
	totalAngle := 0.0
	for i, p1 := range p.points {
		j := (i + 1) % len(p.points)
		p2 := p.points[j]

		p2Angle := p2.Subtract(pt).Atan()
		p1Angle := p1.Subtract(pt).Atan()
		angle := angleDifference(p2Angle, p1Angle)
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

func (p Polygon) EdgeLines() []LineSegment {
	segments := []LineSegment{}
	for i, p1 := range p.points {
		j := (i + 1) % len(p.points)
		p2 := p.points[j]
		segments = append(segments, LineSegment{p1, p2})
	}
	return segments
}

func (p Polygon) IntersectTs(line Line) []float64 {
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

func average(a, b float64) float64 {
	return (a + b) / 2
}

type Circle struct {
	center Point
	radius float64
}

func (c Circle) Len() float64 {
	return math.Pi * 2 * c.radius
}

func (c Circle) Start() Point {
	// usually the plotter plots circles starting and ending at the left extreme point
	return c.center.Add(Vector{-1, 0}.Mult(c.radius))
}

func (c Circle) End() Point {
	// usually the plotter plots circles starting and ending at the left extreme point
	return c.center.Add(Vector{-1, 0}.Mult(c.radius))
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle @%s, r:%.1f", c.center, c.radius)
}

func (c Circle) Inside(p Point) bool {
	distance := p.Subtract(c.center).Len()
	return distance <= c.radius
}

func (c Circle) At(t float64) Point {
	return c.center.Add(Vector{c.radius, 0}.RotateCCW(t))
}

func (c Circle) IsEmpty() bool {
	return c.radius == 0
}

// return the line ts when intersecting with the circle
func (c Circle) IntersectTs(line Line) []float64 {
	w := line.p.Subtract(c.center)
	A := line.v.x*line.v.x + line.v.y*line.v.y
	B := 2 * (line.v.x*w.x + line.v.y*w.y)
	C := w.x*w.x + w.y*w.y - c.radius*c.radius
	ts := quadratic(A, B, C)
	if len(ts) < 2 {
		return nil
	}
	return ts
}

// should return c2 t-values
// IntersectCircleTs returns the angles respective to c2 at which it intersects c1
func (c1 Circle) IntersectCircleTs(c2 Circle) []float64 {
	// distance between the two centers
	dVect := c1.center.Subtract(c2.center)
	d := dVect.Len()
	r1 := c1.radius
	r2 := c2.radius
	// circles too far to intersect
	if d > r1+r2 {
		return nil
	}
	// circles are colinear, return no intersection, even if they are the same circle
	if d == 0 {
		return nil
	}
	// x is distance from c1 at which the chord lies
	x := (d*d + r1*r1 - r2*r2) / (2 * d)
	// y is the half-length of the chord
	y := math.Sqrt(r1*r1 - x*x)
	dUnit := dVect.Unit()
	xVect := dUnit.Mult(x)
	yVect := dUnit.Perp().Mult(y)
	v1 := xVect.Add(yVect)
	v2 := xVect.Add(yVect.Mult(-1))
	p1 := c1.center.Add(v1)
	p2 := c1.center.Add(v2)
	w1 := p1.Subtract(c2.center)
	w2 := p2.Subtract(c2.center)
	if math.Abs(v1.Len()-r1) > 0.001 {
		panic(fmt.Errorf("circle %s intersection vector has wrong length %.1f, want %.1f", c1, v1.Len(), r1))
	}
	if math.Abs(v2.Len()-r1) > 0.001 {
		panic(fmt.Errorf("circle %s intersection vector has wrong length %.1f, want %.1f", c1, v1.Len(), r1))
	}
	if math.Abs(w1.Len()-r2) > 0.001 {
		panic(fmt.Errorf("circle %s intersection vector has wrong length %.1f, want %.1f", c2, w1.Len(), r2))
	}
	if math.Abs(w2.Len()-r2) > 0.001 {
		panic(fmt.Errorf("circle %s intersection vector has wrong length %.1f, want %.1f", c2, w1.Len(), r2))
	}
	t1 := w1.Atan()
	t2 := w2.Atan()
	return []float64{t1, t2}
}

func (c Circle) IntersectLineSegmentT(ls LineSegment) []float64 {
	ts := []float64{}
	line := ls.Line()
	lineTs := c.IntersectTs(line)
	if len(lineTs) == 0 {
		return nil
	}
	for _, t := range lineTs {
		if t >= 0 && t <= 1 {
			// get angle of vector from the center to the point
			v := line.At(t).Subtract(c.center)
			ts = append(ts, v.Atan())
		}
	}
	return ts
}

func (c Circle) XML(color, width string) xmlwriter.Elem {
	// <circle r="45" cx="50" cy="50" fill="red" />
	return xmlwriter.Elem{
		Name: "circle", Attrs: []xmlwriter.Attr{
			{
				Name:  "r",
				Value: fmt.Sprintf("%.1f", c.radius),
			},
			{
				Name:  "cx",
				Value: fmt.Sprintf("%.1f", c.center.x),
			},
			{
				Name:  "cy",
				Value: fmt.Sprintf("%.1f", c.center.y),
			},
			{
				Name:  "fill",
				Value: "none",
			},
			{
				Name:  "stroke-width",
				Value: width,
			},
			{
				Name:  "stroke",
				Value: color,
			},
		},
	}
}

func CircleArc(circle Circle, t1 float64, t2 float64) Path {
	p1 := circle.At(t1)
	p2 := circle.At(t2)
	isLong := false
	if t2-t1 > math.Pi {
		isLong = true
	}
	path := NewPath(p1).AddPathChunk(CircleArcChunk{
		radius:      circle.radius,
		endpoint:    p2,
		isLong:      isLong,
		isClockwise: false,
	})
	return path
}

// type CircleArc struct {
// 	circle Circle
// 	t1     float64
// 	t2     float64
// }

// func (c CircleArc) XML(color, width string) xmlwriter.Elem {
// 	p1 := c.circle.At(c.t1)
// 	p2 := c.circle.At(c.t2)
// 	path := NewPath(p1).AddPathChunk(CircleArcChunk{
// 		radius:   c.circle.radius,
// 		endpoint: p2,
// 	})
// 	return path.XML(color, width)
// 	// return xmlwriter.Elem{
// 	// 	Name: "path", Attrs: []xmlwriter.Attr{
// 	// 		{
// 	// 			Name:  "d",
// 	// 			Value: fmt.Sprintf("M %.1f %.1f A %.1f %.1f 0 0 1 %.1f %.1f", p1.x, p1.y, c.circle.radius, c.circle.radius, p2.x, p2.y),
// 	// 		},
// 	// 		{
// 	// 			Name:  "stroke",
// 	// 			Value: color,
// 	// 		},
// 	// 		{
// 	// 			Name:  "fill",
// 	// 			Value: "none",
// 	// 		},
// 	// 		{
// 	// 			Name:  "stroke-width",
// 	// 			Value: width,
// 	// 		},
// 	// 	},
// 	// }
// }

// func (c CircleArc) IsEmpty() bool {
// 	return c.t1 == c.t2
// }

// func (c CircleArc) String() string {
// 	return fmt.Sprintf("CircleArc %s from %.1f to %.1f", c.circle, c.t1, c.t2)
// }

func quadratic(a, b, c float64) []float64 {
	discriminant := b*b - 4*a*c
	if discriminant < 0.0 {
		return nil
	}
	if discriminant == 0 {
		return []float64{
			-b / (2 * a),
		}
	}
	d := math.Sqrt(discriminant)
	return []float64{
		(-b - d) / (2 * a),
		(-b + d) / (2 * a),
	}
}
