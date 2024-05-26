package main

import (
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"
)

var (
	// differentiating between the various line styles for clarity when visualizing
	pencilStyle = "fill:none;stroke:black;stroke-width:0.264583;stroke-opacity:1"
	brushStyle  = "fill:none;stroke:grey;stroke-width:10;stroke-opacity:1"
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
	canvas.StartviewUnit(width, height, "in", 0, 0, 1000, 1000)
	yPadding := 50
	img := StripImage{
		box:           Box{20, yPadding, 980, 1000 - yPadding},
		nVertical:     5,
		nLines:        30,
		pencilXOffset: 10,
	}
	img.Render(canvas)
	canvas.End()
}

type PlotImage interface {
	Render(*svg.SVG)         // render non-guideline layers (layesr 1+)
	DrawGuideLines(*svg.SVG) // draw guidelines (layer 0)
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

type StripImage struct {
	box           Box
	nVertical     int
	nLines        int
	pencilXOffset int
}

func (s StripImage) String() string {
	return fmt.Sprintf("StripImage %s %d verticals,  %d lines,  %d pixel offset", s.box, s.nVertical, s.nLines, s.pencilXOffset)
}

func (s StripImage) Render(canvas *svg.SVG) {
	fmt.Printf("StripImage %s\n", s)
	xWidth := (s.box.xEnd - s.box.x) / s.nVertical
	padding := (s.box.yEnd - s.box.y) / s.nLines
	strips := []HorizontalStrip{}
	for i := range s.nVertical {
		h := HorizontalStrip{
			box:       Box{x: s.box.x + xWidth*i, y: s.box.y, xEnd: s.box.x + xWidth*(i+1), yEnd: s.box.yEnd},
			padding:   padding,
			layerName: fmt.Sprintf("%d - Brush", i+1),
		}
		strips = append(strips, h)
	}
	canvas.Group(`inkscape:groupmode="layer"`, fmt.Sprintf(`inkscape:label="%s"`, "0 - Pencil"))
	for _, strip := range strips {
		canvas.Line(strip.box.x+s.pencilXOffset, strip.box.y, strip.box.x+s.pencilXOffset, strip.box.yEnd, pencilStyle)
	}
	canvas.Gend()

	for _, strip := range strips {
		strip.drawOnCanvas(canvas)
	}
}

type HorizontalStrip struct {
	box       Box
	padding   int
	layerName string
}

func (h HorizontalStrip) String() string {
	return fmt.Sprintf("HorizontalStrip %s padding %d, name '%s'", h.box, h.padding, h.layerName)
}

func (h HorizontalStrip) drawOnCanvas(canvas *svg.SVG) {
	canvas.Group(`inkscape:groupmode="layer"`, fmt.Sprintf(`inkscape:label="%s"`, h.layerName))
	fmt.Printf("box %s\n", h)
	nLines := (h.box.yEnd - h.box.y) / h.padding
	for i := range nLines + 1 {
		canvas.Line(h.box.x, h.box.y+i*h.padding, h.box.xEnd, h.box.y+i*h.padding, brushStyle)
	}
	canvas.Gend()
}
