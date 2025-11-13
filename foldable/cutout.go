package foldable

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/collections"
	"github.com/libeks/go-plotter-svg/fonts"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/pen"
	"github.com/libeks/go-plotter-svg/primitives"
)

const MATH_PRECISION = 0.001

type FaceID struct {
	Shape
	infills []infill
	Name    string
}

type infill struct {
	color   string  // color for the infill, should be a csv defined color
	pen.Pen         // either pen is defined, or spacing and gap are defined, never both
	spacing float64 // spacing between lines
	angle   float64 // angle of the lines
	gap     float64 // distance from the edge of the polygon
}

func (f FaceID) WithFill(color string, spacing, angle, gap float64) FaceID {
	f.infills = append(f.infills, infill{
		color:   color,
		spacing: spacing,
		angle:   angle,
		gap:     gap,
	})
	return f
}

func (f FaceID) WithBrushFill(color string, p pen.Pen) FaceID {
	f.infills = append(f.infills, infill{
		color: color,
		Pen:   p,
	})
	return f
}

func faceID(s Shape, n string) FaceID {
	return FaceID{
		Shape: s,
		Name:  n,
	}
}

type ConnectionID struct {
	FaceA   string
	FaceB   string
	EdgeAID int
	EdgeBID int
	ConnectionType
}

func smallFlap(faceA, faceB string, edgeA, edgeB int) ConnectionID {
	return ConnectionID{
		FaceA:          faceA,
		FaceB:          faceB,
		EdgeAID:        edgeA,
		EdgeBID:        edgeB,
		ConnectionType: FlapSmallConnection,
	}
}

func flap(faceA, faceB string, edgeA, edgeB int) ConnectionID {
	return ConnectionID{
		FaceA:          faceA,
		FaceB:          faceB,
		EdgeAID:        edgeA,
		EdgeBID:        edgeB,
		ConnectionType: FlapConnection,
	}
}

func link(faceA, faceB string, edgeA, edgeB int) ConnectionID {
	return ConnectionID{
		FaceA:          faceA,
		FaceB:          faceB,
		EdgeAID:        edgeA,
		EdgeBID:        edgeB,
		ConnectionType: FaceConnection,
	}
}

type CutOut struct {
	//
	Faces []FaceID
	// Connections map[int]Connection
	Connections []ConnectionID
}

func NewCutOut(faces []FaceID, connections []ConnectionID) CutOut {
	return CutOut{
		Faces:       faces,
		Connections: connections,
	}
}

type FaceEdge struct {
	Face   string
	EdgeID int
}

// cutoutTrees represents all of the information necessary to render each pattern in the foldable
type cutoutTrees struct {
	faces map[string]*Face // the faces, keyed by their labels
	heads []string         // a list of the roots of the trees, listing their labels
}

func (c CutOut) computeTrees(container primitives.BBox) cutoutTrees {
	// first, convert from a list representation to a face-tree representation
	faceByID := map[string]*Face{}
	// heads are the faces that start a distinct foldable pattern.
	// Each is a face that only has outgoing direct connnections, no incoming
	visitedFaces := map[string]struct{}{}
	heads := []string{}
	// var initialFace *Face
	connectionsCompleted := map[FaceEdge]bool{} // tracks whether all edge connections are accounted for
	for _, face := range c.Faces {
		if _, ok := faceByID[face.Name]; ok {
			fmt.Printf("Edge by the name %s already exists\n", face.Name)
		}
		faceByID[face.Name] = &Face{
			Shape:    face.Shape,
			Name:     face.Name,
			infills:  face.infills,
			Connects: map[int]Connection{},
		}
		for i := range face.Shape.Edges {
			connectionsCompleted[FaceEdge{face.Name, i}] = false
		}
	}
	for i, connection := range c.Connections {
		if connection.ConnectionType == NoneConnection {
			// noop, there is no connection
			edge := FaceEdge{
				connection.FaceA,
				connection.EdgeAID,
			}
			connectionsCompleted[edge] = true
			continue
		}
		faceA, ok := faceByID[connection.FaceA]
		if !ok {
			fmt.Printf("Could not find face named %s for connection number %d\n", connection.FaceA, i)
			panic("Couldn't render")
		}
		faceB, ok := faceByID[connection.FaceB]
		if !ok {
			fmt.Printf("Could not find face named %s for connection number %d\n", connection.FaceB, i)
			panic("Couldn't render")
		}
		if connection.ConnectionType == FaceConnection {
			if _, ok := visitedFaces[faceB.Name]; ok {
				fmt.Printf("Face %s has already been connected to, this creates a cycle in the graph", faceB.Name)
				panic("Couldn't render")
			}
			visitedFaces[faceB.Name] = struct{}{}
		}

		if connection.EdgeAID >= len(faceA.Shape.Edges) {
			fmt.Printf("Face %s with %d faces doesn't have an edge number %d\n", connection.FaceA, len(faceA.Shape.Edges), connection.EdgeAID)
			panic("Couldn't render")
		}
		if connection.EdgeBID >= len(faceB.Shape.Edges) {
			fmt.Printf("Face %s with %d faces doesn't have an edge number %d\n", connection.FaceB, len(faceB.Shape.Edges), connection.EdgeBID)
			panic("Couldn't render")
		}
		// check that the two edges are of the same length
		aLen := faceA.Shape.Edges[connection.EdgeAID].Vector.Len()
		bLen := faceB.Shape.Edges[connection.EdgeBID].Vector.Len()
		if math.Abs(aLen-bLen) > MATH_PRECISION {
			fmt.Printf("The connected edges %s:%d and %s:%d have different lengths (%.3f vs %.3f = diff of %.3f)\n",
				connection.FaceA, connection.EdgeAID,
				connection.FaceB, connection.EdgeBID,
				aLen, bLen, math.Abs(aLen-bLen),
			)
		}

		if connection.ConnectionType == DoubleConnection {
			faceA.Connects[connection.EdgeAID] = Connection{
				Face:      faceB,
				Type:      FlapConnection, // make sure this is a flap, not DoubleConnection
				OtherEdge: connection.EdgeBID,
			}
			// add a full flap on the other side as well
			faceB.Connects[connection.EdgeBID] = Connection{
				Face:      faceA,
				Type:      FlapConnection,
				OtherEdge: connection.EdgeAID,
			}
		} else {
			faceA.Connects[connection.EdgeAID] = Connection{
				Face:      faceB,
				Type:      connection.ConnectionType,
				OtherEdge: connection.EdgeBID,
			}
		}
		edges := []FaceEdge{
			{
				connection.FaceA,
				connection.EdgeAID,
			},
			{
				connection.FaceB,
				connection.EdgeBID,
			},
		}
		for _, edge := range edges {
			if connectionsCompleted[edge] {
				fmt.Printf("Edge %s:%d is already connected elsewhere\n", edge.Face, edge.EdgeID)
			}
			connectionsCompleted[edge] = true
		}
	}
	for edge, ok := range connectionsCompleted {
		if !ok {
			fmt.Printf("Edge %s:%d is not connected to anything\n", edge.Face, edge.EdgeID)
		}
	}
	// for each face that has never appeared as faceB in a face connection, it is a root/head of its own tree
	for name := range faceByID {
		if _, ok := visitedFaces[name]; !ok {
			heads = append(heads, name)
		}
	}
	return cutoutTrees{
		faces: faceByID,
		heads: heads,
	}
}

type BrushLines struct {
	pen.Pen
	Color string
	Lines []lines.LineLike
}

// GeneratePatterns creates a list of foldable patterns, each of which represents a standalone component of the foldable
// which can be placed somewhere on the page. Each pattern will have its bounding box start at the origin
func (c CutOut) GeneratePatterns(container primitives.BBox) []FoldablePattern {
	trees := c.computeTrees(container)
	patterns := []FoldablePattern{}
	for _, headLabel := range trees.heads {
		head := trees.faces[headLabel]
		faceBundle := head.Render(primitives.Origin, 0)
		polygons := []objects.Polygon{}
		annotations := []lines.LineLike{}
		fills := map[string]BrushLines{}
		minAnnotationSize := math.MaxFloat64
		for key, polygon := range faceBundle.FacePolygons {
			polygons = append(polygons, polygon)
			bbox := polygon.LargestContainedSquareBBox()
			bbox = bbox.WithPadding(100)
			annotation := fonts.RenderText(bbox, key, fonts.WithSize(2000), fonts.WithFitToBox())
			// ensure that all annotations are rendered at the same size, which is the smallest size you can render
			if annotation.Size < minAnnotationSize && annotation.Size > 0 {
				minAnnotationSize = annotation.Size
			}
			face := trees.faces[key]
			for _, infill := range face.infills {
				infillLabel := infill.color
				if infill.Pen.Name != "" {
					infillLabel = fmt.Sprintf("%s %s", infill.color, infill.Pen.Name)
				}
				if _, ok := fills[infillLabel]; !ok {
					brush := BrushLines{
						Lines: []lines.LineLike{},
						Color: infill.color,
					}

					if infill.Pen.Name != "" {
						brush.Pen = infill.Pen
					}
					fills[infillLabel] = brush
				}
				if infill.Pen.Name != "" {
					brushLines := fills[infillLabel]
					brushLines.Lines = append(brushLines.Lines, collections.FillPolygonWithPen(polygon, infill.Pen)...)
					fills[infillLabel] = brushLines
				} else {
					infillPoly := polygon.Grow(-infill.gap)
					brushLines := fills[infillLabel]
					brushLines.Lines = append(brushLines.Lines, infillPoly.LineFill(infill.angle, infill.spacing)...)
					fills[infillLabel] = brushLines
				}
			}
		}
		// redo it again with the min annotation size
		for key, polygon := range faceBundle.FacePolygons {
			bbox := polygon.LargestContainedSquareBBox()
			bbox = bbox.WithPadding(100)
			annotations = append(annotations, fonts.RenderText(bbox, key, fonts.WithSize(minAnnotationSize)).CharCurves...)
		}
		polygons = append(polygons, faceBundle.FlapPolygons...)
		patterns = append(patterns, FoldablePattern{
			Edges:       faceBundle.Lines,
			Polygons:    polygons,
			Fill:        fills,
			Annotations: annotations,
		},
		)
	}
	return patterns
}

func (c CutOut) Render(container primitives.BBox) []FoldablePattern {
	patterns := c.GeneratePatterns(container)
	return patterns
}
