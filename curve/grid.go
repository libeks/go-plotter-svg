package curve

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/lines"
)

// TODO: Add non-square grids
// This includes a triangle grid, and a hexagon grid
// each cell would be expressed with (x,y) coordinates, but the "direction" will have less or more than the NWSE model, so this needs to be generalized
// For reference, here is an overview of coordinate systems for:
// * Hexes: https://www.redblobgames.com/grids/hexagons/
// * Triangles: https://simblob.blogspot.com/2007/06/distances-on-triangular-grid.html
//   * Consider a modification of (a,b) = (x, y+z). This converts it to a 2D grid, and each (a,b) value exists uniquely on the grid
//   * Categorize triangles into UP and DOWN. Let (0,0) be an UP triangle. Then even a+b is UP, odd is down.
//   * For UP triangles, there are three adjoining triangles:
//     * up-right direction changes coordinate by 	+(0, -1)
//     * up-left direction changes coordinate by	+(-1, 0)
//     * down direction changes coordinate by 		+(0, +1)
//   * For DOWN triangles, there are three adjoining triangles:
//     * down-left direction chanages coordinate by		+(0, +1)
//     * down-right direction changes coordinate by		+(+1, 0)
//     * up direction changes coordinate by				+(0, -1)
//   * The question is whether it is trivial to move from (a,b) coordinates to (x, y, z), since it seems that information gets lost, no?
//     * No, since the (x,y,z) contains coordinates that do not exist on the grid, such as (1,1,1). But what is the formula
//       for (x,y,z) coordinates that don't exist?
// This generalization can be applied to both Truchet grids, as well as Marching Squares

// Ideas for generalization:
// * Return gridlines
// * Return the perimeter edges
// * Compute all grid points (to get the Marching Squares points)

type Grid struct {
	// TODO: generalize nX, nY for triangle, hex grids
	nX    int
	nY    int
	cells map[cellCoord]*Cell
	*edgePointMapping
	curveMapper CurveMapper
}

// GetEdgePoints returns the t-value of the indexed edge connection
func (g Grid) GetEdgePoint(c *Cell, i int) float64 {
	for _, curve := range c.curves {
		for _, endpoint := range curve.endpoints {
			if endpoint.endpoint.endpoint == i {
				return endpoint.tValue
			}
		}
	}
	panic("Couldn't find GetEdgePoint")
}

func (g Grid) At(x, y int) *Cell {
	// TODO: generalize for triangle and hex grids
	if x < 0 || x >= g.nX || y < 0 || y >= g.nY {
		return nil
	}
	return g.cells[cellCoord{x, y}]
}

func (g Grid) GenerateCurve(cell *Cell, direction connectionEnd) lines.LineLike {
	// first figure out whether this cell has a connection at this direction, and get its t-value, to get the initial point
	edgeTValue := -1.0
	if cell == nil {
		panic("cell is nil")
	}
	for _, curve := range cell.curves {
		for _, endpoint := range curve.endpoints {
			if endpoint.endpoint == direction {
				edgeTValue = endpoint.tValue
				break
			}
		}
	}
	if edgeTValue < 0.0 {
		// this cell doesn't have any connections at this endpoint, so there is no curve to return
		return nil
	}
	startPoint := cell.AtEdge(direction.NWSE, edgeTValue)
	path := lines.NewPath(startPoint)
	for {
		// continue until all path chunks for this curve are exhausted
		if !cell.IsDone() {
			curve, nextCell, nextDirection := cell.VisitFrom(direction) // *Curve, *Cell, *NWSE
			if curve != nil {

				xml := curve.XMLChunk(g.curveMapper, direction)
				if !cell.PointInside(xml.Startpoint()) {
					panic(fmt.Sprintf("Startpoint %v is not inside bounding box %v\n", xml.Startpoint(), cell.BBox))
				}
				if !cell.PointInside(xml.Endpoint()) {
					panic(fmt.Sprintf("Endpoint %v is not inside bounding box %v\n", xml.Endpoint(), cell.BBox))
				}
				if xml.Startpoint().Subtract(startPoint).Len() > 1.0 {
					panic(fmt.Sprintf("Distance is %.2f\n", xml.Endpoint().Subtract(startPoint).Len()))
				}

				path = path.AddPathChunk(xml)
				if nextCell == nil {
					return path
				}
				startPoint = xml.Endpoint()
				cell = nextCell
				direction = g.edgePointMapping.other(nextDirection.endpoint)
			} else {
				// the next curve doesn't exist, we have either hit an edge, or we've come full circle
				return path
			}
		} else {
			// we've come full circle
			return path
		}
	}
}

func (g Grid) GetGridLines() []lines.LineLike {
	ls := []lines.LineLike{}
	for _, cell := range g.cells {
		ls = append(ls, lines.LinesFromBBox(cell.BBox)...)
	}
	return ls
}

func (g Grid) GenerateCurves() []lines.LineLike {
	curves := []lines.LineLike{}
	// start with perimeter
	// first from the top
	for x := range g.nX {
		cell := g.At(x, 0)
		direction := North
		for _, dirIndex := range g.edgePointMapping.endpointsFrom(direction) {
			curves = append(curves, g.GenerateCurve(cell, dirIndex))
		}
	}
	for y := range g.nY {
		cell := g.At(g.nX-1, y)
		direction := East
		for _, dirIndex := range g.edgePointMapping.endpointsFrom(direction) {
			curves = append(curves, g.GenerateCurve(cell, dirIndex))
		}
	}
	for x := g.nX - 1; x >= 0; x-- {
		cell := g.At(x, g.nY-1)
		direction := South
		for _, dirIndex := range g.edgePointMapping.endpointsFrom(direction) {
			curves = append(curves, g.GenerateCurve(cell, dirIndex))
		}
	}
	for y := g.nY - 1; y >= 0; y-- {
		cell := g.At(0, y)
		direction := West
		for _, dirIndex := range g.edgePointMapping.endpointsFrom(direction) {
			curves = append(curves, g.GenerateCurve(cell, dirIndex))
		}
	}
	for x := 0; x < g.nX; x++ {
		for y := 0; y < g.nY; y++ {
			cell := g.At(x, y)
			for _, direction := range []NWSE{North, West, South, East} {
				for _, dirIndex := range g.edgePointMapping.endpointsFrom(direction) {
					c := g.GenerateCurve(cell, dirIndex)
					if c != nil && !c.IsEmpty() {
						curves = append(curves, c)
					}
				}
			}
		}
	}
	return curves
}

type edgeMap struct {
	point connectionPair
	val   float64
}

// Edge represents the edge of a cell,
type Edge struct {
	intersects []edgeMap
}

func (e Edge) GetPoint(i int) float64 {
	for _, intersect := range e.intersects {
		if intersect.point.Has(i) {
			return intersect.val
		}
	}
	return -1
}

// cellCoord represents the indexed coordinates of a specific cell in a grid
type cellCoord struct {
	x int
	y int
}
