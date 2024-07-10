package main

import (
	"math"
	"math/rand"

	svg "github.com/ajstarks/svgo"
)

func main() {
	fname := "gallery/test1.svg"
	sizePx := 10000.0
	padding := 1000.0

	outerBox := Box{0, 0, sizePx, sizePx}
	innerBox := outerBox.WithPadding(padding)
	// scene := getCurlyScene(outerBox)
	// scene := getLinesInsideScene(innerBox, 1000)
	// scene := getLineFieldInObjects(innerBox)
	// scene := radialBoxScene(innerBox)
	// scene := parallelBoxScene(innerBox)
	// scene := parallelSineFieldsScene(innerBox)
	// scene := ParallelCoherentScene(innerBox)
	// scene := CirclesInSquareScene(innerBox)
	// scene := TestDensityScene(innerBox)
	scene := TruchetScene(innerBox)
	SVG{fname: fname,
		width:  "12in",
		height: "9in",
		Scene:  scene,
	}.WriteSVG()
}

func getCurlyBrush(box Box, width, angle float64) []LineLike {
	brushWidth := width
	path := CurlyFill{
		box:     box.WithPadding(brushWidth),
		angle:   angle,
		spacing: float64(brushWidth),
	}
	// return []LineLike{Path{path.GetPath()}}
	return []LineLike{path.GetPath()}
}

func getLinesInsideScene(box Box, n int) Scene {
	// poly := Polygon{[]Point{
	// 	{3000, 3000},
	// 	{7000, 3000},
	// 	{7000, 7000},
	// 	{3000, 7000},
	// }}
	poly := Circle{
		Point{5000, 5000},
		1000,
	}
	return getLinesInsidePolygonScene(box, poly, n)
}

func randRangeMinusPlusOne() float64 {
	return 2 * (rand.Float64() - 0.5)
}

func randInRange(min, max float64) float64 {
	return (max-min)*rand.Float64() + min
}

func randomlyAllocateSegments(segments [][]LineLike, threshold float64) ([]LineLike, []LineLike) {
	layer1 := []LineLike{}
	layer2 := []LineLike{}
	for _, segs := range segments {
		if rand.Float64() > threshold {
			layer1 = append(layer1, segs...)
		} else {
			layer2 = append(layer2, segs...)
		}
	}
	return layer1, layer2
}

type SineDensity struct {
	min    float64
	max    float64
	offset float64
	cycles float64
}

func (d SineDensity) Density(a float64) float64 {
	theta := d.cycles * (a + d.offset) * math.Pi
	dRange := d.max - d.min
	return d.min + dRange*(math.Sin(theta)+1)/2
}

func radialBoxWithCircleExclusion(container Object, center Point, nLines int, radius float64) []LineLike {
	radial := CircularLineField(nLines, center)
	compObject := NewComposite().With(container).Without(Circle{center, radius})
	lines := limitLinesToShape(radial, compObject)
	segments := segmentsToLineLikes(lines)
	return segments
}

func segmentsToLineLikes(segments []LineSegment) []LineLike {
	linelikes := make([]LineLike, len(segments))
	for i, seg := range segments {
		linelikes[i] = seg
	}
	return linelikes
}

// func circlesToLineLikes(circles []Circle) []LineLike {
// 	linelikes := make([]LineLike, len(circles))
// 	for i, seg := range circles {
// 		linelikes[i] = seg
// 	}
// 	return linelikes
// }

type PlotImage interface {
	Render(*svg.SVG)         // render non-guideline layers (layesr 1+)
	DrawGuideLines(*svg.SVG) // draw guidelines (layer 0)
	GetDefs(*svg.SVG)
}
