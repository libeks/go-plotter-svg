package primitives

import (
	"fmt"
	"math"
)

type Vector struct {
	X float64
	Y float64
}

func (v Vector) String() string {
	return fmt.Sprintf("Vector (%.1f, %.1f)", v.X, v.Y)
}

func (v Vector) Mult(t float64) Vector {
	return Vector{t * v.X, t * v.Y}
}

func (v Vector) Add(w Vector) Vector {
	return Vector{v.X + w.X, v.Y + w.Y}
}

func (v Vector) Dot(w Vector) float64 {
	return v.X*w.X + v.Y*w.Y
}

func (v Vector) Len() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector) Point() Point {
	return Point(v) // S1016 complains if I do explicit Point struct
}

// RotateCCW rotates the vector counter clockwise by t in radians
func (v Vector) RotateCCW(t float64) Vector {
	cos := math.Cos(t)
	sin := math.Sin(t)
	return Vector{
		v.X*cos - v.Y*sin,
		v.X*sin + v.Y*cos,
	}
}

// returns the angle theta of the vector counter clockwise wrt. 0x axis
func (v Vector) Atan() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v Vector) Unit() Vector {
	return v.Mult(1 / v.Len())
}

// Perp returns a vector perpendicular to v of the same lenght,
// rotated counter-clockwise by 90deg
func (v Vector) Perp() Vector {
	return Vector{-v.Y, v.X}
}
