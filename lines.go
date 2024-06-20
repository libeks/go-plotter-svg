package main

import (
	"fmt"
	"math"
	"slices"
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

// find out min and max line index. The 0th line goes through the origin (0,0)
func getLineIndexRange(box Box, perpVect Vector) (float64, float64) {
	iSlice := []float64{}
	// v := Vector{math.Cos(angle), math.Sin(angle)}
	// w := v.Perp()
	// w := Vector{-math.Sin(angle), math.Cos(angle)}.Mult(spacing) // v rotated 90deg counter-clockwise
	wSq := perpVect.Dot(perpVect)
	for _, point := range []Point{{box.x, box.y}, {box.x, box.yEnd}, {box.xEnd, box.y}, {box.xEnd, box.yEnd}} {
		pVect := point.Subtract(Point{0, 0})
		i := pVect.Dot(perpVect) / wSq
		iSlice = append(iSlice, i)
	}
	minI := slices.Min(iSlice)
	maxI := slices.Max(iSlice)
	return minI, maxI
}

// LinearLineField returns a set of parallel lines, all oriented in the direction of angle relative to 0x axis,
// only returns the lines that would fall inside the box
func LinearLineField(box Box, angle float64, spacing float64) []Line {
	// find out min and max line index. The 0th line goes through the origin (0,0)
	v := Vector{math.Cos(angle), math.Sin(angle)}.Mult(spacing)
	w := v.Perp()
	// w := Vector{-math.Sin(angle), math.Cos(angle)}.Mult(spacing) // v rotated 90deg counter-clockwise
	minI, maxI := getLineIndexRange(box, w)
	lines := []Line{}
	for i := int(minI) - 1; i <= int(maxI)+1; i++ {
		lines = append(lines, Line{p: w.Mult(float64(i)).Point(), v: v})
	}
	return lines
}

// LinearDensityLineField returns a set of parallel lines, all oriented in the direction of angle relative to 0x axis,
// only returns the lines that would fall inside the box
// densityFn takes input in the range [0;1] and outputs positive values denoting the spacing to respect
// at every increment
func LinearDensityLineField(box Box, angle float64, densityFn func(float64) float64) []Line {
	// find out min and max line index. The 0th line goes through the origin (0,0)
	v := Vector{math.Cos(angle), math.Sin(angle)}
	w := v.Perp()
	// w := Vector{-math.Sin(angle), math.Cos(angle)}.Mult(spacing) // v rotated 90deg counter-clockwise
	minI, maxI := getLineIndexRange(box, w)
	lines := []Line{}
	i := minI
	for i < maxI+1 {
		lines = append(lines, Line{p: w.Mult(float64(i)).Point(), v: v})
		densityVal := densityFn((i - minI) / (maxI - minI))
		fmt.Printf("Densityval at %.1f %.1f \n", (i-minI)/(maxI-minI), densityVal)
		i += densityVal
	}
	return lines
}
