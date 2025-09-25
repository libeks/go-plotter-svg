package box

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

// Box is an axis-aligned box
type Box struct {
	X    float64
	Y    float64
	XEnd float64
	YEnd float64
}

func (b Box) String() string {
	return fmt.Sprintf("Box (%.1f, %.1f) -> (%.1f, %.1f)", b.X, b.Y, b.XEnd, b.YEnd)
}

func (b Box) Lines() []lines.LineLike {
	path := lines.NewPath(primitives.Point{X: b.X, Y: b.Y})
	// find the starting point - extreme point of box in direction perpendicular to

	path = path.AddPathChunk(lines.LineChunk{Start: primitives.Point{X: b.X, Y: b.Y}, End: primitives.Point{X: b.X, Y: b.YEnd}})
	path = path.AddPathChunk(lines.LineChunk{Start: primitives.Point{X: b.X, Y: b.YEnd}, End: primitives.Point{X: b.XEnd, Y: b.YEnd}})
	path = path.AddPathChunk(lines.LineChunk{Start: primitives.Point{X: b.XEnd, Y: b.YEnd}, End: primitives.Point{X: b.XEnd, Y: b.Y}})
	path = path.AddPathChunk(lines.LineChunk{Start: primitives.Point{X: b.XEnd, Y: b.Y}, End: primitives.Point{X: b.X, Y: b.Y}})

	return []lines.LineLike{
		path,
	}
}

func (b Box) BBox() primitives.BBox {
	return primitives.BBox{
		UpperLeft: primitives.Point{
			X: b.X,
			Y: b.Y,
		},
		LowerRight: primitives.Point{
			X: b.XEnd,
			Y: b.YEnd,
		},
	}
}

func (b Box) Corners() []primitives.Point {
	return []primitives.Point{
		{X: b.X, Y: b.Y}, {X: b.X, Y: b.YEnd},
		{X: b.XEnd, Y: b.YEnd}, {X: b.XEnd, Y: b.Y},
	}
}

func (b Box) NWCorner() primitives.Point {
	return primitives.Point{X: b.X, Y: b.Y}
}

func (b Box) NECorner() primitives.Point {
	return primitives.Point{X: b.X, Y: b.YEnd}
}

func (b Box) SWCorner() primitives.Point {
	return primitives.Point{X: b.XEnd, Y: b.Y}
}

func (b Box) SECorner() primitives.Point {
	return primitives.Point{X: b.XEnd, Y: b.YEnd}
}

func (b Box) ClipLineToBox(l lines.Line) *lines.LineSegment {
	ls := []lines.LineSegment{
		{P1: primitives.Point{X: b.X, Y: b.Y}, P2: primitives.Point{X: b.X, Y: b.YEnd}},
		{P1: primitives.Point{X: b.X, Y: b.YEnd}, P2: primitives.Point{X: b.XEnd, Y: b.YEnd}},
		{P1: primitives.Point{X: b.XEnd, Y: b.YEnd}, P2: primitives.Point{X: b.XEnd, Y: b.Y}},
		{P1: primitives.Point{X: b.XEnd, Y: b.Y}, P2: primitives.Point{X: b.X, Y: b.Y}},
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
		b.X + pad,
		b.Y + pad,
		b.XEnd - pad,
		b.YEnd - pad,
	}
}

func (b Box) Translate(v primitives.Vector) Box {
	return Box{
		b.X + v.X,
		b.Y + v.Y,
		b.XEnd + v.X,
		b.YEnd + v.Y,
	}
}

// center of the box in absolute coordinates [0, 10_000]
func (b Box) Center() primitives.Point {
	return primitives.Point{X: b.X + (b.XEnd-b.X)/2, Y: b.Y + (b.YEnd-b.Y)/2}
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
	return primitives.Point{X: getRelativeAroundCenter(b.X + (b.XEnd-b.X)/2), Y: getRelativeAroundCenter(b.Y + (b.YEnd-b.Y)/2)}
}

func (b Box) Width() float64 {
	return b.XEnd - b.X
}

func (b Box) Height() float64 {
	return b.YEnd - b.Y
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
					X:    hh*squareSide + b.X,
					Y:    vv*squareSide + b.Y,
					XEnd: (hh+1)*squareSide + b.X,
					YEnd: (vv+1)*squareSide + b.Y,
				},
				I: h,
				J: v,
			})
		}
	}
	return boxes
}

func BoxFromBBox(b primitives.BBox) Box {
	return Box{
		X:    b.UpperLeft.X,
		Y:    b.UpperLeft.Y,
		XEnd: b.LowerRight.X,
		YEnd: b.LowerRight.Y,
	}
}

type IndexedBox struct {
	Box Box
	I   int
	J   int
}
