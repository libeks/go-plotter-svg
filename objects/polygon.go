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

func (p Polygon) Translate(v primitives.Vector) Polygon {
	points := make([]primitives.Point, len(p.Points))
	for i, pt := range p.Points {
		points[i] = pt.Add(v)
	}
	return Polygon{
		Points: points,
	}
}

func (p Polygon) EdgeLines() []lines.LineSegment {
	segments := []lines.LineSegment{}
	for i, p1 := range p.Points {
		j := (i + 1) % len(p.Points)
		p2 := p.Points[j]
		segments = append(segments, lines.LineSegment{P1: p1, P2: p2})
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

func (p Polygon) BBox() primitives.BBox {
	return primitives.BBoxAroundPoints(p.Points...)
}

// Find the largest square inside the rectangle whose center is the midpoint of all corners
// this is not guaranteed to be the biggets ortogonal square in this polygon
func (p Polygon) findCenterBBox() primitives.BBox {
	var midpoint primitives.Point
	for _, pt := range p.Points {
		midpoint.X += pt.X / float64(len(p.Points))
		midpoint.Y += pt.Y / float64(len(p.Points))
	}
	outsideBox := p.BBox()

	size := max(outsideBox.Width(), outsideBox.Height())
	var bbox primitives.BBox
	// TODO: Do a binary search with a threshold instead
	for {
		bbox = primitives.BBox{
			UpperLeft:  midpoint.Add(primitives.Vector{X: -math.Sqrt(2), Y: -math.Sqrt(2)}.Mult(size / 2)),
			LowerRight: midpoint.Add(primitives.Vector{X: math.Sqrt(2), Y: math.Sqrt(2)}.Mult(size / 2)),
		}
		if p.bboxInside(bbox) {
			return bbox
		}

		size = size * 0.95
		if size < 1.0 {
			panic("Bounding box too small")
		}
	}
}

// Find the largest axis-aligned square that can be inscribed in the polygon
func (p Polygon) LargestContainedSquareBBox() primitives.BBox {
	// Current algo starts with the midpoint of all points and iteratively shrinks the square until it fits,
	// but this is not nearly optimal. Other ideas:
	//
	// * find a square that fits, then try to wiggle it around and increase the side length
	// * start by finding the midpoint of all angles (how does that work for non-triangles)
	// * for each vertical and horizontal, compute the width/height of the polygon, then use the intersection
	//   points as candidates, this also limits how far to adjust for
	// Consider https://cgm.cs.mcgill.ca/~athens/cs507/Projects/2003/DanielSud/
	bbox := p.findCenterBBox()
	oldSize := 0.0
	for bbox.Width() > oldSize {
		oldSize = bbox.Width()
		for i := range 8 {
			candidate := bbox.Translate(primitives.UnitRight.RotateCCW(math.Pi / 4.0 * float64(i)).Mult(bbox.Width() * 0.01))
			// fmt.Printf("Candidate %v\n", candidate)
			candidate = candidate.Scale(1.005)
			// fmt.Printf("Candidate2 %v\n", candidate)
			if math.Abs(1.0-candidate.Width()/candidate.Height()) > 0.001 {
				fmt.Printf("width %f, height %f\n", candidate.Width(), candidate.Height())
				panic("polygon is not a square")
			}
			if p.bboxInside(candidate) {
				// fmt.Printf("old box side %f, new %f\n", oldSize, candidate.Width())
				bbox = candidate
				break
			}
		}
	}
	return bbox
}

// https://stackoverflow.com/a/59299881
func mod(a, b int) int {
	return (a%b + b) % b
}

// Grow increases the polygon outwards (or inwards, if d is negative), by moving every
// edge line outwards perpendicular to itself by the distance d.
// If the resulting polygon is too small, this will return a polygon with no edges
func (p Polygon) Grow(d float64) Polygon {
	// TODO: Fix growing direction for CW and CCW polygons, they each grow in opposite directions
	// for every point,
	//   take the edges that the point falls on,
	//   take their lines,
	//   extend each line perpendicularly outward by d
	//   find the new intersection point of the extended lines
	edges := make([]lines.Line, len(p.Points))
	for i := range len(p.Points) {
		pointA := p.Points[mod(i-1, len(p.Points))] // ensure that it wraps around beatifully
		pointB := p.Points[i]
		edges[i] = lines.Line{P: pointA, V: pointB.Subtract(pointA).Unit()}
	}
	for i, edge := range edges {
		fmt.Printf("old edge %v\n", edge)
		edges[i] = lines.Line{P: edge.P.Add(edge.V.Perp().Unit().Mult(d)), V: edge.V}
		fmt.Printf("new edge %v\n", edges[i])
	}
	points := []primitives.Point{}
	for i := range edges {
		edgeA := edges[mod(i-1, len(edges))] // ensure that it wraps around beatifully
		edgeB := edges[i]
		intersection := edgeA.Intersect(edgeB)
		if intersection == nil {
			return Polygon{}
		}
		fmt.Printf("Intersection of %v and %v is %v\n", edgeA, edgeB, *intersection)
		points = append(points, *intersection)
	}

	return Polygon{Points: points}
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

func (p Polygon) bboxInside(bbox primitives.BBox) bool {
	pts := []primitives.Point{
		bbox.UpperLeft,
		bbox.LowerRight,
		{X: bbox.UpperLeft.X, Y: bbox.LowerRight.Y},
		{X: bbox.LowerRight.X, Y: bbox.UpperLeft.Y},
	}
	for _, pt := range pts {
		if !p.Inside(pt) {
			return false
		}
	}
	return true
}

func PolygonFromBBox(b primitives.BBox) Polygon {
	return Polygon{
		Points: b.Corners(),
	}
}
