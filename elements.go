package main

import (
	"fmt"
	"math"
	"strings"
)

type CurlyFill struct {
	box     Box     // bounding box for strokes. The stroke should never be outside of this box
	angle   float64 // in radians, counter-clockwise from +x direction
	spacing float64 // distance between lines
}

func (f CurlyFill) String() string {
	return fmt.Sprintf("CurlyFill %s angle: %.3f, spacing: %.3f", f.box, f.angle, f.spacing)
}

func (f CurlyFill) GetPath() string {
	commands := []string{}
	// find the starting point - extreme point of box in direction perpendicular to
	w := f.spacing
	boxWidth := float64(f.box.xEnd - f.box.x)
	boxHeight := float64(f.box.yEnd - f.box.y)

	if f.angle <= 0 || f.angle >= math.Pi/2 {
		panic(fmt.Errorf("angle %.3f is not strictly in the first quadrant, not yet supported", f.angle))
	}
	sina := math.Sin(f.angle)
	cosa := math.Cos(f.angle)
	tana := math.Tan(f.angle)
	h := w / cosa

	// start at (0,h)
	x := 0.0
	y := h
	commands = append(commands, fmt.Sprintf("M %.3f %.3f", float64(f.box.x)+x, float64(f.box.y)+y))
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
		commands = append(commands, fmt.Sprintf("L %.3f %.3f", float64(f.box.x)+x, float64(f.box.y)+y))

		x2 := cx + w*sina
		y2 := cy + w*cosa
		commands = append(commands, fmt.Sprintf("A %.3f %.3f 0 0 1 %.3f %.3f", w, w, float64(f.box.x)+x2, float64(f.box.y)+y2))

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
		commands = append(commands, fmt.Sprintf("L %.3f %.3f", float64(f.box.x)+x, float64(f.box.y)+y))

		x2 = cx + w*sina
		y2 = cy + w*cosa
		commands = append(commands, fmt.Sprintf("A %.3f %.3f 0 0 0 %.3f %.3f", w, w, float64(f.box.x)+x2, float64(f.box.y)+y2))
		i += 1
	}
	return strings.Join(commands, " ")
}

type StripImage struct {
	box     Box // bounding box of strokes
	nGroups int // number of groups/layers to draw, spaced evenly in the box according to Direction parameters
	nLines  int // number of lines to draw in a group
	Direction
}

func (s StripImage) String() string {
	return fmt.Sprintf("StripImage %s %d groups,  %d lines", s.box, s.nGroups, s.nLines)
}

func (s StripImage) GetLineLikes() [][]LineLike {
	var barSize float64
	if s.Direction.CardinalDirection == Horizontal {
		barSize = (s.box.xEnd - s.box.x) / float64(s.nGroups)
	} else {
		barSize = (s.box.yEnd - s.box.y) / float64(s.nGroups)
	}
	padding := (s.box.yEnd - s.box.y) / float64(s.nLines)
	linelikes := [][]LineLike{}
	for i := range s.nGroups {
		var box Box
		if s.Direction.CardinalDirection == Horizontal {
			box = Box{x: s.box.x + barSize*float64(i), y: s.box.y, xEnd: s.box.x + barSize*float64(i+1), yEnd: s.box.yEnd}
		} else {
			box = Box{x: s.box.x, y: s.box.y + barSize*float64(i), xEnd: s.box.xEnd, yEnd: s.box.y + barSize*float64(i+1)}
		}

		h := StrokeStrip{
			box:       box,
			padding:   padding,
			Direction: s.Direction,
		}

		linelikes = append(linelikes, h.Lines())
	}
	return linelikes
}

type StrokeStrip struct {
	box     Box
	padding float64
	Direction
}

func (h StrokeStrip) String() string {
	return fmt.Sprintf("StrokeStrip %s padding %.1f", h.box, h.padding)
}

func (h StrokeStrip) Lines() []LineLike {
	var nLines int
	if h.Direction.CardinalDirection == Horizontal {
		nLines = int((h.box.yEnd-h.box.y)/h.padding) + 1
	} else {
		nLines = int((h.box.xEnd-h.box.x)/h.padding) + 1
	}
	lines := make([]LineLike, nLines)

	for i := range nLines {
		j := i
		if h.Direction.OrderDirection == AwayToHome {
			j = nLines - i - 1
		}
		reverse := (h.Direction.StrokeDirection == AwayToHome)
		if h.Direction.Connection == AlternatingDirection && (i%2 == 1) {
			reverse = !reverse
		}
		var line LineSegment
		if h.Direction.CardinalDirection == Horizontal {
			line = LineSegment{h.box.x, h.box.y + float64(j)*h.padding, h.box.xEnd, h.box.y + float64(j)*h.padding}
		} else {
			line = LineSegment{h.box.x + float64(j)*h.padding, h.box.y, h.box.x + float64(j)*h.padding, h.box.yEnd}
		}
		if reverse {
			lines[i] = line.Reverse()
		} else {
			lines[i] = line
		}
		fmt.Printf("Just added line %s\n", lines[i])
	}
	return lines
}
