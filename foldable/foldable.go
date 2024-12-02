package foldable

import (
	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

const (
	flapWidth = 100
)

func Cube(b box.Box, side float64) []lines.LineLike {
	start := primitives.Point{X: b.X, Y: b.Y + side}
	lns := []lines.LineLike{}
	// draws the cube as follows:
	//
	//     +---+
	//  /-\| 4 |/-\ /-\
	// +---+---+---+---+\
	// | 0 | 1 | 2 | 3 ||
	// +---+---+---+---+/
	//  \-/| 5 |\-/ \-/
	//     +---+

	// face 0
	l := lines.NewPath(start)
	end := start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: -side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 1
	start = primitives.Point{X: b.X + side, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 2
	start = primitives.Point{X: b.X + side*2, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 3
	start = primitives.Point{X: b.X + side*3, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 4
	start = primitives.Point{X: b.X + side, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: 0, Y: -side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 5
	start = primitives.Point{X: b.X + side, Y: b.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: -side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 04
	start = primitives.Point{X: b.X, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 24
	start = primitives.Point{X: b.X + side*2, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 34
	start = primitives.Point{X: b.X + side*3, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 05
	start = primitives.Point{X: b.X, Y: b.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 25
	start = primitives.Point{X: b.X + side*2, Y: b.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 35
	start = primitives.Point{X: b.X + side*3, Y: b.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 30
	start = primitives.Point{X: b.X + side*4, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side - flapWidth*2})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	return lns
}
