package voronoi

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"

	"github.com/derekmu/voronoi"
)

func ComputeVoronoi(b primitives.BBox, points []primitives.Point) []objects.Polygon {
	sites := make([]voronoi.Vertex, len(points))
	for i, point := range points {
		sites[i] = voronoi.Vertex{X: point.X, Y: point.Y}
	}
	fmt.Printf("Sites %v\n", sites)
	bbox := voronoi.BBox{Xl: b.UpperLeft.X, Xr: b.LowerRight.X, Yt: b.UpperLeft.Y, Yb: b.LowerRight.Y}
	fmt.Printf("BBox %v\n", bbox)
	diagram := voronoi.ComputeDiagram(sites, bbox, true)
	polygons := make([]objects.Polygon, len(diagram.Cells))
	for i, cell := range diagram.Cells {
		fmt.Printf("cell %v with %d edges\n", cell, len(cell.Halfedges))
		polygons[i] = objects.Polygon{Points: getEdgePoints(cell.Halfedges)}
	}
	fmt.Printf("Polygons %v\n", polygons)
	return polygons
}

// getEdgePoints transforms the list of non-consecutive edges into a list of points
func getEdgePoints(edges []*voronoi.Halfedge) []primitives.Point {
	edgeMap := make(map[primitives.Point][]primitives.Point)
	visited := make(map[primitives.Point]struct{})
	points := []primitives.Point{}
	var startPt primitives.Point
	for _, edge := range edges {
		fmt.Printf("edge\n\tfrom %v\n\tto %v\n\tangle %f\n", edge.Edge.Va, edge.Edge.Vb, edge.Angle)
		ptA := primitives.Point{X: edge.Edge.Va.X, Y: edge.Edge.Va.Y}
		ptB := primitives.Point{X: edge.Edge.Vb.X, Y: edge.Edge.Vb.Y}
		edgeMap[ptA] = append(edgeMap[ptA], ptB)
		edgeMap[ptB] = append(edgeMap[ptB], ptA)
		startPt = ptA
	}
	points = append(points, startPt)
	for range edgeMap {
		fmt.Printf("Start point %v\n", startPt)
		others := edgeMap[startPt]
		for _, other := range others {
			fmt.Printf("Other point %v\n", other)
			_, hasVisited := visited[other]
			if (other.X != startPt.X || other.Y != startPt.Y) && !hasVisited {
				points = append(points, other)
				visited[startPt] = struct{}{}
				startPt = other
				fmt.Printf("Match found, ending early\n")
				break
			}
		}
	}
	fmt.Printf("Points %v\n", points)
	return points

}
