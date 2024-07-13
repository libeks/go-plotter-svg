package truchet

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/samplers"
)

func NewGrid(b box.Box, nx int, dataSource samplers.DataSource, curveConverter func(box.Box, samplers.DataSource) []*Curve) *Grid {
	boxes := b.PartitionIntoSquares(nx)
	cells := make(map[cellCoord]*Cell, len(boxes))
	grid := &Grid{}
	if len(boxes) != nx*nx {
		panic(fmt.Errorf("not right, want %d, got %d", nx*nx, len(boxes)))
	}
	for _, childBox := range boxes {
		cell := &Cell{
			Grid: grid,
			Box:  childBox.Box,
			x:    childBox.I,
			y:    childBox.J,
		}
		cell.PopulateCurves(curveConverter, dataSource)
		cells[cellCoord{childBox.I, childBox.J}] = cell
	}
	grid.nX = nx
	grid.nY = nx
	grid.cells = cells
	return grid
}

type Grid struct {
	nX    int
	nY    int
	cells map[cellCoord]*Cell
	samplers.DataSource
}

func (g Grid) At(x, y int) *Cell {
	if x < 0 || x >= g.nY || y < 0 || y >= g.nX {
		return nil
	}
	return g.cells[cellCoord{x, y}]
}

func (g Grid) GenerateCurve(cell *Cell, direction NWSE) lines.LineLike {
	startPoint := cell.At(direction, 0.5)
	path := lines.NewPath(startPoint)
	for {
		if !cell.IsDone() {
			curve, nextCell, nextDirection := cell.VisitFrom(direction) // *Curve, *Cell, *NWSE
			if curve != nil {
				path = path.AddPathChunk(curve.XMLChunk(direction))
				if nextCell == nil {
					return path
				}
				cell = nextCell
				direction = nextDirection.Opposite()
			} else {
				return path
			}
		} else {
			return path
		}
	}
	// return nil
}

func (g Grid) GererateCurves() []lines.LineLike {
	curves := []lines.LineLike{}
	// start with perimeter
	// first from the top
	for x := range g.nX {
		cell := g.At(x, 0)
		direction := North
		curves = append(curves, g.GenerateCurve(cell, direction))
	}
	for y := range g.nY {
		cell := g.At(g.nX-1, y)
		direction := East
		curves = append(curves, g.GenerateCurve(cell, direction))
	}
	for x := g.nX - 1; x >= 0; x-- {
		cell := g.At(x, g.nY-1)
		direction := South
		curves = append(curves, g.GenerateCurve(cell, direction))
	}
	for y := g.nY - 1; y >= 0; y-- {
		cell := g.At(0, y)
		direction := West
		curves = append(curves, g.GenerateCurve(cell, direction))
	}
	for x := 1; x < g.nX-1; x++ {
		for y := 1; y < g.nY-1; y++ {
			cell := g.At(x, y)
			for _, direction := range []NWSE{North, West, South, East} {
				c := g.GenerateCurve(cell, direction)
				if !c.IsEmpty() {

					curves = append(curves, c)
				}
			}
		}
	}
	return curves
}
