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
	// scene := getLineFieldInObjects(innerBox)
	// scene := radialBoxScene(innerBox)
	// scene := parallelBoxScene(innerBox)
	// scene := parallelSineFieldsScene(innerBox)
	scene := parallelCoherentSineFieldsScene(innerBox)
	SVG{fname: fname,
		width:  "12in",
		height: "9in",
		Scene:  scene,
	}.WriteSVG()
}

func getCurlyScene(box Box) Scene {
	scene := Scene{}.WithGuides()
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
			lines = append(lines, LineSegment{Point{x, y}, Point{x + 100, y}})
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

func partitionIntoSquares(box Box, nHorizontal int) []Box {
	width := box.Width()
	squareSide := width / (float64(nHorizontal))
	boxes := []Box{}
	verticalIteractions := int(box.Height() / float64(squareSide))
	for v := range verticalIteractions {
		vv := float64(v)
		for h := range nHorizontal {
			hh := float64(h)
			boxes = append(boxes, Box{
				x:    hh*squareSide + box.x,
				y:    vv*squareSide + box.y,
				xEnd: (hh+1)*squareSide + box.x,
				yEnd: (vv+1)*squareSide + box.y,
			})
		}
	}
	return boxes
}

func randRangeMinusPlusOne() float64 {
	return 2 * (rand.Float64() - 0.5)
}

func radialBoxScene(box Box) Scene {
	nSegments := 15
	exclusionRadius := 100.0
	wiggle := 200.0
	scene := Scene{}.WithGuides()
	segments := [][]LineLike{}
	boxes := partitionIntoSquares(box, 10)
	for _, minibox := range boxes {
		boxcenter := minibox.Center()
		xwiggle := randRangeMinusPlusOne() * wiggle
		ywiggle := randRangeMinusPlusOne() * wiggle
		center := Point{boxcenter.x + xwiggle, boxcenter.y + ywiggle}
		segments = append(segments, radialBoxWithCircleExclusion(minibox.WithPadding(50).AsPolygon(), center, nSegments, exclusionRadius))
	}
	layer1 := []LineLike{}
	layer2 := []LineLike{}
	for _, segs := range segments {
		if rand.Float64() > 0.5 {
			layer1 = append(layer1, segs...)
		} else {
			layer2 = append(layer2, segs...)
		}
	}
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(layer2).WithOffset(0, 0).WithColor("blue"))
	return scene
}

func randInRange(min, max float64) float64 {
	return (max-min)*rand.Float64() + min
}

func randomlyAllocateSegments(segments [][]LineLike, threshold float64) ([]LineLike, []LineLike) {
	layer1 := []LineLike{}
	layer2 := []LineLike{}
	for _, segs := range segments {
		if rand.Float64() > threshold {
			layer1 = append(layer1, segs...)
		} else {
			layer2 = append(layer2, segs...)
		}
	}
	return layer1, layer2
}

func parallelBoxScene(box Box) Scene {
	minLineWidth := 20.0
	maxLineWidth := 100.0
	minAngle := 0.0
	maxAngle := math.Pi
	// angle := math.Pi / 3
	scene := Scene{}.WithGuides()
	segments := [][]LineLike{}
	boxes := partitionIntoSquares(box, 10)
	for _, minibox := range boxes {
		spacing := randInRange(minLineWidth, maxLineWidth)
		angle := randInRange(minAngle, maxAngle)
		lines := LinearLineField(minibox, angle, spacing)
		lineseg := limitLinesToShape(lines, minibox.WithPadding(50).AsPolygon())
		segments = append(segments, segmentsToLineLikes(lineseg))
	}
	layer1, layer2 := randomlyAllocateSegments(segments, 0.5)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(layer2).WithOffset(0, 0).WithColor("blue"))
	return scene
}

type SineDensity struct {
	min    float64
	max    float64
	offset float64
	cycles float64
}

func (d SineDensity) Density(a float64) float64 {
	theta := d.cycles * (a + d.offset) * math.Pi
	dRange := d.max - d.min
	return d.min + dRange*(math.Sin(theta)+1)/2
}

func parallelSineFieldsScene(box Box) Scene {
	scene := Scene{}.WithGuides()
	layer1 := segmentsToLineLikes(
		limitLinesToShape(
			LinearDensityLineField(
				box, math.Pi/3, SineDensity{min: 20.0, max: 200, cycles: 5, offset: 0}.Density,
			),
			box.AsPolygon(),
		),
	)
	layer2 := segmentsToLineLikes(
		limitLinesToShape(
			LinearDensityLineField(
				box, 0.6, SineDensity{min: 20.0, max: 200, cycles: 3, offset: 0}.Density,
			),
			box.AsPolygon(),
		),
	)
	layer3 := segmentsToLineLikes(
		limitLinesToShape(
			LinearDensityLineField(
				box, 2.0, SineDensity{min: 20.0, max: 200, cycles: 7, offset: 0}.Density,
			),
			box.AsPolygon(),
		),
	)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(layer2).WithOffset(0, 0).WithColor("blue"))
	scene = scene.AddLayer(NewLayer("content3").WithLineLike(layer3).WithOffset(0, 0).WithColor("green"))
	return scene
}

func parallelCoherentSineFieldsScene(box Box) Scene {
	scene := Scene{}.WithGuides()
	layer1 := segmentsToLineLikes(
		limitLinesToShape(
			LinearDensityLineField(
				box, math.Pi/3, SineDensity{min: 20.0, max: 200, cycles: 7, offset: 0}.Density,
			),
			box.AsPolygon(),
		),
	)
	layer2 := segmentsToLineLikes(
		limitLinesToShape(
			LinearDensityLineField(
				box, math.Pi/3, SineDensity{min: 20.0, max: 200, cycles: 7, offset: 0.1}.Density,
			),
			box.AsPolygon(),
		),
	)
	layer3 := segmentsToLineLikes(
		limitLinesToShape(
			LinearDensityLineField(
				box, math.Pi/3, SineDensity{min: 20.0, max: 200, cycles: 7, offset: 0.2}.Density,
			),
			box.AsPolygon(),
		),
	)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(layer2).WithOffset(0, 0).WithColor("blue"))
	scene = scene.AddLayer(NewLayer("content3").WithLineLike(layer3).WithOffset(0, 0).WithColor("green"))
	return scene
}

func radialBoxWithCircleExclusion(container Object, center Point, nLines int, radius float64) []LineLike {
	radial := CircularLineField(nLines, center)
	compObject := NewComposite().With(container).Without(Circle{center, radius})
	lines := limitLinesToShape(radial, compObject)
	segments := segmentsToLineLikes(lines)
	return segments
}

func getLineFieldInObjects(box Box) Scene {
	scene := Scene{}.WithGuides()

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
	radial := CircularLineField(3, Point{5000, 5000})
	fmt.Printf("radial : %s\n", radial)
	lines1 := segmentsToLineLikes(limitLinesToShape(radial, poly1))
	fmt.Printf("linelikes: %s\n", lines1)
	lines2 := segmentsToLineLikes(limitLinesToShape(radial, poly2))
	lines3 := segmentsToLineLikes(limitLinesToShape(radial, poly3))
	lines4 := segmentsToLineLikes(limitLinesToShape(radial, poly4))
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(lines1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(lines2).WithOffset(0, 0).WithColor("blue"))
	scene = scene.AddLayer(NewLayer("content3").WithLineLike(lines3).WithOffset(0, 0).WithColor("yellow"))
	scene = scene.AddLayer(NewLayer("content4").WithLineLike(lines4).WithOffset(0, 0).WithColor("green"))
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
