package curve

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/samplers"
)

var (
	connectNorth = connectionEnd{
		endpoint: 1,
		NWSE:     North,
	}
	connectSouth = connectionEnd{
		endpoint: 3,
		NWSE:     South,
	}
	connectEast = connectionEnd{
		endpoint: 2,
		NWSE:     East,
	}
	connectWest = connectionEnd{
		endpoint: 4,
		NWSE:     West,
	}
)

type marchingSquaresGrid struct {
	Grid
	// edge containers, specifying the position of cell border points
	source     samplers.DataSource
	gridValues map[cellCoord]float64
	gridStates map[cellCoord]bool
	threshold  float64
}

func NewMarchingGrid(b primitives.BBox, nx int, source samplers.DataSource, threshold float64) marchingSquaresGrid {
	boxes := primitives.PartitionIntoSquares(b, nx)
	cells := make(map[cellCoord]*Cell, len(boxes))
	grid := marchingSquaresGrid{
		source:     source,
		gridValues: map[cellCoord]float64{},
		gridStates: map[cellCoord]bool{},
	}
	if len(boxes) != nx*nx {
		panic(fmt.Errorf("not right, want %d, got %d", nx*nx, len(boxes)))
	}
	grid.Grid = Grid{
		nX:               nx,
		nY:               nx,
		cells:            cells,
		edgePointMapping: &EndpointMapping4,
		curveMapper:      MapStraightLines,
	}

	side := boxes[0].Width()
	vectX := primitives.Vector{X: side, Y: 0}
	vectY := primitives.Vector{X: 0, Y: side}
	for x := range nx + 1 {
		xx := float64(x)
		for y := range nx + 1 {
			yy := float64(y)
			coord := cellCoord{x, y}
			val := source.GetValue(primitives.Origin.Add(vectX.Mult(xx)).Add(vectY.Mult(yy)))
			grid.gridValues[coord] = val
			grid.gridStates[coord] = val > threshold
		}
	}

	for _, childBox := range boxes {
		cell := &Cell{
			Grid: &grid.Grid,
			BBox: childBox.BBox,
			x:    childBox.I,
			y:    childBox.J,
		}
		grid.PopulateCellCurveFragments(cell)
		cells[cellCoord{childBox.I, childBox.J}] = cell
	}
	grid.nX = nx
	grid.nY = nx
	grid.cells = cells
	return grid
}

func findInterpolatedTValue(a, b, threshold float64) float64 {
	if a < b {
		width := b - a
		return (threshold - a) / width
	} else {
		return 1 - findInterpolatedTValue(b, a, threshold)
	}
}

func addWNConnection(cell *Cell, wT, nT float64) {
	// TODO: Check that t-values are valid
	cell.curves = append(
		cell.curves,
		&Curve{
			endpoints: []EndpointMidpoint{
				{
					endpoint: connectWest,
					tValue:   wT,
				},
				{
					endpoint: connectNorth,
					tValue:   nT,
				},
			},
			visited: false,
			Cell:    cell,
		})
}

func addWSConnection(cell *Cell, wT, sT float64) {
	// TODO: Check that t-values are valid
	cell.curves = append(
		cell.curves,
		&Curve{
			endpoints: []EndpointMidpoint{
				{
					endpoint: connectWest,
					tValue:   wT,
				},
				{
					endpoint: connectSouth,
					tValue:   sT,
				},
			},
			visited: false,
			Cell:    cell,
		})
}

func addENConnection(cell *Cell, eT, nT float64) {
	// TODO: Check that t-values are valid
	cell.curves = append(
		cell.curves,
		&Curve{
			endpoints: []EndpointMidpoint{
				{
					endpoint: connectEast,
					tValue:   eT,
				},
				{
					endpoint: connectNorth,
					tValue:   nT,
				},
			},
			visited: false,
			Cell:    cell,
		})
}

func addESConnection(cell *Cell, eT, sT float64) {
	// TODO: Check that t-values are valid
	cell.curves = append(
		cell.curves,
		&Curve{
			endpoints: []EndpointMidpoint{
				{
					endpoint: connectEast,
					tValue:   eT,
				},
				{
					endpoint: connectSouth,
					tValue:   sT,
				},
			},
			visited: false,
			Cell:    cell,
		})
}

func addWEConnection(cell *Cell, wT, eT float64) {
	// TODO: Check that t-values are valid
	cell.curves = append(
		cell.curves,
		&Curve{
			endpoints: []EndpointMidpoint{
				{
					endpoint: connectWest,
					tValue:   wT,
				},
				{
					endpoint: connectEast,
					tValue:   eT,
				},
			},
			visited: false,
			Cell:    cell,
		})
}

func addNSConnection(cell *Cell, nT, sT float64) {
	// TODO: Check that t-values are valid
	cell.curves = append(
		cell.curves,
		&Curve{
			endpoints: []EndpointMidpoint{
				{
					endpoint: connectNorth,
					tValue:   nT,
				},
				{
					endpoint: connectSouth,
					tValue:   sT,
				},
			},
			visited: false,
			Cell:    cell,
		})
}

func (g *marchingSquaresGrid) PopulateCellCurveFragments(cell *Cell) {
	// points are numbered with x,y coords relative to the cell itself
	x := cell.x
	y := cell.y
	c00 := cellCoord{x, y}
	c01 := cellCoord{x, y + 1}
	c10 := cellCoord{x + 1, y}
	c11 := cellCoord{x + 1, y + 1}
	p00 := g.gridStates[c00]
	p01 := g.gridStates[c01]
	p10 := g.gridStates[c10]
	p11 := g.gridStates[c11]
	v00 := g.gridValues[c00]
	v01 := g.gridValues[c01]
	v10 := g.gridValues[c10]
	v11 := g.gridValues[c11]
	cell.curves = make([]*Curve, 0)
	wT := findInterpolatedTValue(v00, v01, g.threshold)
	nT := findInterpolatedTValue(v00, v10, g.threshold)
	eT := findInterpolatedTValue(v10, v11, g.threshold)
	sT := findInterpolatedTValue(v01, v11, g.threshold)
	if p00 == p01 == p10 == p11 {
		// no curve fragments in this cell
		return
	}
	if (p00 && !p01 && !p10 && !p00) || !(p00 && !p01 && !p10 && !p00) {
		// only top left corner is lit
		addWNConnection(cell, wT, nT)
		return
	}
	if (!p00 && p01 && !p10 && !p11) || !(!p00 && p01 && !p10 && !p11) {
		// only bottom left corner is lit
		addWSConnection(cell, wT, sT)
		return
	}
	if (!p00 && !p01 && p10 && !p11) || !(!p00 && !p01 && p10 && !p11) {
		// only top right corner is lit
		addENConnection(cell, eT, nT)
		return
	}
	if (!p00 && !p01 && !p10 && p11) || !(!p00 && !p01 && !p10 && p11) {
		// only bottom right corner is lit
		addESConnection(cell, eT, sT)
		return
	}
	if (p00 == p01) && (p10 == p11) {
		// the line is vertical
		addWEConnection(cell, wT, eT)
		return
	}
	if (p00 == p10) && (p01 == p11) {
		// the line is horizontal
		addNSConnection(cell, nT, sT)
		return
	}
	// the only remainder is the x-pattern, the saddle point option. ideally we should be checking the centerpoint
	// but i'll skip that for now
	addWNConnection(cell, wT, nT)
	addESConnection(cell, eT, sT)
}
