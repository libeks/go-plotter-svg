package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/scenes"
	"github.com/libeks/go-plotter-svg/svg"
)

type Config struct {
	fname     string
	sceneName string
	debug     bool
}

func main() {
	args := os.Args[1:]
	sizePx := 10000.0

	config, err := parseFlags(args)
	if err != nil {
		panic(err)
	}

	outerBox := primitives.BBox{
		UpperLeft:  primitives.Point{X: 0, Y: 800},                        // leave space at the top for guides
		LowerRight: primitives.Point{X: sizePx * (12.0 / 9.0), Y: sizePx}, // make sure it spans the 9"x12" canvas
	}
	start := time.Now()
	innerBox := outerBox.WithPadding(500) // enough to no hit the edges
	library := scenes.GatherScenes()
	if config.debug {
		sceneNames := library.GetNames()
		fmt.Printf("These scenes are available:\n")
		for _, name := range sceneNames {
			fmt.Printf("\t%s\n", name)
		}
		return
	}

	sceneFn, err := library.Get(config.sceneName)
	if err != nil {
		panic(err)
	}
	scene := sceneFn(innerBox)

	scene.CalculateStatistics()
	svg.SVG{
		Fname:    config.fname,
		Width:    "12in",
		Height:   "9in",
		Document: scene,
	}.WriteSVG()
	fmt.Printf("Rendering took %s.\n", time.Since(start))
}

func parseFlags(args []string) (Config, error) {
	fname := "gallery/test.svg"
	sceneName := "test-density-v2"
	// n := len(args)
	for len(args) > 0 {
		arg := args[0]
		if arg == "--fname" {
			if len(args) < 2 {
				return Config{}, errors.New("Parameter '--fname' must be followed by a filename")
			}
			fname = args[1]
			args = args[2:]
			continue
		}
		if arg == "--scene" {
			if len(args) < 2 {
				return Config{}, errors.New("Parameter '--scene' must be followed by a scene name")
			}
			sceneName = args[1]
			args = args[2:]
			continue
		}
		if arg == "--list-scenes" {
			return Config{debug: true}, nil
		}
		if arg == "--help" {
			//TODO: write a help doc ehre
		}
		return Config{}, errors.New(fmt.Sprintf("Not sure what to do with parameters %v", args))
	}
	return Config{
		fname:     fname,
		sceneName: sceneName,
	}, nil
}
