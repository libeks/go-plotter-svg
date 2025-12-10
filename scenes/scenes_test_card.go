package scenes

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/collections"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

func testDensityScene(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	quarters := primitives.PartitionIntoSquares(b, 2)
	colors := []string{
		"red",
		"green",
		"blue",
		"orange",
	}
	boxes := lines.LinesFromBBox(b)

	layers := []Layer{}
	for i, quarter := range quarters.BoxIterator() {
		lineLikes := []lines.LineLike{}
		quarterBox := quarter.WithPadding(50)
		testBoxes := primitives.PartitionIntoSquares(quarterBox, 5)
		for _, tbox := range testBoxes.BoxIterator() {
			jj := float64(tbox.J)
			ii := float64(tbox.I)
			var ls []lines.LineLike
			tboxx := tbox.WithPadding(50)
			boxes = append(boxes, lines.LinesFromBBox(tboxx)...)
			spacing := (jj + 1) * 10
			if tbox.I < 4 {
				ls = lines.SegmentsToLineLikes(
					collections.LimitLinesToShape(
						collections.LinearLineField(
							tboxx, ii*math.Pi/4, spacing,
						),
						objects.PolygonFromBBox(tboxx),
					),
				)
			} else {
				ls = collections.LimitCirclesToShape(
					collections.ConcentricCircles(
						tboxx, tboxx.Center(), spacing,
					),
					objects.PolygonFromBBox(tboxx),
				)
			}
			lineLikes = append(lineLikes, ls...)
		}
		layerName := fmt.Sprintf("pen %d", i)
		layers = append(layers, NewLayer(layerName).WithLineLike(lineLikes).WithOffset(0, 0).WithColor(colors[i]))

	}
	// ensure that the frame layer is rendered first, to make sure it isn't added to guides
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(boxes).WithOffset(0, 0).WithColor("grey"))
	for _, layer := range layers {
		scene = scene.AddLayer(layer)
	}
	return scene
}
