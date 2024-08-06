package lines

import (
	"fmt"
	"strings"

	"github.com/shabbyrobe/xmlwriter"

	"github.com/libeks/go-plotter-svg/primitives"
)

const pathThreshold = 0.1

// TODO: don't take start as parameter for NewPath, infer it from chunks
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
	for _, chunk := range p.chunks {
		total += chunk.Length()
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
	start := p.start
	strs := []string{fmt.Sprintf("M %.1f %.1f", start.X, start.Y)}
	for _, xml := range p.chunks {
		if xml.Startpoint().Subtract(start).Len() > pathThreshold {
			panic(fmt.Errorf("path chunks are too far apart (%.1f): %s ", xml.Startpoint().Subtract(start).Len(), xml))
		}
		strs = append(strs, xml.PathXML())
		start = xml.Endpoint()
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
			{ // FIXME: remove
				Name:  "stroke-opacity",
				Value: "0.5",
			},
		},
	}
}

func (p Path) getControlLineString() string {
	strs := []string{fmt.Sprintf("M %.1f %.1f", p.start.X, p.start.Y)}
	for _, xml := range p.chunks {
		strs = append(strs, xml.ControlLines())
	}
	return strings.Join(strs, " ")
}

func (p Path) ControlLineXML(color, width string) xmlwriter.Elem {
	return xmlwriter.Elem{
		Name: "path", Attrs: []xmlwriter.Attr{
			{
				Name:  "d",
				Value: p.getControlLineString(),
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

func (p Path) OffsetLeft(distance float64) LineLike {
	fmt.Printf("offsetting left %s\n", p)
	chunks := make([]PathChunk, len(p.chunks))
	for i, chunk := range p.chunks {
		chunks[i] = chunk.OffsetLeft(distance)
	}
	fmt.Printf("new chunks %v\n", chunks)
	return Path{
		start:  chunks[0].Startpoint(),
		chunks: chunks,
	}
}

func (p Path) Reverse() LineLike {
	if len(p.chunks) == 0 {
		return p // noop, nothing to reverse
	}
	chunks := make([]PathChunk, len(p.chunks))
	for i := range len(p.chunks) {
		chunks[i] = p.chunks[len(p.chunks)-i-1].Reverse()
	}
	return Path{
		start:  chunks[0].Startpoint(),
		chunks: chunks,
	}
}

func (p Path) Bisect(t float64) (Path, Path) {
	if t < 0 || t > 1 {
		panic("invalid t parameter")
	}
	if len(p.chunks) == 0 {
		return p, p // noop
	}
	l := p.Len()
	targetLength := l * t
	var leftChunks, rightChunks []PathChunk
	var cumLen float64
	for i, chunk := range p.chunks {
		chunkLen := chunk.Length()
		if cumLen+chunkLen > targetLength {
			tStart := cumLen / l
			tEnd := (cumLen + chunkLen) / l
			localT := (t - tStart) / (tEnd - tStart)
			left, right := chunk.Bisect(localT) // FIXME: figure out new t
			leftChunks = append(leftChunks, left)
			rightChunks = append(rightChunks, right)
			rightChunks = append(rightChunks, p.chunks[i+1:]...)
			break
		} else {
			leftChunks = append(leftChunks, chunk)
			cumLen += chunk.Length()
		}
	}
	return Path{start: p.start, chunks: leftChunks}, Path{start: rightChunks[0].Startpoint(), chunks: rightChunks}
}

func (p Path) Join(q Path) Path {
	if p.End().Subtract(q.Start()).Len() > 0.1 {
		panic("the two paths don't join at the ends")
	}
	return Path{
		start:  p.start,
		chunks: append(p.chunks, q.chunks...),
	}
}
