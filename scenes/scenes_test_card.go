package scenes

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/collections"
	"github.com/libeks/go-plotter-svg/fonts"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

func densityTestCardV1Scene(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	quarters := primitives.PartitionIntoSquares(b.Square(), 2)
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

func densityTestCardV2Scene(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	nStrips := 11
	strips := primitives.PartitionIntoRectangles(b, nStrips, 1)
	colors := []string{
		"red",
		"green",
		"blue",
		"orange",
		"cyan",
		"magenta",
		"brown",
		"purple",
		"pink",
		"gray",
		"indigo",
	}
	boxes := lines.LinesFromBBox(b)

	layers := []Layer{}
	for i, strip := range strips.BoxIterator() {
		lineLikes := []lines.LineLike{}
		quarterBox := strip.WithPadding(50)
		if i == 0 {
			for _, tbox := range primitives.PartitionIntoRectangles(quarterBox, 1, 10).BoxIterator() {
				fmt.Printf("i %d j %d\n", tbox.I, tbox.J)
				ii := float64(tbox.I)
				spacing := (ii + 1) * 5
				lineLikes = append(lineLikes, fonts.RenderText(tbox.BBox, fmt.Sprintf("%.0f", spacing), fonts.WithSize(200)).CharCurves...)
			}
		} else {
			for _, tbox := range primitives.PartitionIntoRectangles(quarterBox, 2, 10).BoxIterator() {
				ii := float64(tbox.I)
				spacing := (ii + 1) * 5
				tboxx := tbox.WithPadding(50).Square()
				boxes = append(boxes, lines.LinesFromBBox(tboxx)...)
				lineLikes = append(lineLikes,
					collections.FillPolygonWithSpacing(objects.PolygonFromBBox(tboxx), spacing, 0.5*float64(tbox.J))...,
				)
			}
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

func penHeightTestCardScene(b primitives.BBox) Document {
	scene := Document{}.WithGuides()
	nStrips := 10
	strips := primitives.PartitionIntoRectangles(b, nStrips, 1)
	circleRadius := 20.0
	nCircles := 10

	stripLines := []lines.LineLike{}
	layers := []Layer{}
	for i, strip := range strips.BoxIterator() {
		stripLines = append(stripLines, lines.LinesFromBBox(strip.BBox)...)
		lineLikes := []lines.LineLike{}
		spacing := strip.Width() / float64(nCircles)
		for j := range nCircles {
			gap := float64(j) * spacing
			lineLikes = append(lineLikes, lines.FullCircle(strip.BBox.UpperLeft.Add(primitives.Vector{X: circleRadius + gap, Y: circleRadius}), circleRadius))
			lineLikes = append(lineLikes, lines.FullCircle(strip.BBox.UpperLeft.Add(primitives.Vector{X: circleRadius + gap, Y: strip.Height() - circleRadius}), circleRadius))
		}
		layerName := fmt.Sprintf("pen %d", i)
		layers = append(layers, NewLayer(layerName).WithLineLike(lineLikes).WithOffset(0, 0))

	}
	// ensure that the frame layer is rendered first, to make sure it isn't added to guides
	scene = scene.AddLayer(NewLayer("frame").WithLineLike(lines.LinesFromBBox(b)).WithOffset(0, 0).WithColor("grey"))
	scene = scene.AddLayer(NewLayer("strips").WithLineLike(stripLines).WithColor("grey"))
	for _, layer := range layers {
		scene = scene.AddLayer(layer)
	}
	return scene
}

//
