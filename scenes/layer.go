package scenes

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/primitives"
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

func (l Layer) Statistics() string {
	lengths := []float64{}
	upDistances := []float64{}
	start := primitives.Origin
	for _, linelike := range l.linelikes {
		lengths = append(lengths, linelike.Len())

		upDistances = append(upDistances, start.Subtract(linelike.Start()).Len())
		start = linelike.End()
	}

	upDistances = append(upDistances, start.Subtract(primitives.Origin).Len())
	downLen := imageSpaceToMeters(maths.SumFloats(lengths))
	upLen := imageSpaceToMeters(maths.SumFloats(upDistances))
	totalDistance := downLen + upLen
	timeEstimate := metersToTime(totalDistance) + upDownEstimate(len(l.linelikes))
	return fmt.Sprintf("%d curves, down distance %.1fm, up distance %.1fm, total %.1fm traveled\nWould take about %s to plot", len(l.linelikes), downLen, upLen, totalDistance, timeToMinSec(timeEstimate))
}

func (l Layer) RandomizedClosedCurves() Layer {
	if len(l.linelikes) == 0 {
		return l
	}
	for i, line := range l.linelikes {
		if line.Start().Subtract(line.End()).Len() < 0.1 {
			// line is a closed curve
			start, end := line.Bisect(rand.Float64())
			l.linelikes[i] = end.Join(start) // FIXME: I can't combine two linelikes, but I could if I got the underlying PathChunks...
		}
	}
	return l
}

func (l Layer) MinimizePath(allowReverse bool) Layer {
	fmt.Printf("Layer before %s\n", l.Statistics())
	if len(l.linelikes) == 0 {
		return l
	}
	lns := []lines.LineLike{}
	var inputLines []lines.LineLike
	inputLines = append(inputLines, l.linelikes...)
	lns = append(lns, inputLines[0])
	pt := inputLines[0].End()
	inputLines = inputLines[1:]
	for len(inputLines) > 0 {
		minDistance := math.MaxFloat64
		var minIndex int
		var isReversed bool

		for i, line := range inputLines {
			dist := pt.Subtract(line.Start()).Len()
			if dist < minDistance {
				minDistance = dist
				isReversed = false
				minIndex = i
			}
			if allowReverse {
				dist := pt.Subtract(line.End()).Len()
				if dist < minDistance {
					minDistance = dist
					isReversed = true
					minIndex = i
				}
			}
		}
		line := inputLines[minIndex]
		if isReversed {
			line = line.Reverse()
		}
		lns = append(lns, line)
		inputLines[minIndex] = inputLines[len(inputLines)-1]
		inputLines = inputLines[:len(inputLines)-1]
		pt = line.End()
	}
	if len(l.linelikes) != len(lns) {
		panic(fmt.Errorf("did not get the same number of lines %s vs %s", l.linelikes, lns))
	}
	l.linelikes = lns
	// fmt.Printf("Layer after %s\n", l.Statistics())
	return l
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
		contents = append(contents, line.ControlLineXML(color, width))
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
