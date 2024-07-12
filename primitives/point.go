package primitives

import "fmt"

var (
	Origin = Point{X: 0, Y: 0}
)

type Point struct {
	X float64
	Y float64
}

func (p Point) String() string {
	return fmt.Sprintf("Point (%.1f, %.1f)", p.X, p.Y)
}

func (p Point) Add(v Vector) Point {
	return Point{p.X + v.X, p.Y + v.Y}
}

func (p Point) Subtract(p2 Point) Vector {
	return Vector{p.X - p2.X, p.Y - p2.Y}
}