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
