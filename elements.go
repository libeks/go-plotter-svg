package main

import (
	"fmt"

	svg "github.com/ajstarks/svgo"
)

var (
	// differentiating between the various line styles for clarity when visualizing
	pencilStyle = "fill:none;stroke:black;stroke-width:2;stroke-opacity:1"
)

func brushStyle(canvas *svg.SVG, i int) string {
	colors := []string{
		"black",
		"red",
		"orange",
		"yellow",
		"green",
		"cyan",
		"blue",
		"magenta",
	}
	gradientID := fmt.Sprintf("brush%d", i)
	startColor := svg.Offcolor{
		Offset:  0,
		Color:   colors[i%len(colors)],
		Opacity: 1.0,
	}
	endColor := svg.Offcolor{
		Offset:  100,
		Color:   "grey",
		Opacity: 1.0,
	}
	canvas.LinearGradient(gradientID, 0, 0, 100, 0, []svg.Offcolor{startColor, endColor})
	brushStyle := fmt.Sprintf("fill:none;stroke:url(#%s);stroke-width:100;stroke-opacity:0.5", gradientID)
	return brushStyle
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

func (s StripImage) GetLayers(firstLayerID int) []Layer {
	var barSize int
	if s.Direction.CardinalDirection == Horizontal {
		barSize = (s.box.xEnd - s.box.x) / s.nGroups
	} else {
		barSize = (s.box.yEnd - s.box.y) / s.nGroups
	}
	padding := (s.box.yEnd - s.box.y) / s.nLines
	layers := make([]Layer, s.nGroups)
	for i := range s.nGroups {
		var box Box
		if s.Direction.CardinalDirection == Horizontal {
			box = Box{x: s.box.x + barSize*i, y: s.box.y, xEnd: s.box.x + barSize*(i+1), yEnd: s.box.yEnd}
		} else {
			box = Box{x: s.box.x, y: s.box.y + barSize*i, xEnd: s.box.xEnd, yEnd: s.box.y + barSize*(i+1)}
		}

		h := StrokeStrip{
			box:       box,
			padding:   padding,
			layerName: fmt.Sprintf("%d - Brush", i+firstLayerID+1),
			Direction: s.Direction,
		}
		layers[i] = Layer{
			name:  fmt.Sprintf("%d - Brush", i+firstLayerID+1),
			i:     i + firstLayerID,
			lines: h.Lines(),
		}
	}
	return layers
}

type StrokeStrip struct {
	box       Box
	padding   int
	layerName string
	Direction
}

func (h StrokeStrip) String() string {
	return fmt.Sprintf("StrokeStrip %s padding %d, name '%s'", h.box, h.padding, h.layerName)
}

func (h StrokeStrip) Lines() []Line {
	var nLines int
	if h.Direction.CardinalDirection == Horizontal {
		nLines = (h.box.yEnd-h.box.y)/h.padding + 1
	} else {
		nLines = (h.box.xEnd-h.box.x)/h.padding + 1
	}
	lines := make([]Line, nLines)

	for i := range nLines {
		j := i
		if h.Direction.OrderDirection == AwayToHome {
			j = nLines - i - 1
		}
		reverse := (h.Direction.StrokeDirection == AwayToHome)
		if h.Direction.Connection == AlternatingDirection && (i%2 == 1) {
			reverse = !reverse
		}
		var line Line
		if h.Direction.CardinalDirection == Horizontal {
			line = Line{h.box.x, h.box.y + j*h.padding, h.box.xEnd, h.box.y + j*h.padding}
		} else {
			line = Line{h.box.x + j*h.padding, h.box.y, h.box.x + j*h.padding, h.box.yEnd}
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
