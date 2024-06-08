package main

import (
	"fmt"
	"math"

	svg "github.com/ajstarks/svgo"
)

func main() {
	fname := "gallery/test1.svg"
	sizePx := 10000.0
	padding := 1000.0

	outerBox := Box{0, 0, sizePx, sizePx}
	innerBox := outerBox.WithPadding(padding)
	scene := Scene{}.WithGuides()
	boxes := []LineLike{}
	// boxes = append(boxes, outerBox.Lines()...)
	boxes = append(boxes, innerBox.Lines()...)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(boxes).WithOffset(0, 0))
	curlyBrush := getCurlyBrush(innerBox, 400.0, math.Pi/4)
	scene = scene.AddLayer(NewLayer("Curly1").WithLineLike(curlyBrush).WithColor("red").WithWidth(10).WithOffset(-2, 40))
	curlyBrush2 := getCurlyBrush(innerBox, 300.0, math.Pi/3)
	scene = scene.AddLayer(NewLayer("Curly2").WithLineLike(curlyBrush2).WithColor("blue").WithWidth(10).WithOffset(2, -30))
	SVG{fname: fname,
		width:  "12in",
		height: "9in",
		Scene:  scene,
	}.WriteSVG()
}

func getCurlyBrush(box Box, width, angle float64) []LineLike {
	brushWidth := width
	path := CurlyFill{
		box:     box.WithPadding(brushWidth),
		angle:   angle,
		spacing: float64(brushWidth),
	}
	return []LineLike{Path{path.GetPath()}}
}

func getBrushBackForthScene(box Box) Scene {
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
	scene := Scene{}
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	for i, linelikes := range horizontalColumns.GetLineLikes() {
		scene = scene.AddLayer(NewLayer(fmt.Sprintf("Horizontal %d", i)).WithLineLike(linelikes).WithOffset(0, 0))
	}

	for i, linelikes := range verticalColumns.GetLineLikes() {
		scene = scene.AddLayer(NewLayer(fmt.Sprintf("Vertical %d", i)).WithLineLike(linelikes).WithOffset(0, 0))
	}
	return scene
}

type PlotImage interface {
	Render(*svg.SVG)         // render non-guideline layers (layesr 1+)
	DrawGuideLines(*svg.SVG) // draw guidelines (layer 0)
	GetDefs(*svg.SVG)
}
