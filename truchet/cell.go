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

func (c *Cell) GetCellInDirection(direction NWSE) *Cell {
	switch direction {
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
func (c *Cell) VisitFrom(direction NWSE) (*Curve, *Cell, *NWSE) {
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

func (c *Cell) PopulateCurves(curveConverter func(box box.Box, dataSource samplers.DataSource) []*Curve, dataSource samplers.DataSource) {
	c.curves = curveConverter(c.Box, dataSource)
	for _, curve := range c.curves {
		curve.Cell = c
		curve.visited = false
		// curve.CurveType = StraightLine
	}
}

func (c *Cell) At(direction NWSE, t float64) primitives.Point {
	switch direction {
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
