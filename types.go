package main

import (
	"fmt"
	"math"
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
	lines := []LineLike{}
	for i := 1; i < len(s.layers); i++ {
		ii := float64(i)

		for j := 300.0; j < 800.0; j += 100.0 {
			lines = append(lines,
				LineSegment{j + ii*1000, 200, j + ii*1000, 300},
			)
		}
		for j := 300.0; j < 800.0; j += 100.0 {
			lines = append(lines,
				LineSegment{j + ii*1000, 700, j + ii*1000, 800},
			)
		}
		for j := 300.0; j < 800.0; j += 100.0 {
			lines = append(lines,
				LineSegment{200 + ii*1000, j, 300 + ii*1000, j},
			)
		}
		for j := 300.0; j < 800.0; j += 100.0 {
			lines = append(lines,
				LineSegment{700 + ii*1000, j, 800 + ii*1000, j},
			)
		}

	}
	layers = append(layers, NewLayer("GUIDELINES-pen").WithLineLike(lines))
	for i := 1; i < len(s.layers); i++ {
		ii := float64(i)
		layers = append(layers, NewLayer(fmt.Sprintf("GUIDELINES-Layer %d", i)).WithLineLike([]LineLike{
			LineSegment{500.0 + ii*1000, 300.0, 500 + ii*1000, 700},
			LineSegment{300 + ii*1000, 500.0, 700 + ii*1000, 500},
		}).WithColor(layers[i].color).WithWidth(layers[i].width).WithOffset(layers[i].offsetX, layers[i].offsetY))
	}
	return layers
}

// implemented by LineSegment, Path
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

type Vector struct {
	x float64
	y float64
}

func (v Vector) String() string {
	return fmt.Sprintf("Vector (%.1f, %.1f)", v.x, v.y)
}

func (v Vector) Mult(t float64) Vector {
	return Vector{t * v.x, t * v.y}
}

func (v Vector) Add(w Vector) Vector {
	return Vector{v.x + w.x, v.y + w.y}
}

func (v Vector) Dot(w Vector) float64 {
	return v.x*w.x + v.y*w.y
}

func (v Vector) Len() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector) Point() Point {
	return Point{v.x, v.y}
}

// Perp returns a vector perpendicular to v of the same lenght,
// rotated counter-clockwise by 90deg
func (v Vector) Perp() Vector {
	return Vector{-v.y, v.x}
}

type Point struct {
	x float64
	y float64
}

func (p Point) String() string {
	return fmt.Sprintf("Point (%.1f, %.1f)", p.x, p.y)
}

func (p Point) Add(v Vector) Point {
	return Point{p.x + v.x, p.y + v.y}
}

func (p Point) Subtract(p2 Point) Vector {
	return Vector{p.x - p2.x, p.y - p2.y}
}

type Line struct {
	p Point
	v Vector
}

func (l Line) String() string {
	return fmt.Sprintf("Line (%s, %s)", l.p, l.v)
}

// Return a point on the line that is t lenghts of v away from p.
func (l Line) At(t float64) Point {
	return l.p.Add(l.v.Mult(t))
}

func (l Line) Intersect(l2 Line) *Point {
	// TODO: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
	// first line is x1,y1 = l.p, x2,y2 = l.p + l.v. so x2-x1 = l.v.x, y2-y1 = l.v.y
	// second line is x3,y3 = l2.p, x4,y4 = l2.p + l2.v, so x4-x3 = l2.v.x, y4-y3 = l2.v.y
	// determinant is (x1-x2)(y3-y4) - (y1-y2)(x3-x4) = (x2-x1)(y4-y3) - (y2-y1)(x4-x3)
	determinant := l.v.x*l2.v.y - l.v.y*l2.v.x
	if determinant == 0 {
		return nil
	}
	// result is
	// x = ((x1*y2 - y1*x2)(x3-x4) - (x1-x2)(x3*y4 - y3*x4))/determinant
	// y = ((x1*y2 - y1*x2)(y3-y4) - (y1-y2)(x3*y4 - y3*x4))/determinant
	x1x2 := -l.v.x
	x3x4 := -l2.v.x
	y1y2 := -l.v.y
	y3y4 := -l2.v.y
	x1y2y1x2 := (l.p.x*(l.p.y+l.v.y) - l.p.y*(l.p.x+l.v.x))
	x3y4y3x4 := (l2.p.x * (l2.p.y + l2.v.y)) * (l2.p.y * (l2.p.x + l2.v.x))

	x := (x1y2y1x2*x3x4 - x1x2*x3y4y3x4) / determinant
	y := (x1y2y1x2*y3y4 - y1y2*x3y4y3x4) / determinant
	return &Point{x, y}
}

// Return the intersect t parameter of the line l when intersecting line l2
func (l Line) IntersectT(l2 Line) *float64 {
	// TODO: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
	// first line is x1,y1 = l.p, x2,y2 = l.p + l.v. so x2-x1 = l.v.x, y2-y1 = l.v.y
	// second line is x3,y3 = l2.p, x4,y4 = l2.p + l2.v, so x4-x3 = l2.v.x, y4-y3 = l2.v.y
	x1x2 := -l.v.x
	x3x4 := -l2.v.x
	y1y2 := -l.v.y
	y3y4 := -l2.v.y
	x1x3 := l.p.x - l2.p.x
	y1y3 := l2.p.y - l2.p.y

	divisor := (x1x2*y3y4 - y1y2*x3x4)
	if divisor == 0.0 {
		return nil
	}
	t := (x1x3*y3y4 - y1y3*x3x4) / divisor
	return &t
}

// Return the intersection parameters t,u for both lines l and l2
func (l Line) IntersectTU(l2 Line) (*float64, *float64) {
	// TODO: https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection
	// first line is x1,y1 = l.p, x2,y2 = l.p + l.v. so x2-x1 = l.v.x, y2-y1 = l.v.y
	// second line is x3,y3 = l2.p, x4,y4 = l2.p + l2.v, so x4-x3 = l2.v.x, y4-y3 = l2.v.y

	x1x2 := -l.v.x
	x3x4 := -l2.v.x
	y1y2 := -l.v.y
	y3y4 := -l2.v.y
	x1x3 := l.p.x - l2.p.x
	y1y3 := l.p.y - l2.p.y

	divisor := (x1x2*y3y4 - y1y2*x3x4)
	if divisor == 0.0 {
		return nil, nil
	}
	t := (x1x3*y3y4 - y1y3*x3x4) / divisor
	u := (x1x2*y1y3 - y1y2*x1x3) / -divisor // note the divisor is negative here. I initially missed that.
	// fmt.Printf("t:%.3f, u: %.3f\n", t, u)
	return &t, &u
}

func (l Line) IntersectLineSegmentT(ls2 LineSegment) *float64 {
	l2 := ls2.Line()
	t, u := l.IntersectTU(l2)
	// fmt.Printf("Intersecting %s with %s, got t: %+v, u: %+v\n", l, ls2, t, u)
	if t == nil || u == nil {
		return nil
	}
	uu := *u
	if uu <= 1.0 && uu >= 0.0 {
		return t
	}
	return nil
}

type LineSegment struct {
	x1 float64
	y1 float64
	x2 float64
	y2 float64
}

func (l LineSegment) String() string {
	return fmt.Sprintf("LineSegment (%.1f, %.1f) -> (%.1f, %.1f)", l.x1, l.y1, l.x2, l.y2)
}

func (l LineSegment) Reverse() LineSegment {
	return LineSegment{l.x2, l.y2, l.x1, l.y1}
}

func (l LineSegment) XML(color, width string) xmlwriter.Elem {
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

func (l LineSegment) Line() Line {
	return Line{
		p: Point{l.x1, l.y1},
		v: Vector{l.x2 - l.x1, l.y2 - l.y1},
	}
}

func (l LineSegment) IntersectLineT(l2 Line) *float64 {
	l1 := l.Line()
	t := l1.IntersectT(l2)
	if t == nil {
		return nil
	}
	tt := *t
	if tt <= 1.0 && tt >= 0.0 {
		return t
	}
	return nil
}

func (l LineSegment) IntersectLineSegmentT(ls2 LineSegment) *float64 {
	l1 := l.Line()
	l2 := ls2.Line()
	t, u := l1.IntersectTU(l2)
	if t == nil || u == nil {
		return nil
	}
	tt := *t
	uu := *u
	if (tt <= 1.0 && tt >= 0.0) && (uu <= 1.0 && uu >= 0.0) {
		return t
	}
	return nil
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
	return fmt.Sprintf("Box (%.1f, %.1f) -> (%.1f, %.1f)", b.x, b.y, b.xEnd, b.yEnd)
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

func (b Box) ClipLineToBox(l Line) *LineSegment {
	ls := []LineSegment{
		{b.x, b.y, b.x, b.yEnd},
		{b.x, b.yEnd, b.xEnd, b.yEnd},
		{b.xEnd, b.yEnd, b.xEnd, b.y},
		{b.xEnd, b.y, b.x, b.y},
	}
	ts := []float64{}
	for _, lineseg := range ls {
		if t := l.IntersectLineSegmentT(lineseg); t != nil {
			ts = append(ts, *t)
		}
	}
	if len(ts) == 0 {
		return nil
	}
	if len(ts) == 2 {
		p1 := l.At(ts[0])
		p2 := l.At(ts[1])
		return &LineSegment{
			p1.x, p1.y, p2.x, p2.y,
		}
	}
	panic(fmt.Errorf("line had weird number of intersections with box: %v", ts))
}

func (b Box) WithPadding(pad float64) Box {
	return Box{
		b.x + pad,
		b.y + pad,
		b.xEnd - pad,
		b.yEnd - pad,
	}
}

func (b Box) Center() Point {
	return Point{b.x + (b.xEnd-b.x)/2, b.y + (b.yEnd-b.y)/2}
}

func (b Box) Width() float64 {
	return b.xEnd - b.x
}

func (b Box) Height() float64 {
	return b.yEnd - b.y
}

func (b Box) AsPolygon() Object {
	return Polygon{
		points: []Point{
			{b.x, b.y},
			{b.xEnd, b.y},
			{b.xEnd, b.yEnd},
			{b.x, b.yEnd},
		},
	}
}
