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
	box       Box
	nVertical int
	nLines    int
	Direction
}

func (s StripImage) String() string {
	return fmt.Sprintf("StripImage %s %d verticals,  %d lines", s.box, s.nVertical, s.nLines)
}

func (s StripImage) GetLayers() []Layer {
	xWidth := (s.box.xEnd - s.box.x) / s.nVertical
	padding := (s.box.yEnd - s.box.y) / s.nLines
	layers := make([]Layer, s.nVertical)
	for i := range s.nVertical {

		h := HorizontalStrip{
			box:       Box{x: s.box.x + xWidth*i, y: s.box.y, xEnd: s.box.x + xWidth*(i+1), yEnd: s.box.yEnd},
			padding:   padding,
			layerName: fmt.Sprintf("%d - Brush", i+1),
			Direction: s.Direction,
		}
		layers[i] = Layer{
			name:  fmt.Sprintf("%d - Brush", i+1),
			i:     i,
			lines: h.Lines(),
		}
	}
	return layers
}

type HorizontalStrip struct {
	box       Box
	padding   int
	layerName string
	Direction
}

func (h HorizontalStrip) String() string {
	return fmt.Sprintf("HorizontalStrip %s padding %d, name '%s'", h.box, h.padding, h.layerName)
}

func (h HorizontalStrip) Lines() []Line {
	nLines := (h.box.yEnd-h.box.y)/h.padding + 1
	lines := make([]Line, nLines)

	for i := range nLines {
		j := i
		if h.Direction.OrderDirection == BottomToTop {
			j = nLines - i - 1
		}
		if h.Direction.StrokeDirection == LeftToRight {
			lines[i] = Line{h.box.x, h.box.y + j*h.padding, h.box.xEnd, h.box.y + j*h.padding}
		} else {
			lines[i] = Line{h.box.xEnd, h.box.y + j*h.padding, h.box.x, h.box.y + j*h.padding}
		}
		fmt.Printf("Just added line %s\n", lines[i])
	}
	return lines
}
