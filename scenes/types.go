package scenes

import (
	"fmt"
	"time"

	"github.com/libeks/go-plotter-svg/foldable"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

type Document struct {
	pages  []Page
	guides bool
}

func (d Document) WithGuides() Document {
	d.guides = true
	return d
}

func (d Document) NumPages() int {
	return len(d.pages)
}

// Adds a layer to the first page of the document, this is here for backwards compatibility
func (d Document) AddLayer(layer Layer) Document {
	if len(d.pages) == 0 {
		d.pages = []Page{
			Page{}.AddLayer(layer),
		}
	} else {
		d.pages[0] = d.pages[0].AddLayer(layer)
	}
	return d
}

func (d Document) AddPage(p Page) Document {
	d.pages = append(d.pages, p)
	return d
}

func (d Document) CalculateStatistics() {
	fmt.Printf("Document contains %d pages\n", len(d.pages))
	for i, page := range d.pages {
		fmt.Printf("Page #%d:\n", i)
		page.CalculateStatistics()
	}
}

func (d Document) Page(i int) Page {
	if i > len(d.pages)-1 {
		panic(fmt.Sprintf("Document only has %d pages", len(d.pages)))
	}
	return d.pages[i]
}

func (s Document) AddFoldableLayers(p foldable.FoldablePattern) Document {
	// TODO: Customize brush width rendering
	polygons := []lines.LineLike{}
	for _, poly := range p.Polygons {
		polygons = append(polygons, segmentsToLineLikes(poly.EdgeLines())...)
	}
	s = s.AddLayer(NewLayer("black").WithLineLike(p.Edges).WithColor("black").WithWidth(20).MinimizePath(true))
	s = s.AddLayer(NewLayer("red").WithLineLike(polygons).WithColor("red").WithWidth(20).MinimizePath(true))
	s = s.AddLayer(NewLayer("green").WithLineLike(p.Annotations).WithColor("green").WithWidth(20).MinimizePath(true))
	s = s.AddLayer(NewLayer("blue").WithLineLike(lines.LinesFromBBox(p.BBox())).WithColor("blue").WithWidth(20).MinimizePath(true))
	for color, infill := range p.Fill {
		layerName := fmt.Sprintf("fill-%s", color)
		s = s.AddLayer(NewLayer(layerName).WithLineLike(infill).WithColor(color).WithWidth(20).MinimizePath(true))
	}
	return s
}

type Page struct {
	layers []Layer
	guides bool
}

func (s Page) AddLayer(layer Layer) Page {
	s.layers = append(s.layers, layer)
	return s
}

func (s Page) WithGuides() Page {
	s.guides = true
	return s
}

func (s Page) GetLayers() []Layer {
	if !s.guides || len(s.layers) < 2 {
		return s.layers
	}
	// draw guides on the upper edge of the image
	// assume that the 0th layer contains the guidelines
	layers := s.layers
	ls := []lines.LineLike{}
	increment := 25.0
	for i := 1; i < len(s.layers); i++ {
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
	for i := 1; i < len(s.layers); i++ {
		ii := float64(i)
		layers = append(layers, NewLayer(fmt.Sprintf("GUIDELINES-Layer %d", i)).WithLineLike([]lines.LineLike{
			lines.LineSegment{P1: primitives.Point{X: 500.0 + ii*1000, Y: 300.0}, P2: primitives.Point{X: 500 + ii*1000, Y: 700}},
			lines.LineSegment{P1: primitives.Point{X: 300 + ii*1000, Y: 500.0}, P2: primitives.Point{X: 700 + ii*1000, Y: 500}},
		}).WithColor(layers[i].color).WithWidth(layers[i].width).WithOffset(layers[i].offsetX, layers[i].offsetY))
	}
	return layers
}

func (s Page) CalculateStatistics() {
	yesGuides := "without"
	if s.guides {
		yesGuides = "with"
	}

	fmt.Printf("Page has %d layers, %s guides\n", len(s.layers), yesGuides)
	for i, layer := range s.layers {
		fmt.Printf("layer '%s' #%d has %s\n", layer.name, i, layer.Statistics())
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
