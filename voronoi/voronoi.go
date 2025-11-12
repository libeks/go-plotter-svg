package voronoi

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"

	"github.com/derekmu/voronoi"
)

type Edge struct {
	PolyIndex int
	EdgeIndex int
}

type EdgeMap struct {
	From *Edge
	To   *Edge
}

type VoronoiSet struct {
	Polygons []objects.Polygon
	EdgeMap  []EdgeMap
}

func ComputeVoronoiConnections(b primitives.BBox, points []primitives.Point) VoronoiSet {
	sites := make([]voronoi.Vertex, len(points))
	for i, point := range points {
		sites[i] = voronoi.Vertex{X: point.X, Y: point.Y}
	}
	fmt.Printf("Sites %v\n", sites)
	bbox := voronoi.BBox{Xl: b.UpperLeft.X, Xr: b.LowerRight.X, Yt: b.UpperLeft.Y, Yb: b.LowerRight.Y}
	fmt.Printf("BBox %v\n", bbox)
	diagram := voronoi.ComputeDiagram(sites, bbox, true)
	polys := make([]objects.Polygon, len(points))
	polygons := make(map[*voronoi.Cell]PointEdges, len(diagram.Cells))
	polygonIndices := make(map[*voronoi.Cell]int)
	edgeMap := []EdgeMap{}
	for i, cell := range diagram.Cells {
		// fmt.Printf("cell %d %v with %d edges\n", i, cell, len(cell.Halfedges))
		pointEdges := getEdgePoints(cell.Halfedges)
		// for j, edge := range sortEdges(cell.Halfedges) {
		// 	fmt.Printf("\tEdge %d (%v) from %v to %v\n", j, edge.reverse, edge.Edge.Va.Vertex, edge.Edge.Vb.Vertex)
		// }
		// fmt.Printf("pointedges %v\n", pointEdges.EdgeMap)
		polygons[cell] = pointEdges
		polygonIndices[cell] = i
		polys[i] = objects.Polygon{Points: pointEdges.Points}
	}
	visitedCells := make(map[*voronoi.Edge]struct{})
	for i, cell := range diagram.Cells {
		// fmt.Printf("Face %d\n", i)
		for _, edge := range sortEdges(cell.Halfedges) {
			// fmt.Printf("\tEdge %d %v\n", j, edge.Edge)
			// fmt.Printf("\tEdge %d (%v) from %v to %v\n", j, edge.reverse, edge.Edge.Va.Vertex, edge.Edge.Vb.Vertex)
			if _, ok := visitedCells[edge.Edge]; ok {
				continue
			}
			if otherCell := edge.Edge.GetOtherCell(cell); otherCell != nil {
				edgeMap = append(edgeMap, EdgeMap{
					From: &Edge{
						PolyIndex: i,
						EdgeIndex: polygons[cell].EdgeMap[edge.Edge],
					},
					To: &Edge{
						PolyIndex: polygonIndices[otherCell],
						EdgeIndex: polygons[otherCell].EdgeMap[edge.Edge],
					},
				})
				// fmt.Printf("Just added edge map from %d %d to %d %d\n", i, polygons[cell].EdgeMap[edge.Edge], polygonIndices[otherCell], polygons[otherCell].EdgeMap[edge.Edge])
			}
			visitedCells[edge.Edge] = struct{}{}
		}
	}
	// fmt.Printf("Polygons %v\n", polygons)
	return VoronoiSet{
		Polygons: polys,
		EdgeMap:  edgeMap,
	}
}

func ComputeVoronoi(b primitives.BBox, points []primitives.Point) []objects.Polygon {
	return ComputeVoronoiConnections(b, points).Polygons
}

type PointEdges struct {
	Points  []primitives.Point
	EdgeMap map[*voronoi.Edge]int
}

type Halfedge struct {
	*voronoi.Halfedge
	reverse bool
}

// https://stackoverflow.com/a/59299881
// go doesn't do the expected thing for modding, since (-1%5) = -1, but we want to get 4 (to wrap around the index)
func mod(a, b int) int {
	return (a%b + b) % b
}

func sortEdges(edges []*voronoi.Halfedge) []Halfedge {
	// reverse edges first
	// slices.Reverse(edges)
	reverseEdges := make([]*voronoi.Halfedge, len(edges))
	for i, edge := range edges {
		reverseEdges[len(edges)-1-i] = edge
	}
	halfEdges := make([]Halfedge, len(reverseEdges))
	for i, edge := range reverseEdges {
		nextEdge := reverseEdges[mod(i+1, len(reverseEdges))] // ensure that it wraps around beautifully
		var reverse bool
		// if the first vertex points to the following edge, flip the order
		if edge.Edge.Va.Vertex == nextEdge.Edge.Va.Vertex || edge.Edge.Va.Vertex == nextEdge.Edge.Vb.Vertex {
			reverse = true
		}
		halfEdges[i] = Halfedge{
			Halfedge: edge,
			reverse:  reverse,
		}
	}
	return halfEdges
}

// getEdgePoints transforms the list of non-consecutive edges into a list of points
func getEdgePoints(edges []*voronoi.Halfedge) PointEdges {
	points := []primitives.Point{}
	edgeMap := make(map[*voronoi.Edge]int)
	hedges := sortEdges(edges)
	for i, edge := range hedges {
		// fmt.Printf("GG \tEdge %d (%v) from %v to %v\n", i, edge.reverse, edge.Edge.Va.Vertex, edge.Edge.Vb.Vertex)
		edgeMap[edge.Edge] = mod(i-1, len(hedges))
		var pt primitives.Point
		if edge.reverse {
			pt = primitives.Point{X: edge.Edge.Va.X, Y: edge.Edge.Va.Y}
		} else {
			pt = primitives.Point{X: edge.Edge.Vb.X, Y: edge.Edge.Vb.Y}
		}
		points = append(points, pt)
	}
	return PointEdges{
		Points:  points,
		EdgeMap: edgeMap,
	}
}
