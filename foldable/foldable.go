package foldable

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

const (
	flapWidth = 150
	deg30     = math.Pi / 6
	deg45     = math.Pi / 4
)

type FoldablePattern struct {
	Edges       []lines.LineLike
	Fill        []lines.LineLike
	Annotations []lines.LineLike
}

type Edge struct {
	// todo: add support for curves
	primitives.Vector
}

type ConnectionType int

const (
	NoneConnection      ConnectionType = 0
	FlapConnection      ConnectionType = 1
	FlapSmallConnection ConnectionType = 2
	FaceConnection      ConnectionType = 3
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
	Name     string
	Connects map[int]Connection
}

func drawFlap(start primitives.Point, vector primitives.Vector, widthAngle float64) lines.LineLike {
	l := lines.NewPath(start)
	s := start
	// first leg is at 45° left of vector
	end := s.Add(vector.Unit().RotateCCW(-widthAngle).Mult(flapWidth / math.Sin(widthAngle)))
	l = l.AddPathChunk(lines.LineChunk{Start: s, End: end})
	s, end = end, end.Add(vector.Unit().Mult(vector.Len()-(2*flapWidth/math.Tan(widthAngle))))
	l = l.AddPathChunk(lines.LineChunk{Start: s, End: end})
	l = l.AddPathChunk(lines.LineChunk{Start: end, End: start.Add(vector)})
	return l
}

// Render draws this face, as well as any flaps and directly connected Faces
func (f Face) Render(start primitives.Point, angle float64) []lines.LineLike {
	fmt.Printf("angle %f, %v\n", angle, start)
	lns := []lines.LineLike{}
	l := lines.NewPath(start)
	for i, edge := range f.Shape.Edges {
		drawEdge := true
		if c, ok := f.Connects[i]; ok {
			switch c.Type {
			case FlapConnection:
				fmt.Printf("Attaching flap at angle %f\n", angle)
				lns = append(lns, drawFlap(start, edge.Vector.RotateCCW(angle), deg45))
			case FlapSmallConnection:
				fmt.Printf("Attaching flap at angle %f\n", angle)
				lns = append(lns, drawFlap(start, edge.Vector.RotateCCW(angle), deg30))
			case FaceConnection:
				nextFace := c.Face
				if nextFace == nil {
					panic("other face nil")
				}
				childAngle, diff := nextFace.GetEdgeAngle(c.OtherEdge)
				fmt.Printf("childAngle %f, diff %v\n", childAngle, diff)
				edgeAngle := edge.Atan() + angle
				newAngle := -childAngle + edgeAngle + math.Pi
				fmt.Printf("start %v\n", start)
				newStartpoint := start.Add(diff.RotateCCW(newAngle).Mult(-1))
				fmt.Printf("Attaching face at angle %f and at %v\n", newAngle, newStartpoint)
				lns = append(lns, nextFace.Render(newStartpoint, newAngle)...)
				drawEdge = false // don't draw this edge, the render of nextFace will draw it
			}
		}
		end := start.Add(edge.Vector.RotateCCW(angle))
		if drawEdge {
			l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
		} else {
			lns = append(lns, l)
			l = lines.NewPath(end) // start a new path
		}
		start = end

	}
	lns = append(lns, l)
	return lns
}

// WithFlap adds a standard flap on edge i
func (f Face) WithFlap(i int) Face {
	if i >= len(f.Edges) {
		panic("Not enough edges")
	}
	f.Connects[i] = Connection{Type: FlapConnection}
	return f
}

// WithSmallFlap adds a smaller flap if there is not enough space
func (f Face) WithSmallFlap(i int) Face {
	if i >= len(f.Edges) {
		panic("Not enough edges")
	}
	f.Connects[i] = Connection{Type: FlapSmallConnection}
	return f
}

// WithFace adds another face, with this face's edge i connecting to f2's face f2_i
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
