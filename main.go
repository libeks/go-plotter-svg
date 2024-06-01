package main

import (
	"fmt"
	"math"
	"os"

	svg "github.com/ajstarks/svgo"
)

func main() {
	fname := "gallery/test.svg"
	genSVG(fname)
}

func genSVG(fname string) {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	width := 12 // inches, == 1152.000px
	height := 9 // inches, == 864.000px

	canvas := svg.New(f)
	sizePx := 10000
	canvas.StartviewUnit(width, height, "in", 0, 0, sizePx, sizePx)
	padding := 1000

	box := Box{0, 0, sizePx, sizePx}.WithPadding(padding)
	// box := Box{padding, padding, sizePx - padding, sizePx - padding}

	container := NewPlotContainer()
	// layers, lines := getBrushBackForth(box)
	layers, lines := getCurlyBrush(box)

	container = container.WithLayers(layers...)
	container = container.WithLines(lines...)
	canvas.Def()
	defs := container.GetDefs(canvas)
	canvas.DefEnd()

	container.Render(canvas, defs)
	canvas.End()
}

func getCurlyBrush(box Box) ([]Layer, []Line) {
	brushWidth := 20
	path := CurlyFill{
		box:     box.WithPadding(brushWidth),
		angle:   math.Pi / 4,
		spacing: float64(brushWidth),
	}
	layers := []Layer{
		{
			name:  "1 - Brush",
			i:     1,
			paths: []Path{{path.GetPath()}},
		},
	}
	lines := box.Lines()
	return layers, lines
}

func getBrushBackForth(box Box) ([]Layer, []Line) {
	horizontalColumns := &StripImage{
		box:     box,
		nGroups: 1,
		nLines:  30,
		Direction: Direction{
			CardinalDirection: Horizontal,
			StrokeDirection:   AwayToHome,
			OrderDirection:    AwayToHome,
			Connection:        SameDirection,
		},
	}
	verticalColumns := &StripImage{
		box:     box,
		nGroups: 1,
		nLines:  30,
		Direction: Direction{
			CardinalDirection: Vertical,
			StrokeDirection:   AwayToHome,
			OrderDirection:    AwayToHome,
			Connection:        AlternatingDirection,
		},
	}
	layers := horizontalColumns.GetLayers(0)
	layers = append(layers, verticalColumns.GetLayers(len(layers))...)
	lines := box.Lines()
	return layers, lines
}

type PlotImage interface {
	Render(*svg.SVG)         // render non-guideline layers (layesr 1+)
	DrawGuideLines(*svg.SVG) // draw guidelines (layer 0)
	GetDefs(*svg.SVG)
}

type Box struct {
	x    int
	y    int
	xEnd int
	yEnd int
}

func (b Box) String() string {
	return fmt.Sprintf("Box (%d, %d) -> (%d, %d)", b.x, b.y, b.xEnd, b.yEnd)
}

func (b Box) Lines() []Line {
	return []Line{
		{b.x, b.y, b.x, b.yEnd},
		{b.x, b.yEnd, b.xEnd, b.yEnd},
		{b.xEnd, b.yEnd, b.xEnd, b.y},
		{b.xEnd, b.y, b.x, b.y},
	}
}

func (b Box) WithPadding(pad int) Box {
	return Box{
		b.x + pad,
		b.y + pad,
		b.xEnd - pad,
		b.yEnd - pad,
	}
}
