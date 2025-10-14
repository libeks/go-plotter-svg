package main

import (
	"os"

	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/scenes"
	"github.com/libeks/go-plotter-svg/svg"
)

func main() {
	args := os.Args[1:]
	fname := "gallery/test1.svg"
	sizePx := 10000.0
	padding := 1000.0

	if len(args) > 0 {
		fname = args[0]
	}

	outerBox := primitives.BBox{
		UpperLeft:  primitives.Point{X: 0, Y: 0},
		LowerRight: primitives.Point{X: sizePx, Y: sizePx},
	}
	innerBox := outerBox.WithPadding(padding)
	// scene := scenes.getCurlyScene(outerBox)
	// scene := scenes.getLinesInsideScene(innerBox, 1000)
	// scene := scenes.getLineFieldInObjects(innerBox)
	// scene := scenes.radialBoxScene(innerBox)
	// scene := scenes.parallelBoxScene(innerBox)
	// scene := scenes.parallelSineFieldsScene(innerBox)
	// scene := scenes.ParallelCoherentScene(innerBox)
	// scene := scenes.CirclesInSquareScene(innerBox)
	// scene := scenes.TestDensityScene(innerBox)
	// scene := scenes.TruchetScene(innerBox)
	// scene := scenes.SweepTruchetScene(innerBox)
	// scene := scenes.RisingSunScene(innerBox)
	// scene := scenes.CCircleLineSegments(innerBox)
	// scene := scenes.Font(innerBox)
	// scene := scenes.Text(innerBox)
	scene := scenes.PolygonBoxScene(innerBox)
	// scene := scenes.FoldableCubeIDScene(innerBox)
	// scene := scenes.FoldableRhombicuboctahedronID(innerBox)
	// scene := scenes.FoldableRhombiSansCorner(innerBox)
	// scene := scenes.FoldableRightTrianglePrismScene(innerBox)
	// scene := scenes.FoldableRightTrianglePrismIDScene(innerBox)
	// scene := scenes.FoldableCutCubeScene(innerBox)
	// scene := scenes.FoldableRhombiSansCorner(innerBox)
	// scene := scenes.MazeScene(innerBox)
	flipCurves := false
	scene.OptimizeLines(flipCurves)
	scene.CalculateStatistics()
	svg.SVG{
		Fname:  fname,
		Width:  "12in",
		Height: "9in",
		Scene:  scene,
	}.WriteSVG()
}
