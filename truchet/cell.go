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
	box.Box
	x      int
	y      int
	tile   tileSet
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

func (c *Cell) generateCurves(tileset tileSet) {
	// TODO: rename tileset to tile
	c.tile = tileset
	edgePointMap := c.GetEdgePoints()
	curves := make([]*Curve, len(tileset.pairs))
	for i, pair := range tileset.pairs {
		a := pair.a
		aDir := c.Grid.edgePointMapping.getDirection(a)
		b := pair.b
		bDir := c.Grid.edgePointMapping.getDirection(b)
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
	for _, edgePointMapping := range c.Grid.edgePointMapping.pairs {
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
	rand := dataSource.GetValue(c.Box.Center())
	l := len(c.Grid.tileset)
	n := int(rand * float64(l))
	if n == l {
		n = n - 1
	}
	tile := c.Grid.tileset[n]
	c.generateCurves(tile)
}

func (c *Cell) AtEdge(direction endPointTuple, t float64) primitives.Point {
	switch direction.NWSE {
	case North:
		return primitives.Point{X: maths.Interpolate(c.Box.X, c.Box.XEnd, t), Y: c.Box.Y}
	case West:
		return primitives.Point{X: c.Box.X, Y: maths.Interpolate(c.Box.Y, c.Box.YEnd, t)}
	case South:
		return primitives.Point{X: maths.Interpolate(c.Box.X, c.Box.XEnd, t), Y: c.Box.YEnd}
	case East:
		return primitives.Point{X: c.Box.XEnd, Y: maths.Interpolate(c.Box.Y, c.Box.YEnd, t)}
	default:
		panic(fmt.Errorf("got composite direction %d", direction))
	}
}

func (c *Cell) At(horizontal, vertical float64) primitives.Point {
	return primitives.Point{X: maths.Interpolate(c.Box.X, c.Box.XEnd, horizontal), Y: maths.Interpolate(c.Box.Y, c.Box.YEnd, vertical)}
}
