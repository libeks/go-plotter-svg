package main

import (
	"fmt"
	"math"
	"slices"
)

type Object interface {
	Inside(x, y float64) bool
	IntersectTs(line Line) []float64
}

func ClipLineToObject(line Line, obj Object) []LineSegment {
	ts := obj.IntersectTs(line)
	// fmt.Printf("Line intersecting object at %v\n", ts)
	if len(ts) == 0 {
		return nil
	}
	if len(ts) == 1 {
		panic(fmt.Errorf("not sure what to do with only one intersection %v", ts))
	}
	slices.Sort(ts)
	// fmt.Printf("ts: %v\n", ts)
	segments := []LineSegment{}
	for i, t1 := range ts {
		if i == len(ts)-1 {
			break
		}
		t2 := ts[(i + 1)]
		midT := average(t1, t2)
		midPoint := line.At(midT)
		if obj.Inside(midPoint.x, midPoint.y) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			seg := LineSegment{p1.x, p1.y, p2.x, p2.y}
			segments = append(segments, seg)
			// fmt.Printf("adding segment %s at index %d\n", seg, i)
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
		if obj.Inside(midPoint.x, midPoint.y) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			segments = append(segments, LineSegment{
				p1.x, p1.y,
				p2.x, p2.y,
			})
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

func (o CompositeObject) Inside(x, y float64) bool {
	inside := false
	for _, pos := range o.positive {
		if pos.Inside(x, y) {
			inside = true
			break
		}
	}
	if !inside {
		return false
	}
	for _, neg := range o.negative {
		if neg.Inside(x, y) {
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

func positiveATan(y, x float64) float64 {
	angle := math.Atan2(y, x)
	return angle
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

func (p Polygon) Inside(x, y float64) bool {
	// compute the winding angle from the point
	totalAngle := 0.0
	for i, p1 := range p.points {
		j := (i + 1) % len(p.points)
		p2 := p.points[j]
		p2Angle := positiveATan(p2.y-y, p2.x-x)
		p1Angle := positiveATan(p1.y-y, p1.x-x)
		angle := angleDifference(p2Angle, p1Angle)
		totalAngle += angle
	}
	if math.Abs(totalAngle-math.Pi) < 0.01 {
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
		segments = append(segments, LineSegment{p1.x, p1.y, p2.x, p2.y})
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

func average(a, b float64) float64 {
	return (a + b) / 2
}

type Circle struct {
	center Point
	radius float64
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle @%s, r:%.1f", c.center, c.radius)
}

func (c Circle) Inside(x, y float64) bool {
	distance := Point{x, y}.Subtract(c.center).Len()
	return distance <= c.radius
}

func (c Circle) IntersectTs(line Line) []float64 {
	w := line.p.Subtract(c.center)
	// fmt.Printf("w is %s\n", w)
	A := line.v.x*line.v.x + line.v.y*line.v.y
	B := 2 * (line.v.x*w.x + line.v.y*w.y)
	C := w.x*w.x + w.y*w.y - c.radius*c.radius
	ts := quadratic(A, B, C)
	if len(ts) < 2 {
		return nil
	}
	return ts
}

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
