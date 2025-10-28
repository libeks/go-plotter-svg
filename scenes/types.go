package scenes

import (
	"fmt"
	"time"

	"github.com/libeks/go-plotter-svg/foldable"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/pack"
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
			Page{guides: d.guides}.AddLayer(layer),
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

func FromFoldableLayers(shapes []foldable.FoldablePattern, container primitives.BBox) Document {
	// TODO: Customize brush width rendering
	doc := Document{}.WithGuides()
	bboxes := make([]primitives.BBox, len(shapes))
	for i, shape := range shapes {
		box := shape.BBox()
		bboxes[i] = box.Translate(primitives.Origin.Subtract(box.UpperLeft))
	}
	boxPacking := pack.PackOnMultiplePages(bboxes, container, 200)
	shapesByPage := map[int][]foldable.FoldablePattern{}
	for i, v := range boxPacking.Translations {
		shapesByPage[v.Page] = append(shapesByPage[v.Page], shapes[i].Translate(primitives.Origin.Subtract(shapes[i].BBox().UpperLeft).Add(v.Vector)))
	}
	for i := range boxPacking.Pages {
		page := Page{}
		polygons := []lines.LineLike{}
		edges := []lines.LineLike{}
		annotations := []lines.LineLike{}
		bboxLines := []lines.LineLike{}
		fillColors := map[string][]lines.LineLike{}
		objects := shapesByPage[i]
		for _, p := range objects {
			for _, poly := range p.Polygons {
				polygons = append(polygons, lines.SegmentsToLineLikes(poly.EdgeLines())...)
			}
			edges = append(edges, p.Edges...)
			annotations = append(annotations, p.Annotations...)
			bboxLines = append(bboxLines, lines.LinesFromBBox((p.BBox()))...)
			for color, infill := range p.Fill {
				fillColors[color] = append(fillColors[color], infill...)
			}
		}
		page = page.AddLayer(NewLayer("frame").WithLineLike(lines.LinesFromBBox(container)).WithOffset(0, 0))
		page = page.AddLayer(NewLayer("black").WithLineLike(edges).WithColor("black").WithWidth(20).MinimizePath(true))
		page = page.AddLayer(NewLayer("red").WithLineLike(polygons).WithColor("red").WithWidth(20).MinimizePath(true))
		page = page.AddLayer(NewLayer("green").WithLineLike(annotations).WithColor("green").WithWidth(20).MinimizePath(true))
		page = page.AddLayer(NewLayer("blue").WithLineLike(bboxLines).WithColor("blue").WithWidth(20).MinimizePath(true))
		for color, infill := range fillColors {
			layerName := fmt.Sprintf("fill-%s", color)
			page = page.AddLayer(NewLayer(layerName).WithLineLike(infill).WithColor(color).WithWidth(20).MinimizePath(true))
		}
		doc = doc.AddPage(page)
	}
	return doc
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
		offset := primitives.Vector{X: float64(i) * 1000, Y: 0}

		for j := 300.0; j <= 700.0; j += increment {
			len := 75.0
			if j == 500 {
				len = 100.0
			}
			ls = append(ls,
				lines.LineSegment{
					P1: primitives.Point{X: j, Y: 300 - len},
					P2: primitives.Point{X: j, Y: 300},
				}.Translate(offset),
			)
		}
		for j := 300.0; j <= 700.0; j += increment {
			len := 75.0
			if j == 500 {
				len = 100.0
			}
			ls = append(ls,
				lines.LineSegment{
					P1: primitives.Point{X: j, Y: 700},
					P2: primitives.Point{X: j, Y: 700 + len},
				}.Translate(offset),
			)
		}
		for j := 300.0; j <= 700.0; j += increment {
			len := 75.0
			if j == 500 {
				len = 100.0
			}
			ls = append(ls,
				lines.LineSegment{
					P1: primitives.Point{X: 300 - len, Y: j},
					P2: primitives.Point{X: 300, Y: j},
				}.Translate(offset),
			)
		}
		for j := 300.0; j <= 700.0; j += increment {
			len := 75.0
			if j == 500 {
				len = 100.0
			}
			ls = append(ls,
				lines.LineSegment{
					P1: primitives.Point{X: 700, Y: j},
					P2: primitives.Point{X: 700 + len, Y: j},
				}.Translate(offset),
			)
		}
		ls = append(ls,
			lines.LineSegment{
				P1: primitives.Point{X: 450, Y: 500},
				P2: primitives.Point{X: 550, Y: 500},
			}.Translate(offset),
			lines.LineSegment{
				P1: primitives.Point{X: 500, Y: 450},
				P2: primitives.Point{X: 500, Y: 550},
			}.Translate(offset),
		)

	}
	layers = append(layers, NewLayer("GUIDELINES-pen").WithLineLike(ls))
	for i := 1; i < len(s.layers); i++ {
		offset := primitives.Vector{X: float64(i) * 1000, Y: 0}
		layers = append(layers, NewLayer(fmt.Sprintf("GUIDELINES-Layer %d", i)).WithLineLike([]lines.LineLike{
			lines.LineSegment{
				P1: primitives.Point{X: 500, Y: 300},
				P2: primitives.Point{X: 500, Y: 700},
			}.Translate(offset),
			lines.LineSegment{
				P1: primitives.Point{X: 300, Y: 500},
				P2: primitives.Point{X: 700, Y: 500},
			}.Translate(offset),
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
