package scenes

import (
	"github.com/libeks/go-plotter-svg/curve"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/samplers"
)

func getTruchetScene(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(lines.LinesFromBBox(b)).WithOffset(0, 0))
	tileSource := samplers.RandomDataSource{}
	// tileSource := samplers.ConstantDataSource{0}
	// tileSource := samplers.InsideCircleSubDataSource{
	// 	Radius:  0.5,
	// 	Inside:  samplers.RandomChooser{Values: []float64{0, 1}},
	// 	Outside: samplers.ConstantDataSource{Val: 0.5},
	// }
	// tileSource := samplers.InsideCircleSubDataSource{
	// 	Radius:  0.5,
	// 	Inside:  samplers.RandomChooser{Values: []float64{0, 1}},
	// 	Outside: samplers.ConstantDataSource{Val: 0.5},
	// }
	// edgeSource := samplers.RandomDataSource{}
	edgeSource := samplers.Constant(.5) // 0.5 means we'll use default edge values
	// edgeSource := samplers.RandomChooser{Values: []float64{-.25, 1.25}}
	// edgeSource := samplers.RandomChooser{Values: []float64{.3, .7}}
	// edgeSource := samplers.RandomChooser{Values: []float64{0, 1}}
	// truch := curve.Truchet4NonCrossing
	// truch := curve.Truchet4Crossing
	truch := curve.Truchet6NonCrossingSide
	grid := curve.NewTruchetGrid(b, 30, truch, tileSource, edgeSource, curve.MapCircularCurve)
	curves := grid.GenerateCurves()
	// scene = scene.AddLayer(NewLayer("truchet").WithControlLines(curves).WithColor("blue").WithWidth(10))
	scene = scene.AddLayer(NewLayer("truchet").WithLineLike(curves).WithColor("red").WithWidth(10))
	// scene = scene.AddLayer(NewLayer("gridlines").WithLineLike(grid.GetGridLines()).WithColor("black").WithWidth(10))

	return scene
}

func getSweepTruchet(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(lines.LinesFromBBox(b)).WithOffset(0, 0))
	curves1 := curve.NewTruchetGrid(b, 3, curve.Truchet4NonCrossing, samplers.RandomDataSource{}, samplers.Constant(0.5), curve.MapCircularCircleCurve).GenerateCurves()
	curves2 := curve.NewTruchetGrid(b, 6, curve.Truchet4NonCrossing, samplers.RandomDataSource{}, samplers.Constant(0.5), curve.MapCircularCircleCurve).GenerateCurves()
	curves3 := curve.NewTruchetGrid(b, 12, curve.Truchet4NonCrossing, samplers.RandomDataSource{}, samplers.Constant(0.5), curve.MapCircularCircleCurve).GenerateCurves()
	distance := 20.0

	// scene = scene.AddLayer(NewLayer("truchet_offsets_1").WithControlLines(curves1).WithColor("gray").WithWidth(distance))
	scene = scene.AddLayer(NewLayer("truchet_offsets_1").WithLineLike(getOffsetForCurves(curves1, distance, 10)).WithColor("red").WithWidth(distance))
	scene = scene.AddLayer(NewLayer("truchet_offsets_2").WithLineLike(getOffsetForCurves(curves2, distance, 7)).WithColor("green").WithWidth(distance))
	scene = scene.AddLayer(NewLayer("truchet_offsets_3").WithLineLike(getOffsetForCurves(curves3, distance, 5)).WithColor("blue").WithWidth(distance))
	// scene = scene.AddLayer(NewLayer("gridlines").WithLineLike(grid.GetGridLines()).WithColor("black").WithWidth(10))
	return scene
}

func getSpacedTruchets(b primitives.BBox, marchingResolution int, sampler samplers.DataSource, centerThreshold, spacing float64, n int) []lines.LineLike {
	// n/2 lines will be below centerThreshold, n/2 will be above
	baseThreshold := centerThreshold - float64(n)/2*spacing
	var curves = []lines.LineLike{}
	for i := range n {
		curves = append(
			curves,
			curve.NewMarchingGrid(b, marchingResolution, sampler, baseThreshold+spacing*float64(i)).GenerateCurves()...,
		)
	}
	return curves
}
