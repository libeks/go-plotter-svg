package truchet

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/samplers"
)

// Cell represents a Truchet cell in a bigger grid, along with the curve fragments that it corresponds to
type Cell struct {
	*Grid
	primitives.BBox
	x int
	y int
	tile
	curves []*Curve
}

func (c *Cell) String() string {
	return fmt.Sprintf("Cell (%d, %d)", c.x, c.y)
}

// IsDone returns true if all curve fragments are done/visited for this cell
func (c *Cell) IsDone() bool {
	for _, curve := range c.curves {
		if !curve.visited {
			return false
		}
	}
	return true
}

// NextUnseen returns the next curve fragment in this cell that still hasn't been visited, if any
func (c *Cell) NextUnseen() *Curve {
	for _, curve := range c.curves {
		if !curve.visited {
			curve.visited = true
			return curve
		}
	}
	return nil
}

func (g *TruchetGrid) generateCurves(c *Cell, tileset tile) {
	// TODO: rename tileset to tile
	c.tile = tileset
	edgePointMap := g.GetEdgePoints(c)
	curves := make([]*Curve, len(tileset.pairs))
	for i, pair := range tileset.pairs {
		a := pair.a
		aDir := g.TruchetTileSet.EdgePointMapping.getDirection(a)
		b := pair.b
		bDir := g.TruchetTileSet.EdgePointMapping.getDirection(b)
		curves[i] = &Curve{
			endpoints: []EndpointMidpoint{
				{
					endpoint: aDir,
					midpoint: edgePointMap[a],
				},
				{
					endpoint: bDir,
					midpoint: edgePointMap[b],
				},
			},
			visited: false,
			Cell:    c,
		}
	}
	c.curves = curves
}

// GetEdgePoints returns a map from the edge index to its corresponding t-values
func (g TruchetGrid) GetEdgePoints(c *Cell) map[int]float64 {
	// edges contains the t-values on each of the edges of this cell
	edges := map[NWSE]Edge{}
	edges[North] = g.rowEdges[cellCoord{c.x, c.y}]
	edges[South] = g.rowEdges[cellCoord{c.x, c.y + 1}]
	edges[West] = g.columnEdges[cellCoord{c.x, c.y}]
	edges[East] = g.columnEdges[cellCoord{c.x + 1, c.y}]
	fmt.Printf("edges %v\n", edges)

	vals := map[int]float64{}
	for _, edgePointMapping := range g.TruchetTileSet.EdgePointMapping.pairs {
		for _, endPointTuple := range []endPointTuple{edgePointMapping.a, edgePointMapping.b} {
			vals[endPointTuple.endpoint] = edges[endPointTuple.NWSE].GetPoint(endPointTuple.endpoint)
		}
	}
	return vals
}

// GetEdgePoints returns the t-value of the indexed edge connection
func (g TruchetGrid) GetEdgePoint(c *Cell, i int) float64 {
	// TODO: optimize this here code to not have to calculate the whole map
	return g.GetEdgePoints(c)[i]
}

// GetCellInDirection returns the cell in the specified direction, if it exists, otherwise nil
func (c *Cell) GetCellInDirection(direction endPointTuple) *Cell {
	if c.Grid == nil {
		panic("NIL")
	}
	switch direction.NWSE {
	case North:
		return c.Grid.At(c.x, c.y-1)
	case South:
		return c.Grid.At(c.x, c.y+1)
	case West:
		return c.Grid.At(c.x-1, c.y)
	case East:
		return c.Grid.At(c.x+1, c.y)
	default:
		panic(fmt.Errorf("unrecognized direction %s", direction))
	}
}

// return the curve for this cell that starts from direction and next cell (if any)
func (c *Cell) VisitFrom(direction endPointTuple) (*Curve, *Cell, *endPointTuple) {
	for _, curve := range c.curves {
		if nextDir := curve.GetOtherDirection(direction); nextDir != nil {
			if curve.visited {
				continue // curve is already visited, don't double-count
			}
			curve.visited = true
			return curve, c.GetCellInDirection(*nextDir), nextDir
		}
	}
	return nil, nil, nil
}

func getRelativeAroundCenter(v float64) float64 {
	relative := (v) / 10_000
	return 2*relative - 1
}

// center of the box in relative coordinates [0.0, 1.0], assuming that the image is in the range [0, 10_000]
func relativeCenter(b primitives.BBox) primitives.Point {
	return primitives.Point{
		X: getRelativeAroundCenter(b.UpperLeft.X + b.Width()/2.0),
		Y: getRelativeAroundCenter(b.UpperLeft.Y + b.Height()/2.0),
	}
}

// PopulateCurves decides which Truchet tile to use, and populates the curve fragments that fall inside of this cell
func (g *TruchetGrid) PopulateCurves(c *Cell, dataSource samplers.DataSource) {
	rand := dataSource.GetValue(relativeCenter(c.BBox)) // evaluate dataSource in absolute image coordinates
	l := len(g.TruchetTileSet.Tiles)
	n := int(rand * float64(l))
	// rand could produce a value of 1.0, which would map to be outside of the range. We cap this to the last element, since this is a weird edge case
	if n == l {
		n = n - 1
	}
	tile := g.TruchetTileSet.Tiles[n]
	g.generateCurves(c, tile)
}

// AtEdge returns a point on the edge of the cell specified at 'direction', interpolated at 't' on the edge
func (c *Cell) AtEdge(direction endPointTuple, t float64) primitives.Point {
	switch direction.NWSE {
	case North:
		return c.At(t, 0)
	case West:
		return c.At(0, t)
	case South:
		return c.At(t, 1)
	case East:
		return c.At(1, t)
	default:
		panic(fmt.Errorf("got composite direction %d", direction))
	}
}

// At returns the point inside cell, with relative coordinates in the range [0.0, 1.0]
func (c *Cell) At(horizontal, vertical float64) primitives.Point {
	return primitives.Point{
		X: maths.Interpolate(c.BBox.UpperLeft.X, c.BBox.LowerRight.X, horizontal),
		Y: maths.Interpolate(c.BBox.UpperLeft.Y, c.BBox.LowerRight.Y, vertical)}
}
