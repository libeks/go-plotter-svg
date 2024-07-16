package truchet

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/samplers"
)

func NewGrid(b box.Box, nx int, edgeMapping edgePointMapping, tileset []tileSet, dataSource samplers.DataSource) *Grid {
	boxes := b.PartitionIntoSquares(nx)
	cells := make(map[cellCoord]*Cell, len(boxes))
	grid := &Grid{
		edgePointMapping: edgeMapping,
		tileset:          tileset,
	}
	// grid.edgePointMapping = edgeMapping
	if len(boxes) != nx*nx {
		panic(fmt.Errorf("not right, want %d, got %d", nx*nx, len(boxes)))
	}
	horPoints := edgeMapping.getHorizontal()
	vertPoints := edgeMapping.getVertical()
	getIntersects := getRandomIncreasingIntersects
	grid.rowEdges = make(map[cellCoord]Edge, nx+1)
	for i := range nx + 1 { // for each of horizontal edges
		for j := range nx { // for each cell
			intersects := getIntersects(horPoints, float64(j)/float64(nx+1), float64(i)/float64(nx+1))
			grid.rowEdges[cellCoord{j, i}] = Edge{intersects} // flipped order is intentional
		}
	}
	grid.columnEdges = make(map[cellCoord]Edge, nx+1)
	for i := range nx + 1 { // for each of vertical edges
		for j := range nx { // for each cell
			intersects := getIntersects(vertPoints, float64(i)/float64(nx+1), float64(j)/float64(nx+1))
			grid.columnEdges[cellCoord{i, j}] = Edge{intersects}
		}
	}
	for _, childBox := range boxes {
		cell := &Cell{
			Grid: grid,
			Box:  childBox.Box,
			x:    childBox.I,
			y:    childBox.J,
		}
		cell.PopulateCurves(dataSource)
		cells[cellCoord{childBox.I, childBox.J}] = cell
	}
	grid.nX = nx
	grid.nY = nx
	grid.cells = cells
	return grid
}

// getIntersections returns the intersection points for coordinates in unit square
func getRandomIntersects(pointDef []endPointPair, xCoord, yCoord float64) []edgeMap {
	var intersects []edgeMap
	if len(pointDef) == 1 {
		intersects = []edgeMap{
			{
				point: pointDef[0],
				val:   maths.RandInRange(0.2, 0.8),
			},
		}
	} else if len(pointDef) == 2 {
		intersects = []edgeMap{
			{
				point: pointDef[0],
				val:   maths.RandInRange(0.2, 0.4),
			},
			{
				point: pointDef[1],
				val:   maths.RandInRange(0.6, 0.8),
			},
		}
	}
	return intersects
}

// getIntersections returns the intersection points for coordinates in unit square
func getRandomIncreasingIntersects(pointDef []endPointPair, xCoord, yCoord float64) []edgeMap {
	var intersects = make([]edgeMap, len(pointDef))
	spacing := 1 / float64(len(pointDef)+1)
	variance := 0.3 / float64(len(pointDef))
	for i, pt := range pointDef {
		center := spacing * float64(i+1)
		intersects[i] = edgeMap{
			point: pt,
			val:   maths.RandInRange(center-variance*yCoord, center+variance*yCoord),
		}
	}
	return intersects
}

// getIntersections returns the intersection points for coordinates in unit square
func getStaticIntersects(pointDef []endPointPair, xCoord, yCoord float64) []edgeMap {
	var intersects []edgeMap
	if len(pointDef) == 1 {
		intersects = []edgeMap{
			{
				point: pointDef[0],
				val:   0.5,
			},
		}
	} else if len(pointDef) == 2 {
		intersects = []edgeMap{
			{
				point: pointDef[0],
				val:   0.333,
			},
			{
				point: pointDef[1],
				val:   0.666,
			},
		}
	}
	return intersects
}

type edgeMap struct {
	point endPointPair
	val   float64
}

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

type Grid struct {
	nX    int
	nY    int
	cells map[cellCoord]*Cell
	// edge containers, specifying the position of cell border points
	columnEdges map[cellCoord]Edge
	rowEdges    map[cellCoord]Edge
	edgePointMapping
	tileset []tileSet
	samplers.DataSource
}

func (g Grid) At(x, y int) *Cell {
	if x < 0 || x >= g.nY || y < 0 || y >= g.nX {
		return nil
	}
	return g.cells[cellCoord{x, y}]
}

func (g Grid) GenerateCurve(cell *Cell, direction endPointTuple) lines.LineLike {
	startPoint := cell.AtEdge(direction, cell.GetEdgePoint(direction.endpoint))
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
				direction = g.edgePointMapping.other(nextDirection.endpoint)
			} else {
				return path
			}
		} else {
			return path
		}
	}
}

func (g Grid) GetGridLines() []lines.LineLike {
	ls := []lines.LineLike{}
	for _, cell := range g.cells {
		ls = append(ls, cell.Box.Lines()...)
	}
	return ls
}

func (g Grid) GererateCurves() []lines.LineLike {
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
					if !c.IsEmpty() {
						curves = append(curves, c)
					}
				}
			}
		}
	}
	return curves
}
