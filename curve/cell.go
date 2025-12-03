package curve

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/primitives"
)

// Cell represents a Truchet cell in a bigger grid, along with the curve fragments that it corresponds to
type Cell struct {
	*Grid
	primitives.BBox
	x int
	y int
	// tile
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

// GetCellInDirection returns the cell in the specified direction, if it exists, otherwise nil
func (c *Cell) GetCellInDirection(direction connectionEnd) *Cell {
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
func (c *Cell) VisitFrom(direction connectionEnd) (*Curve, *Cell, *connectionEnd) {
	for _, curve := range c.curves {
		if nextDir := curve.GetOtherEnd(direction); nextDir != nil {
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

// AtEdge returns a point on the edge of the cell specified at 'direction', interpolated at 't' on the edge
func (c *Cell) AtEdge(direction NWSE, t float64) primitives.Point {
	switch direction {
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
