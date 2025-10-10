package foldable

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

type FaceID struct {
	Shape
	Name string
}

type ConnectionID struct {
	FaceA   string
	FaceB   string
	EdgeAID int
	EdgeBID int
	ConnectionType
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
		}
		faceB, ok := faceByID[connection.FaceB]
		if !ok {
			fmt.Printf("Could not find face named %s for connection number %d\n", connection.FaceB, i)
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
	faceBundle := initialFace.Render(b.UpperLeft, 0.25)
	fmt.Printf("FaceBundle map %v\n", faceBundle.FaceConfigs)
	return FoldablePattern{
		Edges:       faceBundle.Lines,
		Polygons:    []objects.Polygon{},
		Fill:        map[string]lines.LineLike{},
		Annotations: []lines.LineLike{},
	}
}
