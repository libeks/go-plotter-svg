package foldable

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

const (
	flapWidth = 150
	deg30     = math.Pi / 6
	deg45     = math.Pi / 4
)

type FoldablePattern struct {
	Edges       []lines.LineLike            // edges draw such that there shouldn't be overlap
	Polygons    []objects.Polygon           // a list of polygons that are contained in the fold, including flap polygons
	Fill        map[string][]lines.LineLike // a collection of fills for the faces, mapped by color
	Annotations []lines.LineLike            // contains face labels, etc. that won't always need to be plotted
}

func (p FoldablePattern) BBox() primitives.BBox {
	if len(p.Polygons) == 0 {
		return primitives.BBox{}
	}
	box := p.Polygons[0].BBox()
	for _, poly := range p.Polygons[1:len(p.Polygons)] {
		box = box.Add(poly.BBox())
	}
	return box
}

func (p FoldablePattern) Translate(v primitives.Vector) FoldablePattern {
	edges := make([]lines.LineLike, len(p.Edges))
	for i, edge := range p.Edges {
		edges[i] = edge.Translate(v)
	}
	polygons := make([]objects.Polygon, len(p.Polygons))
	for i, poly := range p.Polygons {
		polygons[i] = poly.Translate(v)
	}
	fills := make(map[string][]lines.LineLike, len(p.Fill))
	for key, subfills := range p.Fill {
		for i, fill := range subfills {
			subfills[i] = fill.Translate(v)
		}
		fills[key] = subfills
	}
	annotations := make([]lines.LineLike, len(p.Annotations))
	for i, annotation := range p.Annotations {
		annotations[i] = annotation.Translate(v)
	}
	return FoldablePattern{
		Edges:       edges,
		Polygons:    polygons,
		Fill:        fills,
		Annotations: annotations,
	}
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

func drawFlap(start primitives.Point, vector primitives.Vector, widthAngle float64) lines.Path {
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

type FaceConfig struct {
	Start primitives.Point
	Angle float64
}

type RenderBundle struct {
	Lines        []lines.LineLike
	FaceConfigs  map[string]FaceConfig
	FacePolygons map[string]objects.Polygon
	FlapPolygons []objects.Polygon
}

// Render draws this face, as well as any flaps and directly connected Faces.
// This is used in a recursive tree-like process to draw each edge exactly once.
func (f Face) Render(start primitives.Point, angle float64) RenderBundle {
	// fmt.Printf("angle %f, %v\n", angle, start)
	lns := []lines.LineLike{}
	l := lines.NewPath(start)
	faceConfigs := map[string]FaceConfig{}
	faceConfigs[f.Name] = FaceConfig{
		Start: start,
		Angle: angle,
	}
	facePolygons := map[string]objects.Polygon{}
	flapPolygons := []objects.Polygon{}
	facePoints := []primitives.Point{start}
	for i, edge := range f.Shape.Edges {
		drawEdge := true
		if c, ok := f.Connects[i]; ok {
			switch c.Type {
			case FlapConnection:
				// fmt.Printf("Attaching flap at angle %f\n", angle)
				flap := drawFlap(start, edge.Vector.RotateCCW(angle), deg45)
				flapPolygons = append(flapPolygons, objects.Polygon{Points: flap.Points()})
				lns = append(lns, flap)
			case FlapSmallConnection:
				// fmt.Printf("Attaching flap at angle %f\n", angle)
				flap := drawFlap(start, edge.Vector.RotateCCW(angle), deg30)
				flapPolygons = append(flapPolygons, objects.Polygon{Points: flap.Points()})
				lns = append(lns, flap)
			case FaceConnection:
				nextFace := c.Face
				if nextFace == nil {
					panic("other face nil")
				}
				childAngle, diff := nextFace.GetEdgeAngle(c.OtherEdge)
				// fmt.Printf("childAngle %f, diff %v\n", childAngle, diff)
				edgeAngle := edge.Atan() + angle
				newAngle := -childAngle + edgeAngle + math.Pi
				// fmt.Printf("start %v\n", start)
				newStartpoint := start.Add(diff.RotateCCW(newAngle).Mult(-1))
				// fmt.Printf("Attaching face at angle %f and at %v\n", newAngle, newStartpoint)
				faceBundle := nextFace.Render(newStartpoint, newAngle)
				for key, faceConfig := range faceBundle.FaceConfigs {
					faceConfigs[key] = faceConfig
				}
				for key, facePolygon := range faceBundle.FacePolygons {
					facePolygons[key] = facePolygon
				}
				for _, flapPolygon := range faceBundle.FlapPolygons {
					flapPolygons = append(flapPolygons, flapPolygon)
				}
				lns = append(lns, faceBundle.Lines...)
				drawEdge = false // don't draw this edge, the render of nextFace will draw it
			}
		}
		end := start.Add(edge.Vector.RotateCCW(angle))
		facePoints = append(facePoints, end)
		if drawEdge {
			l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
		} else {
			lns = append(lns, l)
			l = lines.NewPath(end) // start a new path
		}
		start = end

	}
	lns = append(lns, l)
	fmt.Printf("Adding polygon %s\n", f.Name)
	facePolygons[f.Name] = objects.Polygon{Points: facePoints}
	return RenderBundle{
		Lines:        lns,
		FaceConfigs:  faceConfigs,
		FacePolygons: facePolygons,
		FlapPolygons: flapPolygons,
	}
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
