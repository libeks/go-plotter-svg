package scenes

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/collections"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/samplers"
	"github.com/libeks/go-plotter-svg/truchet"
)

var (
	BrushBackForthScene    = func(b box.Box) Scene { return getBrushBackForthScene(b) }
	CurlyScene             = func(b box.Box) Scene { return getCurlyScene(b) }
	LinesInsideBoxScene    = func(b box.Box) Scene { return getLinesInsideScene(b, 1000) }
	LineFieldScene         = func(b box.Box) Scene { return getLineFieldInObjects(b) }
	RadialBoxScene         = func(b box.Box) Scene { return radialBoxScene(b) }
	ParallelBoxScene       = func(b box.Box) Scene { return parallelBoxScene(b) }
	ParallelSineFieldScene = func(b box.Box) Scene { return parallelSineFieldsScene(b) }
	ParallelCoherentScene  = func(b box.Box) Scene { return parallelCoherentSineFieldsScene(b) }
	CirclesInSquareScene   = func(b box.Box) Scene { return circlesInSquareScene(b) }
	TestDensityScene       = func(b box.Box) Scene { return testDensityScene(b) }
	TruchetScene           = func(b box.Box) Scene { return getTruchetScene(b) }
)

func getLineFieldInObjects(b box.Box) Scene {
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
		Points: []primitives.Point{
			{X: 5100, Y: 3000},
			{X: 7000, Y: 3000},
			{X: 7000, Y: 7000},
			{X: 5100, Y: 7000},
		},
	}
	poly3 := objects.Polygon{
		Points: []primitives.Point{
			{X: 3000, Y: 3000},
			{X: 7000, Y: 3000},
			{X: 7000, Y: 4900},
			{X: 3000, Y: 4900},
		},
	}
	poly4 := objects.Polygon{
		Points: []primitives.Point{
			{X: 3000, Y: 5100},
			{X: 7000, Y: 5100},
			{X: 7000, Y: 7000},
			{X: 3000, Y: 7000},
		},
	}
	radial := collections.CircularLineField(3, primitives.Point{X: 5000, Y: 5000})
	fmt.Printf("radial : %s\n", radial)
	lines1 := segmentsToLineLikes(collections.LimitLinesToShape(radial, poly1))
	fmt.Printf("linelikes: %s\n", lines1)
	lines2 := segmentsToLineLikes(collections.LimitLinesToShape(radial, poly2))
	lines3 := segmentsToLineLikes(collections.LimitLinesToShape(radial, poly3))
	lines4 := segmentsToLineLikes(collections.LimitLinesToShape(radial, poly4))
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(lines1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(lines2).WithOffset(0, 0).WithColor("blue"))
	scene = scene.AddLayer(NewLayer("content3").WithLineLike(lines3).WithOffset(0, 0).WithColor("yellow"))
	scene = scene.AddLayer(NewLayer("content4").WithLineLike(lines4).WithOffset(0, 0).WithColor("green"))
	return scene
}

func parallelCoherentSineFieldsScene(b box.Box) Scene {
	scene := Scene{}.WithGuides()
	layer1 := segmentsToLineLikes(
		collections.LimitLinesToShape(
			collections.LinearDensityLineField(
				b, math.Pi/3, collections.SineDensity{Min: 20.0, Max: 200, Cycles: 7, Offset: 0}.Density,
			),
			b.AsPolygon(),
		),
	)
	layer2 := segmentsToLineLikes(
		collections.LimitLinesToShape(
			collections.LinearDensityLineField(
				b, math.Pi/3, collections.SineDensity{Min: 20.0, Max: 200, Cycles: 7, Offset: 0.1}.Density,
			),
			b.AsPolygon(),
		),
	)
	layer3 := segmentsToLineLikes(
		collections.LimitLinesToShape(
			collections.LinearDensityLineField(
				b, math.Pi/3, collections.SineDensity{Min: 20.0, Max: 200, Cycles: 7, Offset: 0.2}.Density,
			),
			b.AsPolygon(),
		),
	)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(layer2).WithOffset(0, 0).WithColor("blue"))
	scene = scene.AddLayer(NewLayer("content3").WithLineLike(layer3).WithOffset(0, 0).WithColor("green"))
	return scene
}

// circlesInSquareScene are concentric circles in a square
// due to overlap, there tends to be a darkening on the left with certain pens
func circlesInSquareScene(b box.Box) Scene {
	scene := Scene{}.WithGuides()
	layer1 := collections.LimitCirclesToShape(
		collections.ConcentricCircles(
			b, b.Center(), 100,
		),
		b.AsPolygon(),
	)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	return scene
}

func testDensityScene(b box.Box) Scene {
	scene := Scene{}.WithGuides()
	quarters := b.PartitionIntoSquares(2)
	colors := []string{
		"red",
		"green",
		"blue",
		"orange",
	}
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	for i, quarter := range quarters {
		lineLikes := []lines.LineLike{}
		quarterBox := quarter.Box.WithPadding(50)
		testBoxes := quarterBox.PartitionIntoSquares(5)
		fmt.Printf("got %d boxes\n", len(testBoxes))
		for _, tbox := range testBoxes {
			jj := float64(tbox.J)
			ii := float64(tbox.I)
			var ls []lines.LineLike
			tboxx := tbox.Box.WithPadding(50)
			spacing := (jj + 1) * 10
			if tbox.I < 4 {
				ls = segmentsToLineLikes(
					collections.LimitLinesToShape(
						collections.LinearLineField(
							tboxx, ii*math.Pi/4, spacing,
						),
						tboxx.AsPolygon(),
					),
				)
			} else {
				ls = collections.LimitCirclesToShape(
					collections.ConcentricCircles(
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
	return scene
}

func parallelSineFieldsScene(b box.Box) Scene {
	scene := Scene{}.WithGuides()
	layer1 := segmentsToLineLikes(
		collections.LimitLinesToShape(
			collections.LinearDensityLineField(
				b, math.Pi/3, collections.SineDensity{Min: 20.0, Max: 200, Cycles: 5, Offset: 0}.Density,
			),
			b.AsPolygon(),
		),
	)
	layer2 := segmentsToLineLikes(
		collections.LimitLinesToShape(
			collections.LinearDensityLineField(
				b, 0.6, collections.SineDensity{Min: 20.0, Max: 200, Cycles: 3, Offset: 0}.Density,
			),
			b.AsPolygon(),
		),
	)
	layer3 := segmentsToLineLikes(
		collections.LimitLinesToShape(
			collections.LinearDensityLineField(
				b, 2.0, collections.SineDensity{Min: 20.0, Max: 200, Cycles: 7, Offset: 0}.Density,
			),
			b.AsPolygon(),
		),
	)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(layer2).WithOffset(0, 0).WithColor("blue"))
	scene = scene.AddLayer(NewLayer("content3").WithLineLike(layer3).WithOffset(0, 0).WithColor("green"))
	return scene
}

func parallelBoxScene(b box.Box) Scene {
	minLineWidth := 20.0
	maxLineWidth := 100.0
	minAngle := 0.0
	maxAngle := math.Pi
	// angle := math.Pi / 3
	scene := Scene{}.WithGuides()
	segments := [][]lines.LineLike{}
	boxes := b.PartitionIntoSquares(10)
	for _, minibox := range boxes {
		spacing := maths.RandInRange(minLineWidth, maxLineWidth)
		angle := maths.RandInRange(minAngle, maxAngle)
		lines := collections.LinearLineField(minibox.Box, angle, spacing)
		lineseg := collections.LimitLinesToShape(lines, minibox.Box.WithPadding(50).AsPolygon())
		segments = append(segments, segmentsToLineLikes(lineseg))
	}
	layer1, layer2 := collections.RandomlyAllocateSegments(segments, 0.5)
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(layer2).WithOffset(0, 0).WithColor("blue"))
	return scene
}

func radialBoxScene(b box.Box) Scene {
	nSegments := 15
	exclusionRadius := 100.0
	wiggle := 200.0
	scene := Scene{}.WithGuides()
	segments := [][]lines.LineLike{}
	boxes := b.PartitionIntoSquares(10)
	for _, minibox := range boxes {
		boxcenter := minibox.Box.Center()
		xwiggle := maths.RandRangeMinusPlusOne() * wiggle
		ywiggle := maths.RandRangeMinusPlusOne() * wiggle
		center := primitives.Point{X: boxcenter.X + xwiggle, Y: boxcenter.Y + ywiggle}
		segments = append(segments, radialBoxWithCircleExclusion(minibox.Box.WithPadding(50).AsPolygon(), center, nSegments, exclusionRadius))
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
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(layer1).WithOffset(0, 0).WithColor("red"))
	scene = scene.AddLayer(NewLayer("content2").WithLineLike(layer2).WithOffset(0, 0).WithColor("blue"))
	return scene
}

func getLinesInsidePolygonScene(b box.Box, poly objects.Object, n int) Scene {
	scene := Scene{}
	ls := []lines.LineLike{}
	for {
		if len(ls) == n {
			break
		}
		x := rand.Float64()*(b.XEnd-b.X) + b.X
		y := rand.Float64()*(b.YEnd-b.Y) + b.Y
		if poly.Inside(primitives.Point{X: x, Y: y}) {
			ls = append(ls, lines.LineSegment{P1: primitives.Point{X: x, Y: y}, P2: primitives.Point{X: x + 100, Y: y}})
		}
	}
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	scene = scene.AddLayer(NewLayer("content").WithLineLike(ls).WithOffset(0, 0))
	return scene
}

func getLinesInsideScene(b box.Box, n int) Scene {
	poly := objects.Circle{
		Center: primitives.Point{X: 5000, Y: 5000},
		Radius: 1000,
	}
	return getLinesInsidePolygonScene(b, poly, n)
}

func radialBoxWithCircleExclusion(container objects.Object, center primitives.Point, nLines int, radius float64) []lines.LineLike {
	radial := collections.CircularLineField(nLines, center)
	compObject := objects.NewComposite().With(container).Without(objects.Circle{Center: center, Radius: radius})
	lines := collections.LimitLinesToShape(radial, compObject)
	segments := segmentsToLineLikes(lines)
	return segments
}

func getBrushBackForthScene(b box.Box) Scene {
	horizontalColumns := &collections.StripImage{
		Box:     b,
		NGroups: 1,
		NLines:  30,
		Direction: collections.Direction{
			CardinalDirection: collections.Horizontal,
			StrokeDirection:   collections.AwayToHome,
			OrderDirection:    collections.AwayToHome,
			Connection:        collections.SameDirection,
		},
	}
	verticalColumns := &collections.StripImage{
		Box:     b,
		NGroups: 1,
		NLines:  30,
		Direction: collections.Direction{
			CardinalDirection: collections.Vertical,
			StrokeDirection:   collections.AwayToHome,
			OrderDirection:    collections.AwayToHome,
			Connection:        collections.AlternatingDirection,
		},
	}
	scene := Scene{}
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	for i, linelikes := range horizontalColumns.GetLineLikes() {
		scene = scene.AddLayer(NewLayer(fmt.Sprintf("Horizontal %d", i)).WithLineLike(linelikes).WithOffset(0, 0))
	}

	for i, linelikes := range verticalColumns.GetLineLikes() {
		scene = scene.AddLayer(NewLayer(fmt.Sprintf("Vertical %d", i)).WithLineLike(linelikes).WithOffset(0, 0))
	}
	return scene
}

func getCurlyScene(b box.Box) Scene {
	scene := Scene{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	curlyBrush := getCurlyBrush(b, 400.0, math.Pi/4)
	scene = scene.AddLayer(NewLayer("Curly1").WithLineLike(curlyBrush).WithColor("red").WithWidth(10).WithOffset(-2, 40))
	curlyBrush2 := getCurlyBrush(b, 300.0, math.Pi/3)
	scene = scene.AddLayer(NewLayer("Curly2").WithLineLike(curlyBrush2).WithColor("blue").WithWidth(10).WithOffset(2, -30))
	return scene
}

func getCurlyBrush(b box.Box, width, angle float64) []lines.LineLike {
	brushWidth := width
	path := collections.CurlyFill{
		Box:     b.WithPadding(brushWidth),
		Angle:   angle,
		Spacing: float64(brushWidth),
	}
	return []lines.LineLike{path.GetPath()}
}

func getTruchetScene(b box.Box) Scene {
	scene := Scene{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(b.Lines()).WithOffset(0, 0))
	dataSource := samplers.RandomDataSource{}
	// dataSource := samplers.ConstantDataSource{Val: .5}
	// grid := NewGrid(box, 30, dataSource, truchetTiles)
	// truchetTile := truchet.TruchetUnderPairs
	grid := truchet.NewGrid(b, 50, truchet.EndpointMapping4, truchet.TruchetPairs, dataSource)
	// grid := truchet.NewGrid(b, 40, truchet.EndpointMapping6Side, truchet.Truchet6Pairs, dataSource)
	curves := grid.GererateCurves()
	scene = scene.AddLayer(NewLayer("truchet").WithLineLike(curves).WithColor("red").WithWidth(10))
	// scene = scene.AddLayer(NewLayer("gridlines").WithLineLike(grid.GetGridLines()).WithColor("black").WithWidth(10))

	return scene
}
