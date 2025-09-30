package primitives

import (
	"fmt"
	"math"
	"slices"
)

type BBox struct {
	UpperLeft  Point
	LowerRight Point
}

func (b BBox) IsEmpty() bool {
	// A bbox that is of width or height 0 isn't necessarily empty, this could happen when intersecting
	// with a vertical or horizontal line

	// positive area means it is not empty
	if b.Area() > 0 {
		return false
	}
	// if both corners are the same (their distance is zero), it is empty
	if b.UpperLeft.Subtract(b.LowerRight).Len() == 0.0 {
		return true
	}
	// if the points are not ordered correctly, the bbox is empty
	if b.UpperLeft.X > b.LowerRight.X || b.UpperLeft.Y > b.LowerRight.Y {
		// fmt.Printf("bbox is empty due to border condition\n")
		return true
	}
	// otherwise the area is 0, but the points are not the same
	return false
}

func (b BBox) Width() float64 {
	return max(b.LowerRight.X-b.UpperLeft.X, 0.0)
}

func (b BBox) Height() float64 {
	return max(b.LowerRight.Y-b.UpperLeft.Y, 0.0)
}

func (b BBox) Area() float64 {
	// clamp both coords to 0.0 in case the order of the points is inverted
	return b.Width() * b.Height()
}

// Add combines two bounding boxes and returns a box that contains both
func (b BBox) Add(c BBox) BBox {
	return BBoxAroundPoints(b.UpperLeft, b.LowerRight, c.UpperLeft, c.LowerRight)
}

func (b BBox) Translate(v Vector) BBox {
	b.UpperLeft = b.UpperLeft.Add(v)
	b.LowerRight = b.LowerRight.Add(v)
	return b
}

// Scale the box around the centerpoint by ratio r
// It does so by moving the two corners in the right direction
func (b BBox) Scale(r float64) BBox {
	size := (r - 1) * b.Width()
	upperLeft := UnitRight.RotateCCW((5.0 / 4.0) * math.Pi).Mult(size)
	lowerRight := UnitRight.RotateCCW((1.0 / 4.0) * math.Pi).Mult(size)
	fmt.Printf("ul %v, lr %v\n", upperLeft, lowerRight)
	return BBox{
		UpperLeft:  b.UpperLeft.Add(upperLeft),
		LowerRight: b.LowerRight.Add(lowerRight),
	}
}

func (b BBox) Intersect(c BBox) (BBox, bool) {
	upperLeftX := max(b.UpperLeft.X, c.UpperLeft.X)
	upperLeftY := max(b.UpperLeft.Y, c.UpperLeft.Y)
	lowerRightX := min(b.LowerRight.X, c.LowerRight.X)
	lowerRightY := min(b.LowerRight.Y, c.LowerRight.Y)
	if upperLeftX > lowerRightX || upperLeftY > lowerRightY {
		return BBox{}, false
	}
	return BBox{UpperLeft: Point{upperLeftX, upperLeftY}, LowerRight: Point{lowerRightX, lowerRightY}}, true
}

func (b BBox) DoesIntersect(c BBox) bool {
	_, doesIntersect := b.Intersect(c)
	return doesIntersect
}

func BBoxAroundPoints(pts ...Point) BBox {
	if len(pts) == 0 {
		return BBox{UpperLeft: Origin, LowerRight: Origin}
	}
	xes := make([]float64, len(pts))
	ys := make([]float64, len(pts))
	for i, pt := range pts {
		xes[i] = pt.X
		ys[i] = pt.Y
	}
	minX := slices.Min(xes)
	maxX := slices.Max(xes)
	minY := slices.Min(ys)
	maxY := slices.Max(ys)
	return BBox{UpperLeft: Point{X: minX, Y: minY}, LowerRight: Point{X: maxX, Y: maxY}}
}
