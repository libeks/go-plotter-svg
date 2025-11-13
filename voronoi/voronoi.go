package voronoi

import (
	"github.com/libeks/go-plotter-svg/maths"
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
	// fmt.Printf("Sites %v\n", sites)
	pointIndex := map[primitives.Point]int{}
	for i, point := range points {
		pointIndex[point] = i
	}
	bbox := voronoi.BBox{Xl: b.UpperLeft.X, Xr: b.LowerRight.X, Yt: b.UpperLeft.Y, Yb: b.LowerRight.Y}
	diagram := voronoi.ComputeDiagram(sites, bbox, true)
	polys := make([]objects.Polygon, len(points))
	polygons := make(map[*voronoi.Cell]PointEdges, len(diagram.Cells))
	polygonIndices := make(map[*voronoi.Cell]int)
	edgeMap := []EdgeMap{}
	// reorder cells to match the input point order
	cells := make([]*voronoi.Cell, len(diagram.Cells))
	for _, cell := range diagram.Cells {
		i := pointIndex[primitives.Point{X: cell.Site.X, Y: cell.Site.Y}]
		cells[i] = cell
	}
	for i, cell := range cells {
		pointEdges := getEdgePoints(cell.Halfedges)
		polygons[cell] = pointEdges
		polygonIndices[cell] = i
		polys[i] = objects.Polygon{Points: pointEdges.Points}
	}
	visitedCells := make(map[*voronoi.Edge]struct{})
	for i, cell := range cells {
		for _, edge := range sortEdges(cell.Halfedges) {
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
			}
			visitedCells[edge.Edge] = struct{}{}
		}
	}
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

func sortEdges(edges []*voronoi.Halfedge) []Halfedge {
	// reverse edges first
	// slices.Reverse(edges)
	reverseEdges := make([]*voronoi.Halfedge, len(edges))
	for i, edge := range edges {
		reverseEdges[len(edges)-1-i] = edge
	}
	halfEdges := make([]Halfedge, len(reverseEdges))
	for i, edge := range reverseEdges {
		nextEdge := reverseEdges[maths.Mod(i+1, len(reverseEdges))] // ensure that it wraps around beautifully
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
		edgeMap[edge.Edge] = maths.Mod(i-1, len(hedges))
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
