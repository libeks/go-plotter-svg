package main

import (
	"fmt"
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
	padding := 500

	box := Box{padding, padding, sizePx - padding, sizePx - padding}

	container := NewPlotContainer()

	brushColumns := &StripImage{
		box:       box,
		nVertical: 1,
		nLines:    30,
		Direction: Direction{
			StrokeDirection: RightToLeft,
			OrderDirection:  BottomToTop,
			Connection:      SameDirection,
		},
	}
	container = container.WithLayers(brushColumns.GetLayers()...)
	container = container.WithLines(box.Lines()...)
	canvas.Def()
	defs := container.GetDefs(canvas)
	canvas.DefEnd()

	container.Render(canvas, defs)
	canvas.End()
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
