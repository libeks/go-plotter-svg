package scenes

import (
	"fmt"
	"time"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/primitives"
)

type Scene struct {
	Layers []Layer
	Guides bool
}

func (s Scene) AddLayer(layer Layer) Scene {
	s.Layers = append(s.Layers, layer)
	return s
}

func (s Scene) WithGuides() Scene {
	s.Guides = true
	return s
}

func (s Scene) GetLayers() []Layer {
	if !s.Guides || len(s.Layers) < 2 {
		return s.Layers
	}
	// draw guides on the upper edge of the image
	// assume that the 0th layer contains the guidelines
	layers := s.Layers
	ls := []lines.LineLike{}
	increment := 25.0
	for i := 1; i < len(s.Layers); i++ {
		ii := float64(i)

		for j := 300.0; j <= 700.0; j += increment {
			len := 75.0
			if j == 500 {
				len = 100.0
			}
			ls = append(ls,
				lines.LineSegment{P1: primitives.Point{X: j + ii*1000, Y: 300 - len}, P2: primitives.Point{X: j + ii*1000, Y: 300}},
			)
		}
		for j := 300.0; j <= 700.0; j += increment {
			len := 75.0
			if j == 500 {
				len = 100.0
			}
			ls = append(ls,
				lines.LineSegment{P1: primitives.Point{X: j + ii*1000, Y: 700}, P2: primitives.Point{X: j + ii*1000, Y: 700 + len}},
			)
		}
		for j := 300.0; j <= 700.0; j += increment {
			len := 75.0
			if j == 500 {
				len = 100.0
			}
			ls = append(ls,
				lines.LineSegment{P1: primitives.Point{X: 300 - len + ii*1000, Y: j}, P2: primitives.Point{X: 300 + ii*1000, Y: j}},
			)
		}
		for j := 300.0; j <= 700.0; j += increment {
			len := 75.0
			if j == 500 {
				len = 100.0
			}
			ls = append(ls,
				lines.LineSegment{P1: primitives.Point{X: 700 + ii*1000, Y: j}, P2: primitives.Point{X: 700 + len + ii*1000, Y: j}},
			)
		}

	}
	layers = append(layers, NewLayer("GUIDELINES-pen").WithLineLike(ls))
	for i := 1; i < len(s.Layers); i++ {
		ii := float64(i)
		layers = append(layers, NewLayer(fmt.Sprintf("GUIDELINES-Layer %d", i)).WithLineLike([]lines.LineLike{
			lines.LineSegment{P1: primitives.Point{X: 500.0 + ii*1000, Y: 300.0}, P2: primitives.Point{X: 500 + ii*1000, Y: 700}},
			lines.LineSegment{P1: primitives.Point{X: 300 + ii*1000, Y: 500.0}, P2: primitives.Point{X: 700 + ii*1000, Y: 500}},
		}).WithColor(layers[i].color).WithWidth(layers[i].width).WithOffset(layers[i].offsetX, layers[i].offsetY))
	}
	return layers
}

func (s Scene) CalculateStatistics() {
	yesGuides := "without"
	if s.Guides {
		yesGuides = "with"
	}

	fmt.Printf("Scene has %d layers, %s guides\n", len(s.Layers), yesGuides)
	for i, layer := range s.Layers {
		lengths := []float64{}
		upDistances := []float64{}
		start := primitives.Origin
		for _, linelike := range layer.linelikes {
			// fmt.Printf("line %s has len %.1f\n", linelike, linelike.Len())
			lengths = append(lengths, linelike.Len())
			end := linelike.End()
			upDistances = append(upDistances, end.Subtract(start).Len())
			start = end
		}
		end := primitives.Origin
		upDistances = append(upDistances, end.Subtract(start).Len())
		downLen := imageSpaceToMeters(maths.SumFloats(lengths))
		upLen := imageSpaceToMeters(maths.SumFloats(upDistances))
		totalDistance := downLen + upLen
		timeEstimate := metersToTime(totalDistance)
		fmt.Printf("layer %d has %d curves, down distance %.1fm, up distance %.1fm, total %.1fm traveled\n", i, len(layer.linelikes), downLen, upLen, totalDistance)
		fmt.Printf("Would take about %s to plot\n", timeToMinSec(timeEstimate))
	}
}

func (s Scene) OptimizeLines(flipCurves bool) Scene {
	// TODO: Actually fill in
	return s
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
