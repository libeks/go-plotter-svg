package box

import (
	"math"

	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

// func (b Box) ClipLineToBox(l lines.Line) *lines.LineSegment {
// 	ls := []lines.LineSegment{
// 		{P1: b.NWCorner(), P2: b.NECorner()},
// 		{P1: b.NECorner(), P2: b.SECorner()},
// 		{P1: b.SECorner(), P2: b.SWCorner()},
// 		{P1: b.SWCorner(), P2: b.NWCorner()},
// 	}
// 	ts := []float64{}
// 	for _, lineseg := range ls {
// 		if t := l.IntersectLineSegmentT(lineseg); t != nil {
// 			ts = append(ts, *t)
// 		}
// 	}
// 	if len(ts) == 0 {
// 		return nil
// 	}
// 	if len(ts) == 2 {
// 		p1 := l.At(ts[0])
// 		p2 := l.At(ts[1])
// 		return &lines.LineSegment{P1: p1, P2: p2}
// 	}
// 	panic(fmt.Errorf("line had weird number of intersections with box: %v", ts))
// }

func RelativeMinusPlusOneCenter(b, parentBox primitives.BBox) primitives.Point {
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
func RelativeCenter(b primitives.BBox) primitives.Point {
	return primitives.Point{
		X: getRelativeAroundCenter(b.UpperLeft.X + b.Width()/2.0),
		Y: getRelativeAroundCenter(b.UpperLeft.Y + b.Height()/2.0),
	}
}

func CircleInsideBox(b primitives.BBox) objects.Circle {
	return objects.Circle{
		Center: b.Center(),
		Radius: b.Width() / 2,
	}
}

func PartitionIntoSquares(b primitives.BBox, nHorizontal int) []IndexedBox {
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
				I: h,
				J: v,
			})
		}
	}
	return boxes
}

type IndexedBox struct {
	primitives.BBox
	I int
	J int
}
