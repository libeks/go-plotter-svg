package foldable

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/fonts"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

const MATH_PRECISION = 0.001

type FaceID struct {
	Shape
	infill
	Name string
}

type infill struct {
	color   string  // color for the infill, should be a csv defined color
	spacing float64 // spacing between lines
	angle   float64 // angle of the lines
	gap     float64 // distance from the edge of the polygon
}

func (f FaceID) WithFill(color string, spacing, angle, gap float64) FaceID {
	f.infill = infill{
		color:   color,
		spacing: spacing,
		angle:   angle,
		gap:     gap,
	}
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

func (c CutOut) Render(b primitives.BBox) FoldablePattern {
	// first, convert from a list representation to a face-tree representation
	faceByID := map[string]*Face{}
	var initialFace *Face
	connectionsCompleted := map[FaceEdge]bool{} // tracks whether all edge connections are accounted for
	for _, face := range c.Faces {
		if _, ok := faceByID[face.Name]; ok {
			fmt.Printf("Edge by the name %s already exists\n", face.Name)
		}
		faceByID[face.Name] = &Face{
			Shape:    face.Shape,
			Name:     face.Name,
			infill:   face.infill,
			Connects: map[int]Connection{},
		}
		for i := range face.Shape.Edges {
			connectionsCompleted[FaceEdge{face.Name, i}] = false
		}
	}
	for i, connection := range c.Connections {
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
		faceA.Connects[connection.EdgeAID] = Connection{
			Face:      faceB,
			Type:      connection.ConnectionType,
			OtherEdge: connection.EdgeBID,
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
	// set the initial face
	initialFace = faceByID[c.Connections[0].FaceA]
	faceBundle := initialFace.Render(b.UpperLeft, 0)
	polygons := []objects.Polygon{}
	annotations := []lines.LineLike{}
	fills := map[string][]lines.LineLike{}
	minAnnotationSize := math.MaxFloat64
	for key, polygon := range faceBundle.FacePolygons {
		polygons = append(polygons, polygon)
		bbox := polygon.LargestContainedSquareBBox()
		bbox = bbox.WithPadding(100)
		annotation := fonts.RenderText(bbox, key, fonts.WithSize(2000), fonts.WithFitToBox())
		if annotation.Size < minAnnotationSize {
			minAnnotationSize = annotation.Size
		}
		// annotations = append(annotations, fonts.RenderText(bbox, key, fonts.WithSize(2000), fonts.WithFitToBox()).CharCurves...)

		face := faceByID[key]
		if face.infill.color != "" {
			if _, ok := fills[face.infill.color]; !ok {
				fills[face.infill.color] = []lines.LineLike{}
			}
			infillPoly := polygon.Grow(-face.infill.gap)
			fills[face.infill.color] = append(fills[face.infill.color], infillPoly.LineFill(face.infill.angle, face.infill.spacing)...)
		}
	}
	// redo it again with the min annotation size
	for key, polygon := range faceBundle.FacePolygons {
		bbox := polygon.LargestContainedSquareBBox()
		bbox = bbox.WithPadding(100)
		annotations = append(annotations, fonts.RenderText(bbox, key, fonts.WithSize(minAnnotationSize)).CharCurves...)
	}
	polygons = append(polygons, faceBundle.FlapPolygons...)
	return FoldablePattern{
		Edges:       faceBundle.Lines,
		Polygons:    polygons,
		Fill:        fills,
		Annotations: annotations,
	}
}
