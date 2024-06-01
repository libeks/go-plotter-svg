package main

import (
	"fmt"

	svg "github.com/ajstarks/svgo"
)

type Path struct {
	s string
}

func (p Path) String() string {
	return fmt.Sprintf("Path (%s)", p.s)
}

type Line struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func (l Line) String() string {
	return fmt.Sprintf("Line (%d, %d) -> (%d,%d)", l.x1, l.y1, l.x2, l.y2)
}

func (l Line) Reverse() Line {
	return Line{l.x2, l.y2, l.x1, l.y1}
}

type Layer struct {
	name  string // name, should start with "(i+1) .*"
	i     int    // render order
	lines []Line
	paths []Path
}

func (l Layer) String() string {
	return fmt.Sprintf("Layer %s %d %v %v", l.name, l.i, l.lines, l.paths)
}

func NewPlotContainer() PlotContainer {
	return PlotContainer{}
}

type PlotContainer struct {
	layers []Layer
	lines  []Line
}

func (c PlotContainer) WithLayers(l ...Layer) PlotContainer {
	return PlotContainer{
		layers: append(c.layers, l...),
		lines:  c.lines,
	}
}

func (c PlotContainer) WithLines(l ...Line) PlotContainer {
	return PlotContainer{
		layers: c.layers,
		lines:  append(c.lines, l...),
	}
}

func (c PlotContainer) Render(canvas *svg.SVG, defs []string) {
	canvas.Group(`inkscape:groupmode="layer"`, fmt.Sprintf(`inkscape:label="%s"`, "0 - Pencil"))
	for _, line := range c.lines {
		canvas.Line(line.x1, line.y1, line.x2, line.y2, pencilStyle)
	}
	canvas.Gend()

	for layerID, layer := range c.layers {
		canvas.Group(`inkscape:groupmode="layer"`, fmt.Sprintf(`inkscape:label="%s"`, layer.name))
		if len(layer.lines) > 0 && len(layer.paths) > 0 {
			panic(fmt.Errorf("Layer %s has both lines and paths", layer))
		}
		for _, line := range layer.lines {

			// add +1 to endpoint x,y coord to ensure line gradient can render
			// vertical/horizontal lines cannot be rendered with a gradient in SVG
			dx := 0
			dy := 0
			if line.x1 == line.x2 {
				dx = 1
			}
			if line.y1 == line.y2 {
				dy = 1
			}
			canvas.Line(line.x1, line.y1, line.x2+dx, line.y2+dy, defs[layerID]) // add +1 to endpoint x coord to ensure line gradient can render
			fmt.Printf("Just wrote line %s\n", line)
		}
		for _, path := range layer.paths {
			canvas.Path(path.s, `stroke="black"`, `fill="none"`, `stroke-width="10"`)
		}
		canvas.Gend()
	}
}

func (c *PlotContainer) GetDefs(canvas *svg.SVG) []string {
	brushStyles := make([]string, len(c.layers))
	for i := range c.layers {
		style := brushStyle(canvas, i)
		brushStyles[i] = style
	}
	return brushStyles
}
