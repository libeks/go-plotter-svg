package lines

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/primitives"
)

type Line struct {
	P primitives.Point
	V primitives.Vector
}

func (l Line) String() string {
	return fmt.Sprintf("Line (%s, %s)", l.P, l.V)
}

// Return a point on the line that is t lenghts of v away from p.
func (l Line) At(t float64) primitives.Point {
	return l.P.Add(l.V.Mult(t))
}

func (l Line) Intersect(l2 Line) *primitives.Point {
	// TODO: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
	// first line is x1,y1 = l.p, x2,y2 = l.p + l.v. so x2-x1 = l.v.x, y2-y1 = l.v.y
	// second line is x3,y3 = l2.p, x4,y4 = l2.p + l2.v, so x4-x3 = l2.v.x, y4-y3 = l2.v.y
	// determinant is (x1-x2)(y3-y4) - (y1-y2)(x3-x4) = (x2-x1)(y4-y3) - (y2-y1)(x4-x3)
	determinant := l.V.X*l2.V.Y - l.V.Y*l2.V.X
	if determinant == 0 {
		return nil
	}
	// result is
	// x = ((x1*y2 - y1*x2)(x3-x4) - (x1-x2)(x3*y4 - y3*x4))/determinant
	//   = ()
	// y = ((x1*y2 - y1*x2)(y3-y4) - (y1-y2)(x3*y4 - y3*x4))/determinant
	x1x2 := -l.V.X
	x3x4 := -l2.V.X
	y1y2 := -l.V.Y
	y3y4 := -l2.V.Y
	x1y2y1x2 := (l.P.X*(l.P.Y+l.V.Y) - l.P.Y*(l.P.X+l.V.X))                 // x1*(y2) - y1*x2
	x3y4y3x4 := (l2.P.X * (l2.P.Y + l2.V.Y)) - (l2.P.Y * (l2.P.X + l2.V.X)) // x4*(y4+y3) - y4*(x3+x4)

	x := (x1y2y1x2*x3x4 - x1x2*x3y4y3x4) / determinant
	y := (x1y2y1x2*y3y4 - y1y2*x3y4y3x4) / determinant
	return &primitives.Point{X: x, Y: y}
}

// Return the intersect t parameter of the line l when intersecting line l2
func (l Line) IntersectT(l2 Line) *float64 {
	// TODO: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
	// first line is x1,y1 = l.p, x2,y2 = l.p + l.v. so x2-x1 = l.v.x, y2-y1 = l.v.y
	// second line is x3,y3 = l2.p, x4,y4 = l2.p + l2.v, so x4-x3 = l2.v.x, y4-y3 = l2.v.y
	x1x2 := -l.V.X
	x3x4 := -l2.V.X
	y1y2 := -l.V.Y
	y3y4 := -l2.V.Y
	x1x3 := l.P.X - l2.P.X
	y1y3 := l2.P.Y - l2.P.Y

	divisor := (x1x2*y3y4 - y1y2*x3x4)
	if divisor == 0.0 {
		return nil
	}
	t := (x1x3*y3y4 - y1y3*x3x4) / divisor
	return &t
}

// Return the intersection parameters t,u for both lines l and l2
func (l Line) IntersectTU(l2 Line) (*float64, *float64) {
	// TODO: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
	// first line is x1,y1 = l.p, x2,y2 = l.p + l.v. so x2-x1 = l.v.x, y2-y1 = l.v.y
	// second line is x3,y3 = l2.p, x4,y4 = l2.p + l2.v, so x4-x3 = l2.v.x, y4-y3 = l2.v.y

	x1x2 := -l.V.X
	x3x4 := -l2.V.X
	y1y2 := -l.V.Y
	y3y4 := -l2.V.Y
	x1x3 := l.P.X - l2.P.X
	y1y3 := l.P.Y - l2.P.Y

	divisor := (x1x2*y3y4 - y1y2*x3x4)
	if divisor == 0.0 {
		return nil, nil
	}
	t := (x1x3*y3y4 - y1y3*x3x4) / divisor
	u := (x1x2*y1y3 - y1y2*x1x3) / -divisor // note the divisor is negative here. I initially missed that.
	return &t, &u
}

func (l Line) IntersectLineSegmentT(ls2 LineSegment) *float64 {
	l2 := ls2.Line()
	t, u := l.IntersectTU(l2)
	if t == nil || u == nil {
		return nil
	}
	uu := *u
	if uu <= 1.0 && uu >= 0.0 {
		return t
	}
	return nil
}
