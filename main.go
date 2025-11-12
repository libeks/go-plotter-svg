package main

import (
	"fmt"
	"os"
	"time"

	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/scenes"
	"github.com/libeks/go-plotter-svg/svg"
)

func main() {
	args := os.Args[1:]
	fname := "gallery/test.svg"
	sizePx := 10000.0

	if len(args) > 0 {
		fname = args[0]
	}

	outerBox := primitives.BBox{
		UpperLeft:  primitives.Point{X: 0, Y: 800},                        // leave space at the top for guides
		LowerRight: primitives.Point{X: sizePx * (12.0 / 9.0), Y: sizePx}, // make sure it spans the 9"x12" canvas
	}
	start := time.Now()
	innerBox := outerBox.WithPadding(500) // enough to no hit the edges
	fmt.Printf("InnerBox %v\n", innerBox)
	// scene := scenes.getCurlyScene(outerBox)
	// scene := scenes.getLinesInsideScene(innerBox, 1000)
	// scene := scenes.getLineFieldInObjects(innerBox)
	// scene := scenes.radialBoxScene(innerBox)
	// scene := scenes.parallelBoxScene(innerBox)
	// scene := scenes.parallelSineFieldsScene(innerBox)
	// scene := scenes.ParallelCoherentScene(innerBox)
	// scene := scenes.CirclesInSquareScene(innerBox)
	// scene := scenes.TestDensityScene(innerBox)
	// scene := scenes.BoxFillScene(innerBox)
	// scene := scenes.TruchetScene(innerBox)
	// scene := scenes.SweepTruchetScene(innerBox)
	// scene := scenes.RisingSunScene(innerBox)
	// scene := scenes.CCircleLineSegments(innerBox)
	// scene := scenes.Font(innerBox)
	// scene := scenes.Text(innerBox)
	// scene := scenes.PolygonBoxScene(innerBox)
	// scene := scenes.FoldableCubeIDScene(innerBox)
	// scene := scenes.FoldableRhombicuboctahedronID(innerBox)
	// scene := scenes.FoldableRhombiSansCorner(innerBox)
	scene := scenes.FoldableVoronoi(innerBox)
	// scene := scenes.FoldableRightTrianglePrismScene(innerBox)
	// scene := scenes.FoldableRightTrianglePrismIDScene(innerBox)
	// scene := scenes.FoldableCutCubeScene(innerBox)
	// scene := scenes.FoldableRhombiSansCorner(innerBox)
	// scene := scenes.MazeScene(innerBox)
	// scene := scenes.RectanglePackingScene(innerBox)

	scene.CalculateStatistics()
	svg.SVG{
		Fname:    fname,
		Width:    "12in",
		Height:   "9in",
		Document: scene,
	}.WriteSVG()
	fmt.Printf("Rendering took %s.\n", time.Since(start))
}
