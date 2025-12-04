package curve

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/lines"
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

// func (g marchingSquaresGrid) GererateCurves() []lines.LineLike {
// 	return g.Grid.GererateCurves()
// }

func NewMarchingGrid(b primitives.BBox, nx int, source samplers.DataSource, threshold float64) marchingSquaresGrid {
	boxes := primitives.PartitionIntoSquares(b, nx)
	cells := make(map[cellCoord]*Cell, len(boxes))
	grid := marchingSquaresGrid{
		source:     source,
		threshold:  threshold,
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
	if a == b {
		// fmt.Printf("a %.2f, b%.2f, thresh %.2f => 0.5\n", a, b, threshold)
		return 0.5
	}
	if a < b {
		width := b - a
		// fmt.Printf("a %.2f, b%.2f, thresh %.2f => %.2f\n", a, b, threshold, (threshold-a)/width)
		return (threshold - a) / width
	} else {
		width := a - b
		// fmt.Printf("a %.2f, b%.2f, thresh %.2f => %.2f\n", a, b, threshold, (a-threshold)/width)
		return (a - threshold) / width
	}
}

func checkTValue(a float64) {
	if (a < 0.0) || (a > 1.0) {
		panic(fmt.Sprintf("incorrect t value %.2f", a))
	}
}

func addWNConnection(cell *Cell, wT, nT float64) {
	// TODO: Check that t-values are valid
	checkTValue(wT)
	checkTValue(nT)
	// fmt.Printf("Adding WN from %.1f to %.1f\n", wT, nT)
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
	checkTValue(wT)
	checkTValue(sT)
	// fmt.Printf("Adding WS from %.1f to %.1f\n", wT, sT)
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
	checkTValue(eT)
	checkTValue(nT)
	// fmt.Printf("Adding EN from %.1f to %.1f\n", eT, nT)
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
	checkTValue(eT)
	checkTValue(sT)
	// fmt.Printf("Adding ES from %.1f to %.1f\n", sT, sT)
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
	checkTValue(wT)
	checkTValue(eT)
	// fmt.Printf("Adding WE from %.1f to %.1f\n", wT, eT)
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
	checkTValue(nT)
	checkTValue(sT)
	// TODO: Check that t-values are valid
	// fmt.Printf("Adding NS from %.1f to %.1f\n", nT, sT)
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

	if (p00 == p01) && (p01 == p10) && (p10 == p11) {
		// no curve fragments in this cell
		// fmt.Printf("Noop\n")
		return
	}
	// fmt.Printf("Cell (%d, %d), has %v %v %v %v\n", cell.x, cell.y, p00, p01, p10, p11)
	// fmt.Printf("w-t (%.1f), n-t (%.1f), e-t (%.1f), s-t (%.1f)\n", wT, nT, eT, sT)
	if (p00 && !p01 && !p10 && !p11) || (!p00 && p01 && p10 && p11) {
		// only top left corner is lit
		// fmt.Printf("Top left\n")
		// fmt.Printf("v00 %.2f, v01 %.2f, %.2f\n", v00, v01, wT)
		// wt := findInterpolatedTValue(v00, v01, g.threshold)
		// fmt.Printf("Got %.2f\n", wt)
		addWNConnection(cell, wT, nT)
		// fmt.Printf("Added WN\n")
		return
	}
	if (!p00 && p01 && !p10 && !p11) || (p00 && !p01 && p10 && p11) {
		// only bottom left corner is lit
		// fmt.Printf("Bottom left\n")
		addWSConnection(cell, wT, sT)
		// fmt.Printf("Added WS\n")
		return
	}
	if (!p00 && !p01 && p10 && !p11) || (p00 && p01 && !p10 && p11) {
		// fmt.Printf("Top right\n")
		// only top right corner is lit
		addENConnection(cell, eT, nT)
		// fmt.Printf("Added EN\n")
		return
	}
	if (!p00 && !p01 && !p10 && p11) || (p00 && p01 && p10 && !p11) {
		// only bottom right corner is lit
		// fmt.Printf("Bottom right\n")
		addESConnection(cell, eT, sT)
		// fmt.Printf("Added ES\n")
		return
	}
	if (p00 == p01) && (p10 == p11) {
		// the line is vertical
		// fmt.Printf("Vertical\n")
		addNSConnection(cell, nT, sT)
		// fmt.Printf("Added WE\n")
		return
	}
	if (p00 == p10) && (p01 == p11) {
		// the line is horizontal
		// fmt.Printf("Horizontal\n")
		addWEConnection(cell, wT, eT)
		// fmt.Printf("Added NS\n")
		return
	}
	// the only remainder is the x-pattern, the saddle point option. ideally we should be checking the centerpoint
	// but i'll skip that for now
	fmt.Printf("Saddle\n")
	addWNConnection(cell, wT, nT)
	fmt.Printf("Added WN\n")
	addESConnection(cell, eT, sT)
	fmt.Printf("Added ES\n")
}

func (g *marchingSquaresGrid) GetControlPoints() []lines.LineLike {
	lnes := []lines.LineLike{}
	for coord, state := range g.gridStates {
		cell := g.cells[coord]
		// fmt.Printf("cell %v %v\n", coord, cell)
		if cell != nil {
			pt := g.cells[coord].UpperLeft
			if state {
				lnes = append(lnes, lines.Cross(pt, 30)...)
			} else {
				lnes = append(lnes, lines.FullCircle(pt, 30))
			}
		}
	}
	return lnes
}
