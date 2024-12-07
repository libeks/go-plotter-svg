package foldable

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

const (
	flapWidth = 150
)

type Edge struct {
	// todo: add support for curves
	primitives.Vector
}

type ConnectionType int

const (
	NoneConnection ConnectionType = 0
	FlapConnection ConnectionType = 1
	FaceConnection ConnectionType = 2
)

type Connection struct {
	Type      ConnectionType
	Face      *Face
	OtherEdge int // index of the edge of the other Face to connect to
}

func NewFace(shape Shape) Face {
	return Face{
		Shape:    shape,
		Connects: make(map[int]Connection, len(shape.Edges)),
	}
}

type Face struct {
	Shape
	Connects map[int]Connection
}

func drawFlap(start primitives.Point, vector primitives.Vector) lines.LineLike {
	l := lines.NewPath(start)
	s := start
	// first leg is at 45° left of vector
	end := s.Add(vector.Unit().RotateCCW(-math.Pi / 4).Mult(math.Sqrt(2 * flapWidth * flapWidth)))
	l = l.AddPathChunk(lines.LineChunk{Start: s, End: end})
	s, end = end, end.Add(vector.Unit().Mult(vector.Len()-2*flapWidth))
	l = l.AddPathChunk(lines.LineChunk{Start: s, End: end})
	l = l.AddPathChunk(lines.LineChunk{Start: end, End: start.Add(vector)})
	return l
}

// Render draws this face, as well as any flaps and directly connected Faces
func (f Face) Render(start primitives.Point, angle float64) []lines.LineLike {
	fmt.Printf("angle %f\n", angle)
	lns := []lines.LineLike{}
	l := lines.NewPath(start)
	for i, edge := range f.Shape.Edges {
		drawEdge := true
		if c, ok := f.Connects[i]; ok {
			if c.Type == FlapConnection {
				lns = append(lns, drawFlap(start, edge.Vector))
			} else if c.Type == FaceConnection {
				nextFace := c.Face
				if nextFace == nil {
					panic("other face nil")
				}
				atan, diff := nextFace.GetEdgeAngle(c.OtherEdge)
				lns = append(lns, nextFace.Render(start.Add(diff.Mult(-1)), atan)...)
				drawEdge = false // don't draw this edge, the render of nextFace will draw it
			}
		}
		end := start.Add(edge.Vector.RotateCCW(0))
		if drawEdge {
			l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
		} else {
			lns = append(lns, l)
			l = lines.NewPath(end) // start a new path
		}
		start = end

	}
	// fmt.Printf("lns %v\n", lns)
	lns = append(lns, l)
	return lns
}

func (f Face) WithFlap(i int) Face {
	if i >= len(f.Edges) {
		panic("Not enough edges")
	}
	f.Connects[i] = Connection{Type: FlapConnection}
	return f
}

func (f Face) WithFace(i int, f2 Face, f2_i int) Face {
	if i >= len(f.Edges) {
		panic("Not enough edges")
	}
	f.Connects[i] = Connection{
		Type:      FaceConnection,
		Face:      &f2,
		OtherEdge: f2_i,
	}
	return f
}
