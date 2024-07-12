package main

import (
	"fmt"
	"math"
	"slices"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

func limitLinesToShape(ls []lines.Line, shape objects.Object) []lines.LineSegment {
	segments := []lines.LineSegment{}
	for _, line := range ls {
		segs := ClipLineToObject(line, shape)
		segments = append(segments, segs...)
	}
	return segments
}

func CircularLineField(n int, center primitives.Point) []lines.Line {
	linelikes := []lines.Line{}
	for i := range n {
		ii := float64(i)
		angle := math.Pi * ii / (float64(n))
		linelikes = append(linelikes, lines.Line{
			P: center,
			V: primitives.Vector{X: 100 * math.Sin(angle), Y: 100 * math.Cos(angle)},
		})
	}
	return linelikes
}

func limitCirclesToShape(circles []objects.Circle, shape objects.Object) []lines.LineLike {
	segments := []lines.LineLike{}
	for _, circle := range circles {
		segs := ClipCircleToObject(circle, shape)
		segments = append(segments, segs...)
	}
	return segments
}

func concentricCircles(b box.Box, center primitives.Point, spacing float64) []objects.Circle {
	// find out the maximum radius
	maxDist := 0.0
	for _, p := range b.Corners() {
		v := p.Subtract(center)
		dist := v.Len()
		if dist > maxDist {
			maxDist = dist
		}
	}
	nCircles := maxDist / spacing
	circles := []objects.Circle{}
	for i := range int(nCircles) + 1 {
		ii := float64(i)
		circles = append(circles, objects.Circle{Center: center, Radius: spacing * ii})
	}
	return circles
}

// find out min and max line index. The 0th line goes through the origin (0,0)
func getLineIndexRange(b box.Box, perpVect primitives.Vector) (float64, float64) {
	iSlice := []float64{}
	wSq := perpVect.Dot(perpVect)
	for _, point := range []primitives.Point{{X: b.X, Y: b.Y}, {X: b.X, Y: b.YEnd}, {X: b.XEnd, Y: b.Y}, {X: b.XEnd, Y: b.YEnd}} {
		pVect := point.Subtract(primitives.Point{X: 0, Y: 0})
		i := pVect.Dot(perpVect) / wSq
		iSlice = append(iSlice, i)
	}
	minI := slices.Min(iSlice)
	maxI := slices.Max(iSlice)
	return minI, maxI
}

// LinearLineField returns a set of parallel lines, all oriented in the direction of angle relative to 0x axis,
// only returns the lines that would fall inside the box
func LinearLineField(b box.Box, angle float64, spacing float64) []lines.Line {
	// find out min and max line index. The 0th line goes through the origin (0,0)
	v := primitives.Vector{X: math.Cos(angle), Y: math.Sin(angle)}.Mult(spacing)
	w := v.Perp()
	minI, maxI := getLineIndexRange(b, w)
	ls := []lines.Line{}
	for i := int(minI) - 1; i <= int(maxI)+1; i++ {
		ls = append(ls, lines.Line{P: w.Mult(float64(i)).Point(), V: v})
	}
	return ls
}

// LinearDensityLineField returns a set of parallel lines, all oriented in the direction of angle relative to 0x axis,
// only returns the lines that would fall inside the box
// densityFn takes input in the range [0;1] and outputs positive values denoting the spacing to respect
// at every increment
func LinearDensityLineField(b box.Box, angle float64, densityFn func(float64) float64) []lines.Line {
	// find out min and max line index. The 0th line goes through the origin (0,0)
	v := primitives.Vector{X: math.Cos(angle), Y: math.Sin(angle)}
	w := v.Perp()
	minI, maxI := getLineIndexRange(b, w)
	ls := []lines.Line{}
	i := minI
	for i < maxI+1 {
		ls = append(ls, lines.Line{P: w.Mult(float64(i)).Point(), V: v})
		densityVal := densityFn((i - minI) / (maxI - minI))
		fmt.Printf("Densityval at %.1f %.1f \n", (i-minI)/(maxI-minI), densityVal)
		i += densityVal
	}
	return ls
}
