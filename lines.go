package main

import (
	"math"
)

func limitLinesToShape(lines []Line, shape Object) []LineSegment {
	segments := []LineSegment{}
	for _, line := range lines {
		segs := ClipLineToObject(line, shape)
		// segs := shape.IntersectLine(line)
		// fmt.Printf("line %s intersects with %s at %v\n", line, shape, segs)
		segments = append(segments, segs...)
	}
	return segments
}

func CircularLineField(n int, center Point) []Line {
	lines := []Line{}
	for i := range n {
		ii := float64(i)
		angle := math.Pi * ii / (float64(n))
		lines = append(lines, Line{
			p: center,
			v: Vector{100 * math.Sin(angle), 100 * math.Cos(angle)},
		})
	}
	return lines
}

// func ParallelLineField()
