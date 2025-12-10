package scenes

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/curve"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/samplers"
)

func getRandomMarchingSquares(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	fmt.Printf("b Width %.1f height %.1f\n", b.Width(), b.Height())
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(lines.LinesFromBBox(b)).WithOffset(0, 0))
	scaleTan := 0.0004
	scaleSin := 0.5
	sampler := samplers.Displace(
		samplers.Mult(
			samplers.Lambda(func(p primitives.Point) float64 {
				return math.Tan(p.X*scaleTan) - math.Tan(p.Y*scaleTan)
			}),
			samplers.Lambda(func(p primitives.Point) float64 {
				return math.Sin(p.X*scaleSin) - math.Cos(p.Y*scaleSin)
			}),
			samplers.HighCenterRelativeDataSource{Scale: .01}),
		primitives.Vector{X: -6000, Y: -4000},
	)

	marchingResolution := 300
	spacing := 0.05
	// baseThreshold := -.2
	curves := getSpacedTruchets(b, marchingResolution, sampler, 0.0, spacing, 15)
	// var curves = []lines.LineLike{}
	// for i := range 15 {
	// 	curves = append(
	// 		curves,
	// 		curve.NewMarchingGrid(b, marchingResolution, sampler, baseThreshold+spacing*float64(i)).GenerateCurves()...,
	// 	)
	// }
	scene = scene.AddLayer(NewLayer("curve").WithLineLike(curves).WithColor("black").WithWidth(10.0))
	return scene
}

func getThreeColorCircleMarchingSquares(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(lines.LinesFromBBox(b)).WithOffset(0, 0))
	// c1 := samplers.CircleRadius{Center: primitives.Point{X: 5000, Y: 3000}}
	c2 := samplers.ScalarMultiple(
		samplers.PointDistance(primitives.Point{X: 1000, Y: 4000}),
		1.0,
	)
	c3 := samplers.PointDistance(primitives.Point{X: 3000, Y: 4000})
	c4 := samplers.PointDistance(primitives.Point{X: 5000, Y: 4000})
	// c5 := samplers.CircleRadius{Center: primitives.Point{X: 7000, Y: 4000}}
	sampler_cyan := samplers.Add(
		c2, samplers.ScalarMultiple(c3, -1.0), //c4,
	)
	sampler_magenta := samplers.Add(
		samplers.ScalarMultiple(c2, -1.0), c4, //c5,
	)
	sampler_yellow := samplers.Add(
		c3, samplers.ScalarMultiple(c4, -1.0), //c5,
	)

	marchingResolution := 150
	spacing := 18.0
	baseThreshold := 1500.0
	var curves_cyan = []lines.LineLike{}
	var curves_magenta = []lines.LineLike{}
	var curves_yellow = []lines.LineLike{}
	for i := range 50 {
		curves_cyan = append(curves_cyan, curve.NewMarchingGrid(b, marchingResolution, sampler_cyan, baseThreshold+spacing*float64(i)).GenerateCurves()...)
		curves_magenta = append(curves_magenta, curve.NewMarchingGrid(b, marchingResolution, sampler_magenta, baseThreshold+spacing*float64(i)).GenerateCurves()...)
		curves_yellow = append(curves_yellow, curve.NewMarchingGrid(b, marchingResolution, sampler_yellow, baseThreshold+spacing*float64(i)).GenerateCurves()...)

	}
	scene = scene.AddLayer(NewLayer("curve-cyan").WithLineLike(curves_cyan).WithColor("cyan").WithWidth(10.0))
	scene = scene.AddLayer(NewLayer("curve-magenta").WithLineLike(curves_magenta).WithColor("magenta").WithWidth(10.0))
	scene = scene.AddLayer(NewLayer("curve-yellow").WithLineLike(curves_yellow).WithColor("yellow").WithWidth(10.0))
	return scene
}

func getCircleMarchingSquares(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(lines.LinesFromBBox(b)).WithOffset(0, 0))
	// sampler := samplers.CircleRadius{Center: primitives.Point{X: 5000, Y: 5000}}
	sampler := samplers.Min(
		samplers.PointDistance(primitives.Point{X: 5000, Y: 5000}),
		samplers.PointDistance(primitives.Point{X: 3000, Y: 4000}),
		samplers.PointDistance(primitives.Point{X: 2500, Y: 5000}),
	)
	marchingResolution := 200
	marchingGrid1 := curve.NewMarchingGrid(b, marchingResolution, sampler, 1103)
	curves1 := marchingGrid1.GenerateCurves()
	// control1 := marchingGrid1.GetControlPoints()
	marchingGrid2 := curve.NewMarchingGrid(b, marchingResolution, sampler, 1115)
	curves2 := marchingGrid2.GenerateCurves()

	scene = scene.AddLayer(NewLayer("curve1").WithLineLike(curves1).WithColor("red").WithWidth(10.0))
	scene = scene.AddLayer(NewLayer("curve2").WithLineLike(curves2).WithColor("green").WithWidth(10.0))
	return scene
}

func getCircleArtifactMarchingSquares(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(lines.LinesFromBBox(b)).WithOffset(0, 0))
	sampler := samplers.Add(
		samplers.Min(
			samplers.PointDistance(primitives.Point{X: 5000, Y: 3000}),
			samplers.PointDistance(primitives.Point{X: 2000, Y: 3500}),
			samplers.PointDistance(primitives.Point{X: 3000, Y: 2000}),
			samplers.PointDistance(primitives.Point{X: 3000, Y: 5300}),
			samplers.PointDistance(primitives.Point{X: 5443, Y: 5300}),
		),
		samplers.Lambda(
			func(p primitives.Point) float64 {
				return math.Sin(p.X*0.0125)*10 + math.Sin(p.Y*0.0125)*10
				// return 0.0
			},
		),
	)
	marchingResolution := 250
	spacing := 18
	var curves = []lines.LineLike{}
	for i := range 50 {
		curves = append(curves, curve.NewMarchingGrid(b, marchingResolution, sampler, 1080+float64(spacing*i)).GenerateCurves()...)

	}
	scene = scene.AddLayer(NewLayer("curve").WithLineLike(curves).WithColor("black").WithWidth(10.0))
	return scene
}

func getPerlinMarchingSquares(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(lines.LinesFromBBox(b)).WithOffset(0, 0))
	offset := primitives.Vector{X: 10000, Y: 10000}
	sampler := samplers.Add(
		samplers.NewPerlinNoise(0.00001, offset),
		samplers.NewPerlinNoise(0.00005, offset),
		samplers.NewPerlinNoise(0.0001, offset),
		samplers.NewPerlinNoise(0.0005, offset),
	)
	marchingResolution := 250
	spacing := 0.02
	baseThreshold := -.2
	var curves = []lines.LineLike{}
	for i := range 20 {
		curves = append(
			curves,
			curve.NewMarchingGrid(b, marchingResolution, sampler, baseThreshold+spacing*float64(i)).GenerateCurves()...,
		)

	}
	scene = scene.AddLayer(NewLayer("curve").WithLineLike(curves).WithColor("black").WithWidth(10.0))
	return scene
}
