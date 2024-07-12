package main

import (
	"fmt"

	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

type Scene struct {
	layers []Layer
	guides bool
}

func (s Scene) AddLayer(layer Layer) Scene {
	s.layers = append(s.layers, layer)
	return s
}

func (s Scene) WithGuides() Scene {
	s.guides = true
	return s
}

func (s Scene) Layers() []Layer {
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

func NewLayer(annotation string) Layer {
	return Layer{name: annotation}
}

type Layer struct {
	name      string
	linelikes []lines.LineLike
	offsetX   float64
	offsetY   float64
	color     string
	width     float64
}

func (l Layer) WithLineLike(linelikes []lines.LineLike) Layer {
	l.linelikes = append(l.linelikes, linelikes...)
	return l
}

func (l Layer) WithOffset(x, y float64) Layer {
	l.offsetX = x
	l.offsetY = y
	return l
}

func (l Layer) WithColor(color string) Layer {
	l.color = color
	return l
}

func (l Layer) WithWidth(width float64) Layer {
	l.width = width
	return l
}

func (l Layer) String() string {
	return fmt.Sprintf("Layer %s %v", l.name, l.linelikes)
}

func (l Layer) XML(i int) xmlwriter.Elem {
	color := "black"
	if l.color != "" {
		color = l.color
	}
	width := "3"
	if l.width > 0 {
		width = fmt.Sprintf("%.1f", l.width)
	}
	contents := []xmlwriter.Writable{}
	for _, line := range l.linelikes {
		contents = append(contents, line.XML(color, width))
	}
	return xmlwriter.Elem{
		Name: "g", Attrs: []xmlwriter.Attr{
			{Name: "inkscape:groupmode", Value: "layer"},
			{Name: "inkscape:label", Value: fmt.Sprintf("%d - %s", i, l.name)},
			{Name: "id", Value: "g5"},
			{Name: "transform", Value: fmt.Sprintf("translate(%.1f %.1f)", l.offsetX, l.offsetY)}, // no translation for now
		},
		Content: contents,
	}
}
