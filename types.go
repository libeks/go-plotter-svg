package main

import (
	"fmt"
	"strings"

	"github.com/shabbyrobe/xmlwriter"
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
	for i, layer := range layers {
		fmt.Printf("Layer %d has %d objects\n", i, len(layer.linelikes))
	}
	lines := []LineLike{}
	for i := 1; i < len(s.layers); i++ {
		ii := float64(i)

		for j := 300.0; j < 800.0; j += 100.0 {
			lines = append(lines,
				Line{j + ii*1000, 200, j + ii*1000, 300},
			)
		}
		for j := 300.0; j < 800.0; j += 100.0 {
			lines = append(lines,
				Line{j + ii*1000, 700, j + ii*1000, 800},
			)
		}
		for j := 300.0; j < 800.0; j += 100.0 {
			lines = append(lines,
				Line{200 + ii*1000, j, 300 + ii*1000, j},
			)
		}
		for j := 300.0; j < 800.0; j += 100.0 {
			lines = append(lines,
				Line{700 + ii*1000, j, 800 + ii*1000, j},
			)
		}

	}
	layers = append(layers, NewLayer("GUIDELINES-pen").WithLineLike(lines))
	for i := 1; i < len(s.layers); i++ {
		ii := float64(i)
		layers = append(layers, NewLayer(fmt.Sprintf("GUIDELINES-Layer %d", i)).WithLineLike([]LineLike{
			Line{500.0 + ii*1000, 300.0, 500 + ii*1000, 700},
			Line{300 + ii*1000, 500.0, 700 + ii*1000, 500},
		}).WithColor(layers[i].color).WithWidth(layers[i].width).WithOffset(layers[i].offsetX, layers[i].offsetY))
	}
	for i, layer := range layers {
		fmt.Printf("layer[%d]: %v\n", i, layers[i])
		fmt.Printf("Layer %d has %d objects\n", i, len(layer.linelikes))

	}
	return layers
}

// implemented by Line, Path
type LineLike interface {
	XML(color, width string) xmlwriter.Elem
	String() string
}

type Path struct {
	s string
}

func (p Path) String() string {
	return fmt.Sprintf("Path (%s)", p.s)
}

func (p Path) XML(color, width string) xmlwriter.Elem {
	return xmlwriter.Elem{
		Name: "path", Attrs: []xmlwriter.Attr{
			{
				Name:  "d",
				Value: p.s,
			},
			{
				Name:  "stroke",
				Value: color,
			},
			{
				Name:  "fill",
				Value: "none",
			},
			{
				Name:  "stroke-width",
				Value: width,
			},
		},
	}
}

type Line struct {
	x1 float64
	y1 float64
	x2 float64
	y2 float64
}

func (l Line) String() string {
	return fmt.Sprintf("Line (%.1f, %.1f) -> (%.1f, %.1f)", l.x1, l.y1, l.x2, l.y2)
}

func (l Line) Reverse() Line {
	return Line{l.x2, l.y2, l.x1, l.y1}
}

func (l Line) XML(color, width string) xmlwriter.Elem {
	return xmlwriter.Elem{
		Name: "line", Attrs: []xmlwriter.Attr{
			{
				Name:  "x1",
				Value: fmt.Sprintf("%.1f", l.x1),
			},
			{
				Name:  "x2",
				Value: fmt.Sprintf("%.1f", l.x2),
			},
			{
				Name:  "y1",
				Value: fmt.Sprintf("%.1f", l.y1),
			},
			{
				Name:  "y2",
				Value: fmt.Sprintf("%.1f", l.y2),
			},
			{
				Name:  "stroke",
				Value: color,
			},
			{
				Name:  "fill",
				Value: "none",
			},
			{
				Name:  "stroke-width",
				Value: width,
			},
		},
	}
}

func NewLayer(annotation string) Layer {
	return Layer{name: annotation}
}

type Layer struct {
	name      string
	linelikes []LineLike
	offsetX   float64
	offsetY   float64
	color     string
	width     float64
}

func (l Layer) WithLineLike(linelikes []LineLike) Layer {
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

type Box struct {
	x    float64
	y    float64
	xEnd float64
	yEnd float64
}

func (b Box) String() string {
	return fmt.Sprintf("Box (%d, %d) -> (%d, %d)", b.x, b.y, b.xEnd, b.yEnd)
}

func (b Box) Lines() []LineLike {
	commands := []string{}
	// find the starting point - extreme point of box in direction perpendicular to

	commands = append(commands, fmt.Sprintf("M %.3f %.3f", b.x, b.y))

	commands = append(commands, fmt.Sprintf("L %.3f %.3f", b.x, float64(b.yEnd)))
	commands = append(commands, fmt.Sprintf("L %.3f %.3f", b.xEnd, float64(b.yEnd)))
	commands = append(commands, fmt.Sprintf("L %.3f %.3f", b.xEnd, float64(b.y)))
	commands = append(commands, fmt.Sprintf("L %.3f %.3f", b.x, float64(b.y)))

	return []LineLike{
		Path{
			strings.Join(commands, " "),
		},
	}
}

func (b Box) WithPadding(pad float64) Box {
	return Box{
		b.x + pad,
		b.y + pad,
		b.xEnd - pad,
		b.yEnd - pad,
	}
}
