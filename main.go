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
	sceneName := "test-density-v2"

	if len(args) > 0 {
		fname = args[0]
		args = args[1:]
	}
	if len(args) > 0 {
		sceneName = args[0]
	}

	outerBox := primitives.BBox{
		UpperLeft:  primitives.Point{X: 0, Y: 800},                        // leave space at the top for guides
		LowerRight: primitives.Point{X: sizePx * (12.0 / 9.0), Y: sizePx}, // make sure it spans the 9"x12" canvas
	}
	start := time.Now()
	innerBox := outerBox.WithPadding(500) // enough to no hit the edges
	library := scenes.GatherScenes()
	sceneFn, err := library.Get(sceneName)
	if err != nil {
		panic(err)
	}
	scene := sceneFn(innerBox)

	scene.CalculateStatistics()
	svg.SVG{
		Fname:    fname,
		Width:    "12in",
		Height:   "9in",
		Document: scene,
	}.WriteSVG()
	fmt.Printf("Rendering took %s.\n", time.Since(start))
}
