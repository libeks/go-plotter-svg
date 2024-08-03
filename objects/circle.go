package objects

import (
	"fmt"
	"math"

	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/primitives"
)

type Circle struct {
	Center primitives.Point
	Radius float64
}

func (c Circle) Len() float64 {
	return math.Pi * 2 * c.Radius
}

func (c Circle) Start() primitives.Point {
	// usually the plotter plots circles starting and ending at the left extreme point
	return c.Center.Add(primitives.Vector{X: -1, Y: 0}.Mult(c.Radius))
}

func (c Circle) End() primitives.Point {
	// usually the plotter plots circles starting and ending at the left extreme point
	return c.Center.Add(primitives.Vector{X: -1, Y: 0}.Mult(c.Radius))
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle @%s, r:%.1f", c.Center, c.Radius)
}

func (c Circle) Inside(p primitives.Point) bool {
	distance := p.Subtract(c.Center).Len()
	return distance <= c.Radius
}

func (c Circle) At(t float64) primitives.Point {
	return c.Center.Add(primitives.Vector{X: c.Radius, Y: 0}.RotateCCW(t))
}

func (c Circle) IsEmpty() bool {
	return c.Radius == 0
}

// return the line ts when intersecting with the circle
func (c Circle) IntersectTs(line lines.Line) []float64 {
	w := line.P.Subtract(c.Center)
	A := line.V.X*line.V.X + line.V.Y*line.V.Y
	B := 2 * (line.V.X*w.X + line.V.Y*w.Y)
	C := w.X*w.X + w.Y*w.Y - c.Radius*c.Radius
	ts := maths.Quadratic(A, B, C)
	if len(ts) < 2 {
		return nil
	}
	return ts
}

// should return c2 t-values
// IntersectCircleTs returns the angles respective to c2 at which it intersects c1
func (c1 Circle) IntersectCircleTs(c2 Circle) []float64 {
	// distance between the two centers
	dVect := c1.Center.Subtract(c2.Center)
	d := dVect.Len()
	r1 := c1.Radius
	r2 := c2.Radius
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
	p1 := c1.Center.Add(v1)
	p2 := c1.Center.Add(v2)
	w1 := p1.Subtract(c2.Center)
	w2 := p2.Subtract(c2.Center)
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

func (c Circle) IntersectLineSegmentT(ls lines.LineSegment) []float64 {
	ts := []float64{}
	line := ls.Line()
	lineTs := c.IntersectTs(line)
	if len(lineTs) == 0 {
		return nil
	}
	for _, t := range lineTs {
		if t >= 0 && t <= 1 {
			// get angle of vector from the center to the point
			v := line.At(t).Subtract(c.Center)
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
				Value: fmt.Sprintf("%.1f", c.Radius),
			},
			{
				Name:  "cx",
				Value: fmt.Sprintf("%.1f", c.Center.X),
			},
			{
				Name:  "cy",
				Value: fmt.Sprintf("%.1f", c.Center.Y),
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

func (c Circle) ControlLineXML(color, width string) xmlwriter.Elem {
	return c.XML(color, width)
}

func (c Circle) OffsetLeft(distance float64) lines.LineLike {
	// assume that a circle is drawn counter-clockwise

	newRadius := c.Radius + distance
	// radiusRatio := newRadius / c.Radius
	return Circle{
		Center: c.Center,
		Radius: newRadius,
	}
}
