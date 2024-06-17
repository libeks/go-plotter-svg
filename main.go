package main

import (
	"fmt"
	"math"
	"math/rand"

	svg "github.com/ajstarks/svgo"
)

func main() {
	fname := "gallery/test1.svg"
	sizePx := 10000.0
	padding := 1000.0

	outerBox := Box{0, 0, sizePx, sizePx}
	innerBox := outerBox.WithPadding(padding)
	// scene := getCurlyScene(outerBox)
	// scene := getLinesInsideScene(innerBox, 1000)
	// testLine()
	scene := getLineFieldInObjects(innerBox)
	SVG{fname: fname,
		width:  "12in",
		height: "9in",
		Scene:  scene,
	}.WriteSVG()
}

func getCurlyScene(box Box) Scene {
	scene := Scene{}.WithGuides()
	// innerBox := box.WithPadding(padding)
	// boxes := []LineLike{}
	// boxes = append(boxes, outerBox.Lines()...)
	// boxes = append(boxes, box.Lines()...)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	curlyBrush := getCurlyBrush(box, 400.0, math.Pi/4)
	scene = scene.AddLayer(NewLayer("Curly1").WithLineLike(curlyBrush).WithColor("red").WithWidth(10).WithOffset(-2, 40))
	curlyBrush2 := getCurlyBrush(box, 300.0, math.Pi/3)
	scene = scene.AddLayer(NewLayer("Curly2").WithLineLike(curlyBrush2).WithColor("blue").WithWidth(10).WithOffset(2, -30))
	return scene
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

func getLinesInsideScene(box Box, n int) Scene {
	// poly := Polygon{[]Point{
	// 	{3000, 3000},
	// 	{7000, 3000},
	// 	{7000, 7000},
	// 	{3000, 7000},
	// }}
	poly := Circle{
		Point{5000, 5000},
		1000,
	}
	return getLinesInsidePolygonScene(box, poly, n)
}

func getLinesInsidePolygonScene(box Box, poly Object, n int) Scene {
	scene := Scene{}
	lines := []LineLike{}
	for {
		if len(lines) == n {
			break
		}
		x := rand.Float64()*(box.xEnd-box.x) + box.x
		y := rand.Float64()*(box.yEnd-box.y) + box.y
		if poly.Inside(x, y) {
			lines = append(lines, LineSegment{x, y, x + 100, y})
		}
	}
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(lines).WithOffset(0, 0))
	return scene
}

func testLine() {

	line := Line{Point{0, 0}, Vector{1, 1}}
	circle := Circle{
		Point{3, 2},
		1.0,
	}
	segments := limitLinesToShape([]Line{line}, circle)
	fmt.Printf("segments: %v\n", segments)
}

func getLineFieldInObjects(box Box) Scene {
	scene := Scene{}.WithGuides()

	// poly := Circle{
	// 	Point{5000, 5000},
	// 	2000,
	// }
	poly1 := Polygon{
		[]Point{
			{3000, 3000},
			{4900, 3000},
			{4900, 7000},
			{3000, 7000},
		},
	}
	poly2 := Polygon{
		[]Point{
			{5100, 3000},
			{7000, 3000},
			{7000, 7000},
			{5100, 7000},
		},
	}
	poly3 := Polygon{
		[]Point{
			{3000, 3000},
			{7000, 3000},
			{7000, 4900},
			{3000, 4900},
		},
	}
	poly4 := Polygon{
		[]Point{
			{3000, 5100},
			{7000, 5100},
			{7000, 7000},
			{3000, 7000},
		},
	}
	radial := CircularLineField(100, Point{5000, 5000})
	// radial2 := CircularLineField(10, Point{5000, 5000})
	// segments :=
	// segments2 := limitLinesToShape(radial2, poly)
	lines1 := segmentsToLineLikes(limitLinesToShape(radial, poly1))
	lines2 := segmentsToLineLikes(limitLinesToShape(radial, poly2))
	lines3 := segmentsToLineLikes(limitLinesToShape(radial, poly3))
	lines4 := segmentsToLineLikes(limitLinesToShape(radial, poly4))
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(lines1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(lines2).WithOffset(0, 0).WithColor("blue"))
	scene = scene.AddLayer(NewLayer("content3").WithLineLike(lines3).WithOffset(0, 0).WithColor("yellow"))
	scene = scene.AddLayer(NewLayer("content4").WithLineLike(lines4).WithOffset(0, 0).WithColor("green"))
	// scene = scene.AddLayer(NewLayer("content2").WithLineLike(lines2).WithOffset(0, 0))
	return scene
}

func segmentsToLineLikes(segments []LineSegment) []LineLike {
	linelikes := make([]LineLike, len(segments))
	for i, seg := range segments {
		linelikes[i] = seg
	}
	return linelikes
}

type PlotImage interface {
	Render(*svg.SVG)         // render non-guideline layers (layesr 1+)
	DrawGuideLines(*svg.SVG) // draw guidelines (layer 0)
	GetDefs(*svg.SVG)
}
