package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

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
	calculateStatistics(scene)
	SVG{fname: fname,
		width:  "12in",
		height: "9in",
		Scene:  scene,
	}.WriteSVG()
}

func calculateStatistics(scene Scene) {
	yesGuides := "no"
	if scene.guides {
		yesGuides = "with "
	}

	fmt.Printf("Scene has %d layers, %s guides\n", len(scene.layers), yesGuides)
	for i, layer := range scene.layers {
		lengths := []float64{}
		upDistances := []float64{}
		start := Point{0, 0}
		for _, linelike := range layer.linelikes {
			lengths = append(lengths, linelike.Len())
			end := linelike.End()
			upDistances = append(upDistances, end.Subtract(start).Len())
			start = end
		}
		end := Point{0, 0}
		upDistances = append(upDistances, end.Subtract(start).Len())
		downLen := imageSpaceToMeters(sumFloats(lengths))
		upLen := imageSpaceToMeters(sumFloats(upDistances))
		totalDistance := downLen + upLen
		timeEstimate := metersToTime(totalDistance)
		fmt.Printf("layer %d has %d curves, down distance %.1fm, up distance %.1fm, total %.1fm traveled\n", i, len(layer.linelikes), downLen, upLen, totalDistance)
		fmt.Printf("Would take about %s to plot\n", timeToMinSec(timeEstimate))
	}
}

func timeToMinSec(d time.Duration) string {
	minutes := int(d / time.Minute)
	seconds := int((d - time.Duration(float64(minutes)*float64(time.Minute))) / time.Second)
	return fmt.Sprintf("%dm%ds", minutes, seconds)
}

func metersToTime(m float64) time.Duration {
	return time.Duration(22.6 * float64(time.Second) * m)
}

func imageSpaceToMeters(l float64) float64 {
	const unitPerMeter = 44092.0
	return l / unitPerMeter
}

func sumFloats(l []float64) float64 {
	total := 0.0
	for _, v := range l {
		total += v
	}
	return total
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
