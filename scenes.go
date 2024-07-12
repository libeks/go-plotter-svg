package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

var (
	BrushBackForthScene    = func(box Box) Scene { return getBrushBackForthScene(box) }
	CurlyScene             = func(box Box) Scene { return getCurlyScene(box) }
	LinesInsideBoxScene    = func(box Box) Scene { return getLinesInsideScene(box, 1000) }
	LineFieldScene         = func(box Box) Scene { return getLineFieldInObjects(box) }
	RadialBoxScene         = func(box Box) Scene { return radialBoxScene(box) }
	ParallelBoxScene       = func(box Box) Scene { return parallelBoxScene(box) }
	ParallelSineFieldScene = func(box Box) Scene { return parallelSineFieldsScene(box) }
	ParallelCoherentScene  = func(box Box) Scene { return parallelCoherentSineFieldsScene(box) }
	CirclesInSquareScene   = func(box Box) Scene { return circlesInSquareScene(box) }
	TestDensityScene       = func(box Box) Scene { return testDensityScene(box) }
	TruchetScene           = func(box Box) Scene { return getTruchetScene(box) }
)

func getLineFieldInObjects(box Box) Scene {
	scene := Scene{}.WithGuides()

	poly1 := objects.Polygon{
		Points: []primitives.Point{
			{X: 3000, Y: 3000},
			{X: 4900, Y: 3000},
			{X: 4900, Y: 7000},
			{X: 3000, Y: 7000},
		},
	}
	poly2 := objects.Polygon{
		[]primitives.Point{
			{X: 5100, Y: 3000},
			{X: 7000, Y: 3000},
			{X: 7000, Y: 7000},
			{X: 5100, Y: 7000},
		},
	}
	poly3 := objects.Polygon{
		[]primitives.Point{
			{X: 3000, Y: 3000},
			{X: 7000, Y: 3000},
			{X: 7000, Y: 4900},
			{X: 3000, Y: 4900},
		},
	}
	poly4 := objects.Polygon{
		[]primitives.Point{
			{X: 3000, Y: 5100},
			{X: 7000, Y: 5100},
			{X: 7000, Y: 7000},
			{X: 3000, Y: 7000},
		},
	}
	radial := CircularLineField(3, primitives.Point{5000, 5000})
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

// circlesInSquareScene are concentric circles in a square
// due to overlap, there tends to be a darkening on the left with certain pens
func circlesInSquareScene(box Box) Scene {
	scene := Scene{}.WithGuides()
	layer1 := limitCirclesToShape(
		concentricCircles(
			box, box.Center(), 100,
		),
		box.AsPolygon(),
	)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	return scene
}

func testDensityScene(box Box) Scene {
	scene := Scene{}.WithGuides()
	quarters := partitionIntoSquares(box, 2)
	colors := []string{
		"red",
		"green",
		"blue",
		"orange",
	}
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	for i, quarter := range quarters {
		lineLikes := []lines.LineLike{}
		quarterBox := quarter.box.WithPadding(50)
		testBoxes := partitionIntoSquares(quarterBox, 5)
		fmt.Printf("got %d boxes\n", len(testBoxes))
		for _, tbox := range testBoxes {
			jj := float64(tbox.j)
			ii := float64(tbox.i)
			var ls []lines.LineLike
			tboxx := tbox.box.WithPadding(50)
			spacing := (jj + 1) * 10
			if tbox.i < 4 {
				ls = segmentsToLineLikes(
					limitLinesToShape(
						LinearLineField(
							tboxx, ii*math.Pi/4, spacing,
						),
						tboxx.AsPolygon(),
					),
				)
			} else {
				ls = limitCirclesToShape(
					concentricCircles(
						tboxx, tboxx.Center(), spacing,
					),
					tboxx.AsPolygon(),
				)
			}
			lineLikes = append(lineLikes, ls...)
		}
		layerName := fmt.Sprintf("pen %d", i)
		scene = scene.AddLayer(NewLayer(layerName).WithLineLike(lineLikes).WithOffset(0, 0).WithColor(colors[i]))
	}

	// scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	return scene
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

func parallelBoxScene(box Box) Scene {
	minLineWidth := 20.0
	maxLineWidth := 100.0
	minAngle := 0.0
	maxAngle := math.Pi
	// angle := math.Pi / 3
	scene := Scene{}.WithGuides()
	segments := [][]lines.LineLike{}
	boxes := partitionIntoSquares(box, 10)
	for _, minibox := range boxes {
		spacing := randInRange(minLineWidth, maxLineWidth)
		angle := randInRange(minAngle, maxAngle)
		lines := LinearLineField(minibox.box, angle, spacing)
		lineseg := limitLinesToShape(lines, minibox.box.WithPadding(50).AsPolygon())
		segments = append(segments, segmentsToLineLikes(lineseg))
	}
	layer1, layer2 := randomlyAllocateSegments(segments, 0.5)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(layer2).WithOffset(0, 0).WithColor("blue"))
	return scene
}

func radialBoxScene(box Box) Scene {
	nSegments := 15
	exclusionRadius := 100.0
	wiggle := 200.0
	scene := Scene{}.WithGuides()
	segments := [][]lines.LineLike{}
	boxes := partitionIntoSquares(box, 10)
	for _, minibox := range boxes {
		boxcenter := minibox.box.Center()
		xwiggle := randRangeMinusPlusOne() * wiggle
		ywiggle := randRangeMinusPlusOne() * wiggle
		center := primitives.Point{X: boxcenter.X + xwiggle, Y: boxcenter.Y + ywiggle}
		segments = append(segments, radialBoxWithCircleExclusion(minibox.box.WithPadding(50).AsPolygon(), center, nSegments, exclusionRadius))
	}
	layer1 := []lines.LineLike{}
	layer2 := []lines.LineLike{}
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

func getLinesInsidePolygonScene(box Box, poly objects.Object, n int) Scene {
	scene := Scene{}
	ls := []lines.LineLike{}
	for {
		if len(ls) == n {
			break
		}
		x := rand.Float64()*(box.xEnd-box.x) + box.x
		y := rand.Float64()*(box.yEnd-box.y) + box.y
		if poly.Inside(primitives.Point{X: x, Y: y}) {
			ls = append(ls, lines.LineSegment{P1: primitives.Point{X: x, Y: y}, P2: primitives.Point{X: x + 100, Y: y}})
		}
	}
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(ls).WithOffset(0, 0))
	return scene
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

func getCurlyScene(box Box) Scene {
	scene := Scene{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	curlyBrush := getCurlyBrush(box, 400.0, math.Pi/4)
	scene = scene.AddLayer(NewLayer("Curly1").WithLineLike(curlyBrush).WithColor("red").WithWidth(10).WithOffset(-2, 40))
	curlyBrush2 := getCurlyBrush(box, 300.0, math.Pi/3)
	scene = scene.AddLayer(NewLayer("Curly2").WithLineLike(curlyBrush2).WithColor("blue").WithWidth(10).WithOffset(2, -30))
	return scene
}

func getTruchetScene(box Box) Scene {
	scene := Scene{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	dataSource := RandomDataSource{}
	// dataSource := ConstantDataSource{1.0}
	// grid := NewGrid(box, 30, dataSource, truchetTiles)
	grid := NewGrid(box, 30, dataSource, truchetTilesWithStrikeThrough)
	curves := grid.GererateCurves()
	scene = scene.AddLayer(NewLayer("Curly1").WithLineLike(curves).WithColor("red").WithWidth(10))

	return scene
}

type IndexedBox struct {
	box Box
	i   int
	j   int
}

func partitionIntoSquares(box Box, nHorizontal int) []IndexedBox {
	width := box.Width()
	squareSide := width / (float64(nHorizontal))
	boxes := []IndexedBox{}
	verticalIterations := int(box.Height() / float64(squareSide))
	if verticalIterations < nHorizontal && math.Abs(box.Height()-(float64(nHorizontal)*float64(squareSide))) < 0.1 {
		verticalIterations = nHorizontal
	}
	for v := range verticalIterations {
		vv := float64(v)
		for h := range nHorizontal {
			hh := float64(h)
			boxes = append(boxes, IndexedBox{
				box: Box{
					x:    hh*squareSide + box.x,
					y:    vv*squareSide + box.y,
					xEnd: (hh+1)*squareSide + box.x,
					yEnd: (vv+1)*squareSide + box.y,
				},
				i: h,
				j: v,
			})
		}
	}
	return boxes
}
