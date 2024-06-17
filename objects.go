package main

import (
	"fmt"
	"math"
	"slices"
)

type CompositeObject struct {
	base        Object
	newObject   Object
	subtraction bool
}

type Object interface {
	Inside(x, y float64) bool
	IntersectLine(line Line) []LineSegment
	IntersectLineSegment(ls LineSegment) []LineSegment
	// Outline() LineLike
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

func (p Polygon) IntersectLine(line Line) []LineSegment {
	ts := []float64{}
	for _, segment := range p.EdgeLines() {
		if t := line.IntersectLineSegmentT(segment); t != nil {
			ts = append(ts, *t)
		}
	}
	fmt.Printf("Line intersecting polygon at %v\n", ts)
	if len(ts) == 0 {
		return nil
	}
	if len(ts) == 1 {
		panic(fmt.Errorf("not sure what to do with only one intersection %v", ts))
	}
	slices.Sort(ts)
	segments := []LineSegment{}
	for i, t1 := range ts {
		t2 := ts[(i+1)%len(ts)]
		midT := average(t1, t2)
		midPoint := line.At(midT)
		if p.Inside(midPoint.x, midPoint.y) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			segments = append(segments, LineSegment{p1.x, p1.y, p2.x, p2.y})
		}
	}
	return segments
}

func (p Polygon) IntersectLineSegment(ls LineSegment) []LineSegment {

	ts := []float64{0.0, 1.0} // the segment's endpoints should also be considered
	for _, segment := range p.EdgeLines() {
		if t := ls.IntersectLineSegmentT(segment); t != nil {
			ts = append(ts, *t)
		}
	}
	slices.Sort(ts)
	line := ls.Line()
	segments := []LineSegment{}
	// todo, generalize this midpoint calculation, ensure that line segments are not broken up
	for i, t1 := range ts {
		t2 := ts[(i+1)%len(ts)]
		midT := average(t1, t2)
		midPoint := line.At(midT)
		if p.Inside(midPoint.x, midPoint.y) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			segments = append(segments, LineSegment{p1.x, p1.y, p2.x, p2.y})
		}
	}
	return segments
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

func (c Circle) IntersectLine(line Line) []LineSegment {
	w := line.p.Subtract(c.center)
	fmt.Printf("w is %s\n", w)
	A := line.v.x*line.v.x + line.v.y*line.v.y
	B := 2 * (line.v.x*w.x + line.v.y*w.y)
	C := w.x*w.x + w.y*w.y - c.radius*c.radius
	ts := quadratic(A, B, C)
	if len(ts) < 2 {
		return nil
	}
	p1 := line.At(ts[0])
	p2 := line.At(ts[1])
	for i, p := range []Point{p1, p2} {
		if r := p.Subtract(c.center).Len(); r-c.radius > 0.1 {
			panic(fmt.Errorf("t: %.1f point %s not on circle %s, r is %.1f, not %.1f", ts[i], p, c, r, c.radius))
		}
	}
	return []LineSegment{
		{
			p1.x, p1.y,
			p2.x, p2.y,
		},
	}
}

func (c Circle) IntersectLineSegment(ls LineSegment) []LineSegment {
	line := ls.Line()
	pc := line.p.Subtract(c.center)
	A := line.v.x*line.v.x + line.v.y*line.v.y
	B := 2 * (line.v.x*pc.x + line.v.y*pc.y)
	C := pc.x*pc.x + pc.y*pc.y
	ts := []float64{0.0, 1.0}
	ts = append(ts, quadratic(A, B, C)...)
	slices.Sort(ts)
	segments := []LineSegment{}
	for i, t1 := range ts {
		t2 := ts[(i+1)%len(ts)]
		midT := average(t1, t2)
		midPoint := line.At(midT)
		if c.Inside(midPoint.x, midPoint.y) {
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
	fmt.Printf("a: %.1f, b: %.1f, c: %.1f, d^2: %.1f, d: %.1f, x1: %.1f, x2: %.1f\n", a, b, c, discriminant, d, (-b-d)/(2*a), (-b+d)/(2*a))
	return []float64{
		(-b - d) / (2 * a),
		(-b + d) / (2 * a),
	}
}
