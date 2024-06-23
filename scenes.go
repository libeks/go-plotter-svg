package main

import (
	"fmt"
	"math"
	"math/rand"
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
)

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

func getLinesInsidePolygonScene(box Box, poly Object, n int) Scene {
	scene := Scene{}
	lines := []LineLike{}
	for {
		if len(lines) == n {
			break
		}
		x := rand.Float64()*(box.xEnd-box.x) + box.x
		y := rand.Float64()*(box.yEnd-box.y) + box.y
		if poly.Inside(Point{x, y}) {
			lines = append(lines, LineSegment{Point{x, y}, Point{x + 100, y}})
		}
	}
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(box.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(lines).WithOffset(0, 0))
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
