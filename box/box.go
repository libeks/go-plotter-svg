package box

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

// Box is an axis-aligned box, a wrapper around primitives.BBox with more advanced actions
type Box struct {
	primitives.BBox
}

func (b Box) String() string {
	return fmt.Sprintf("Box (%.1f, %.1f) -> (%.1f, %.1f)", b.UpperLeft.X, b.UpperLeft.Y, b.LowerRight.X, b.LowerRight.Y)
}

func (b Box) Corners() []primitives.Point {
	return []primitives.Point{
		b.NWCorner(), b.NECorner(),
		b.SWCorner(), b.SECorner(),
	}
}

func (b Box) NWCorner() primitives.Point {
	return b.UpperLeft
}

func (b Box) NECorner() primitives.Point {
	return primitives.Point{X: b.UpperLeft.X, Y: b.LowerRight.Y}
}

func (b Box) SWCorner() primitives.Point {
	return primitives.Point{X: b.LowerRight.X, Y: b.UpperLeft.Y}
}

func (b Box) SECorner() primitives.Point {
	return b.LowerRight
}

func (b Box) ClipLineToBox(l lines.Line) *lines.LineSegment {
	ls := []lines.LineSegment{
		{P1: b.NWCorner(), P2: b.NECorner()},
		{P1: b.NECorner(), P2: b.SECorner()},
		{P1: b.SECorner(), P2: b.SWCorner()},
		{P1: b.SWCorner(), P2: b.NWCorner()},
	}
	ts := []float64{}
	for _, lineseg := range ls {
		if t := l.IntersectLineSegmentT(lineseg); t != nil {
			ts = append(ts, *t)
		}
	}
	if len(ts) == 0 {
		return nil
	}
	if len(ts) == 2 {
		p1 := l.At(ts[0])
		p2 := l.At(ts[1])
		return &lines.LineSegment{P1: p1, P2: p2}
	}
	panic(fmt.Errorf("line had weird number of intersections with box: %v", ts))
}

func (b Box) WithPadding(pad float64) Box {
	return Box{
		BBox: b.BBox.WithPadding(pad),
	}
}

func (b Box) Translate(v primitives.Vector) Box {
	return Box{
		BBox: primitives.BBox{
			UpperLeft:  b.BBox.UpperLeft.Add(v),
			LowerRight: b.BBox.LowerRight.Add(v),
		},
	}
}

func (b Box) RelativeMinusPlusOneCenter(parentBox Box) primitives.Point {
	center := b.Center()
	parentCenter := parentBox.Center()
	return primitives.Point{
		X: 2 * (center.X - parentCenter.X) / parentBox.Width(),
		Y: 2 * (center.Y - parentCenter.Y) / parentBox.Height(),
	}
}

func getRelativeAroundCenter(v float64) float64 {
	relative := (v) / 10_000
	return 2*relative - 1
}

// center of the box in relative coordinates [0.0, 1.0], assuming that the image is in the range [0, 10_000]
func (b Box) RelativeCenter() primitives.Point {
	return primitives.Point{
		X: getRelativeAroundCenter(b.UpperLeft.X + (b.LowerRight.X-b.UpperLeft.X)/2),
		Y: getRelativeAroundCenter(b.UpperLeft.Y + (b.LowerRight.Y-b.UpperLeft.Y)/2),
	}
}

func (b Box) AsPolygon() objects.Polygon {
	return objects.Polygon{
		Points: b.Corners(),
	}
}

func (b Box) CircleInsideBox() objects.Circle {
	return objects.Circle{
		Center: b.Center(),
		Radius: b.Width() / 2,
	}
}

func (b Box) PartitionIntoSquares(nHorizontal int) []IndexedBox {
	width := b.Width()
	squareSide := width / (float64(nHorizontal))
	boxes := []IndexedBox{}
	verticalIterations := int(b.Height() / float64(squareSide))
	if verticalIterations < nHorizontal && math.Abs(b.Height()-(float64(nHorizontal)*float64(squareSide))) < 0.1 {
		verticalIterations = nHorizontal
	}
	for v := range verticalIterations {
		vv := float64(v)
		for h := range nHorizontal {
			hh := float64(h)
			boxes = append(boxes, IndexedBox{
				Box: Box{
					BBox: primitives.BBox{
						UpperLeft: primitives.Point{
							X: hh*squareSide + b.UpperLeft.X,
							Y: vv*squareSide + b.UpperLeft.Y,
						},
						LowerRight: primitives.Point{
							X: (hh+1)*squareSide + b.UpperLeft.X,
							Y: (vv+1)*squareSide + b.UpperLeft.Y,
						},
					},
				},
				I: h,
				J: v,
			})
		}
	}
	return boxes
}

type IndexedBox struct {
	Box Box
	I   int
	J   int
}
