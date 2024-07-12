package lines

import (
	"fmt"
	"strings"

	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/primitives"
)

func NewPath(start primitives.Point) Path {
	return Path{
		start: start,
	}
}

type Path struct {
	start  primitives.Point
	chunks []PathChunk
}

func (p Path) AddPathChunk(chunk PathChunk) Path {
	p.chunks = append(p.chunks, chunk)
	return p
}

func (p Path) Len() float64 {
	total := 0.0
	start := p.start
	for _, chunk := range p.chunks {
		total += chunk.Length(start)
		start = chunk.Endpoint()
	}
	return total
}

func (p Path) Start() primitives.Point {
	return p.start
}

func (p Path) End() primitives.Point {
	if len(p.chunks) > 0 {
		return p.chunks[len(p.chunks)-1].Endpoint()
	}
	return p.start
}

func (p Path) String() string {
	return fmt.Sprintf("Path (%s)", p.pathString())
}

func (p Path) pathString() string {
	strs := []string{fmt.Sprintf("M %.1f %.1f", p.start.X, p.start.Y)}
	for _, xml := range p.chunks {
		strs = append(strs, xml.XMLChunk())
	}
	return strings.Join(strs, " ")
}

func (p Path) IsEmpty() bool {
	return len(p.chunks) == 0
}

func (p Path) XML(color, width string) xmlwriter.Elem {
	return xmlwriter.Elem{
		Name: "path", Attrs: []xmlwriter.Attr{
			{
				Name:  "d",
				Value: p.pathString(),
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
