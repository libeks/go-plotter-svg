package scenes

import (
	"fmt"

	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/lines"
)

func NewLayer(annotation string) Layer {
	return Layer{name: annotation}
}

type Layer struct {
	name         string
	linelikes    []lines.LineLike
	controllines []lines.LineLike
	offsetX      float64
	offsetY      float64
	color        string
	width        float64
}

func (l Layer) WithLineLike(linelikes []lines.LineLike) Layer {
	l.linelikes = append(l.linelikes, linelikes...)
	return l
}

func (l Layer) WithControlLines(linelikes []lines.LineLike) Layer {
	l.controllines = append(l.controllines, linelikes...)
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
	for _, line := range l.controllines {
		contents = append(contents, line.GuideXML(color, width))
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
