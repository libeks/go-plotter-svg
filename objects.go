package main

import (
	"fmt"
	"math"
	"slices"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/objects"
)

func ClipLineToObject(line lines.Line, obj objects.Object) []lines.LineSegment {
	ts := obj.IntersectTs(line)
	if len(ts) == 0 {
		return nil
	}
	if len(ts) == 1 {
		panic(fmt.Errorf("not sure what to do with only one intersection %v", ts))
	}
	slices.Sort(ts)
	segments := []lines.LineSegment{}
	for i, t1 := range ts {
		if i == len(ts)-1 {
			break
		}
		t2 := ts[(i + 1)]
		midT := maths.Average(t1, t2)
		midPoint := line.At(midT)
		if obj.Inside(midPoint) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			seg := lines.LineSegment{P1: p1, P2: p2}
			segments = append(segments, seg)
		}
	}
	return segments
}

func ClipLineSegmentToObject(ls lines.LineSegment, obj objects.Object) []lines.LineSegment {
	line := ls.Line()
	ts := append(obj.IntersectTs(line), 0.0, 1.0)
	ts = filterToRange(ts, 0.0, 1.0)
	slices.Sort(ts)
	segments := []lines.LineSegment{}
	for i, t1 := range ts {
		t2 := ts[(i+1)%len(ts)]
		midT := maths.Average(t1, t2)
		midPoint := line.At(midT)
		if obj.Inside(midPoint) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			segments = append(segments, lines.LineSegment{
				P1: p1, P2: p2,
			})
		}
	}
	return segments
}

func ClipCircleToObject(c objects.Circle, obj objects.Object) []lines.LineLike {
	ts := obj.IntersectCircleTs(c)
	if len(ts) == 0 {
		return []lines.LineLike{c}
	}
	slices.Sort(ts)
	segments := []lines.LineLike{}
	for i, t1 := range ts {
		t2 := ts[(i+1)%len(ts)]
		midT := maths.Average(t1, t2)
		midPoint := c.At(midT)
		if obj.Inside(midPoint) {
			segments = append(segments, CircleArc(c, t1, t2))
		}
	}
	return segments
}

func filterToRange(vals []float64, min, max float64) []float64 {
	rets := []float64{}
	for _, val := range vals {
		if val <= max && val >= min {
			rets = append(rets, val)
		}
	}
	return rets
}

func CircleArc(circle objects.Circle, t1 float64, t2 float64) lines.Path {
	p1 := circle.At(t1)
	p2 := circle.At(t2)
	isLong := false
	if t2-t1 > math.Pi {
		isLong = true
	}
	path := lines.NewPath(p1).AddPathChunk(lines.CircleArcChunk{
		Radius:      circle.Radius,
		End:         p2,
		IsLong:      isLong,
		IsClockwise: false,
	})
	return path
}

// type CircleArc struct {
// 	circle Circle
// 	t1     float64
// 	t2     float64
// }

// func (c CircleArc) XML(color, width string) xmlwriter.Elem {
// 	p1 := c.circle.At(c.t1)
// 	p2 := c.circle.At(c.t2)
// 	path := NewPath(p1).AddPathChunk(CircleArcChunk{
// 		radius:   c.circle.radius,
// 		endpoint: p2,
// 	})
// 	return path.XML(color, width)
// 	// return xmlwriter.Elem{
// 	// 	Name: "path", Attrs: []xmlwriter.Attr{
// 	// 		{
// 	// 			Name:  "d",
// 	// 			Value: fmt.Sprintf("M %.1f %.1f A %.1f %.1f 0 0 1 %.1f %.1f", p1.x, p1.y, c.circle.radius, c.circle.radius, p2.x, p2.y),
// 	// 		},
// 	// 		{
// 	// 			Name:  "stroke",
// 	// 			Value: color,
// 	// 		},
// 	// 		{
// 	// 			Name:  "fill",
// 	// 			Value: "none",
// 	// 		},
// 	// 		{
// 	// 			Name:  "stroke-width",
// 	// 			Value: width,
// 	// 		},
// 	// 	},
// 	// }
// }

// func (c CircleArc) IsEmpty() bool {
// 	return c.t1 == c.t2
// }

// func (c CircleArc) String() string {
// 	return fmt.Sprintf("CircleArc %s from %.1f to %.1f", c.circle, c.t1, c.t2)
// }
