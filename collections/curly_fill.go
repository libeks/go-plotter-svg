package collections

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

type CurlyFill struct {
	box.Box         // bounding box for strokes. The stroke should never be outside of this box
	Angle   float64 // in radians, counter-clockwise from +x direction
	Spacing float64 // distance between lines
}

func (f CurlyFill) String() string {
	return fmt.Sprintf("CurlyFill %s angle: %.3f, spacing: %.3f", f.Box, f.Angle, f.Spacing)
}

func (f CurlyFill) GetPath() lines.Path {
	// TODO: Refactor to use line field instead
	// find the starting point - extreme point of box in direction perpendicular to
	w := f.Spacing
	boxWidth := f.Box.Width()
	boxHeight := f.Box.Height()

	if f.Angle <= 0 || f.Angle >= math.Pi/2 {
		panic(fmt.Errorf("angle %.3f is not strictly in the first quadrant, not yet supported", f.Angle))
	}
	sina := math.Sin(f.Angle)
	cosa := math.Cos(f.Angle)
	tana := math.Tan(f.Angle)
	h := w / cosa

	// start at (0,h)
	x := 0.0
	y := h
	path := lines.NewPath(primitives.Point{X: f.Box.X + x, Y: f.Box.Y + y})
	i := 0
	for {
		ii := float64(i)
		cx := ((4*ii+2)*h - w) / tana // iterate over 2,6,10, ...
		cy := w

		if cx > boxWidth-w {
			cx = boxWidth - w
			dy := (boxWidth - w) * tana
			cy = ((4*ii + 2) * h) - dy
		}
		if cy > boxHeight-w {
			break
		}

		x = cx - w*sina
		y = cy - w*cosa
		end := primitives.Point{X: f.Box.X + x, Y: f.Box.Y + y}
		path = path.AddPathChunk(lines.LineChunk{End: end})

		x2 := cx + w*sina
		y2 := cy + w*cosa
		path = path.AddPathChunk(lines.CircleArcChunk{
			Radius:      w,
			Center:      primitives.Point{X: f.Box.X + cx, Y: f.Box.Y + cy},
			Start:       end,
			End:         primitives.Point{X: f.Box.X + x2, Y: f.Box.Y + y2},
			IsLong:      false,
			IsClockwise: true,
		})

		cx = w
		cy = (4*(ii+1))*h - w*tana // iterate over 4,8,12,...
		if cy > boxHeight-w {
			cy = boxHeight - w
			cx = ((4*(ii+1))*h - boxHeight + w) / tana
		}
		if cx > boxWidth-w {
			break
		}

		x = cx - w*sina
		y = cy - w*cosa
		end = primitives.Point{X: f.Box.X + x, Y: f.Box.Y + y}
		path = path.AddPathChunk(lines.LineChunk{End: end})

		x2 = cx + w*sina
		y2 = cy + w*cosa
		path = path.AddPathChunk(lines.CircleArcChunk{
			Radius:      w,
			Center:      primitives.Point{X: f.Box.X + cx, Y: f.Box.Y + cy},
			Start:       end,
			End:         primitives.Point{X: f.Box.X + x2, Y: f.Box.Y + y2},
			IsLong:      false,
			IsClockwise: false,
		})
		i += 1
	}
	return path
}

type StripImage struct {
	box.Box     // bounding box of strokes
	NGroups int // number of groups/layers to draw, spaced evenly in the box according to Direction parameters
	NLines  int // number of lines to draw in a group
	Direction
}

func (s StripImage) String() string {
	return fmt.Sprintf("StripImage %s %d groups,  %d lines", s.Box, s.NGroups, s.NLines)
}

func (s StripImage) GetLineLikes() [][]lines.LineLike {
	var barSize float64
	if s.Direction.CardinalDirection == Horizontal {
		barSize = s.Box.Width() / float64(s.NGroups)
	} else {
		barSize = s.Box.Height() / float64(s.NGroups)
	}
	padding := (s.Box.YEnd - s.Box.Y) / float64(s.NLines)
	linelikes := [][]lines.LineLike{}
	for i := range s.NGroups {
		var b box.Box
		if s.Direction.CardinalDirection == Horizontal {
			b = box.Box{X: s.Box.X + barSize*float64(i), Y: s.Box.Y, XEnd: s.Box.X + barSize*float64(i+1), YEnd: s.Box.YEnd}
		} else {
			b = box.Box{X: s.Box.X, Y: s.Box.Y + barSize*float64(i), XEnd: s.Box.XEnd, YEnd: s.Box.Y + barSize*float64(i+1)}
		}

		h := StrokeStrip{
			box:       b,
			padding:   padding,
			Direction: s.Direction,
		}

		linelikes = append(linelikes, h.Lines())
	}
	return linelikes
}
