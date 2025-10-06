package truchet

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/samplers"
)

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

func (c *Cell) IsDone() bool {
	for _, curve := range c.curves {
		if !curve.visited {
			return false
		}
	}
	return true
}

func (c *Cell) NextUnseen() *Curve {
	for _, curve := range c.curves {
		if !curve.visited {
			curve.visited = true
			return curve
		}
	}
	return nil
}

func (c *Cell) generateCurves(tileset tile) {
	// TODO: rename tileset to tile
	c.tile = tileset
	edgePointMap := c.GetEdgePoints()
	curves := make([]*Curve, len(tileset.pairs))
	for i, pair := range tileset.pairs {
		a := pair.a
		aDir := c.Grid.TruchetTileSet.EdgePointMapping.getDirection(a)
		b := pair.b
		bDir := c.Grid.TruchetTileSet.EdgePointMapping.getDirection(b)
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

func (c *Cell) GetEdgePoints() map[int]float64 {
	edges := map[NWSE]Edge{}
	edges[North] = c.Grid.rowEdges[cellCoord{c.x, c.y}]
	edges[South] = c.Grid.rowEdges[cellCoord{c.x, c.y + 1}]
	edges[West] = c.Grid.columnEdges[cellCoord{c.x, c.y}]
	edges[East] = c.Grid.columnEdges[cellCoord{c.x + 1, c.y}]
	vals := map[int]float64{}
	for _, edgePointMapping := range c.Grid.TruchetTileSet.EdgePointMapping.pairs {
		for _, endPointTuple := range []endPointTuple{edgePointMapping.a, edgePointMapping.b} {
			vals[endPointTuple.endpoint] = edges[endPointTuple.NWSE].GetPoint(endPointTuple.endpoint)
		}
	}
	return vals
}

func (c *Cell) GetEdgePoint(i int) float64 {
	// TODO: optimize this here code to not have to calculate the whole map
	return c.GetEdgePoints()[i]
}

func (c *Cell) GetCellInDirection(direction endPointTuple) *Cell {
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

func (c *Cell) PopulateCurves(dataSource samplers.DataSource) {
	rand := dataSource.GetValue(box.RelativeCenter(c.BBox)) // evaluate dataSource in absolute image coordinates
	l := len(c.Grid.TruchetTileSet.Tiles)
	n := int(rand * float64(l))
	if n == l {
		n = n - 1
	}
	tile := c.Grid.TruchetTileSet.Tiles[n]
	c.generateCurves(tile)
}

func (c *Cell) AtEdge(direction endPointTuple, t float64) primitives.Point {
	switch direction.NWSE {
	case North:
		return primitives.Point{X: maths.Interpolate(c.BBox.UpperLeft.X, c.BBox.LowerRight.X, t), Y: c.BBox.UpperLeft.Y}
	case West:
		return primitives.Point{X: c.BBox.UpperLeft.X, Y: maths.Interpolate(c.BBox.UpperLeft.Y, c.BBox.LowerRight.Y, t)}
	case South:
		return primitives.Point{X: maths.Interpolate(c.BBox.UpperLeft.X, c.BBox.LowerRight.X, t), Y: c.BBox.LowerRight.Y}
	case East:
		return primitives.Point{X: c.BBox.LowerRight.X, Y: maths.Interpolate(c.BBox.UpperLeft.Y, c.BBox.LowerRight.Y, t)}
	default:
		panic(fmt.Errorf("got composite direction %d", direction))
	}
}

func (c *Cell) At(horizontal, vertical float64) primitives.Point {
	return primitives.Point{
		X: maths.Interpolate(c.BBox.UpperLeft.X, c.BBox.LowerRight.X, horizontal),
		Y: maths.Interpolate(c.BBox.UpperLeft.Y, c.BBox.LowerRight.Y, vertical)}
}
