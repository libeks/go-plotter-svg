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
	cells := make(map[cellCoord]*Cell, len(boxes.BoxIterator()))
	grid := marchingSquaresGrid{
		source:     source,
		threshold:  threshold,
		gridValues: map[cellCoord]float64{},
		gridStates: map[cellCoord]bool{},
	}
	// if len(boxes) != nx*nx {
	// 	panic(fmt.Errorf("not right, want %d, got %d", nx*nx, len(boxes)))
	// }
	grid.Grid = Grid{
		nX:               boxes.NX,
		nY:               boxes.NY,
		cells:            cells,
		edgePointMapping: &EndpointMapping4,
		curveMapper:      MapStraightLines,
	}

	side := boxes.BoxWidth
	vectX := primitives.Vector{X: side, Y: 0}
	vectY := primitives.Vector{X: 0, Y: side}
	for x := range boxes.NX + 1 {
		xx := float64(x)
		for y := range boxes.NY + 1 {
			yy := float64(y)
			coord := cellCoord{x, y}
			val := source.GetValue(primitives.Origin.Add(vectX.Mult(xx)).Add(vectY.Mult(yy)))
			grid.gridValues[coord] = val
			grid.gridStates[coord] = val > threshold
		}
	}

	for _, childBox := range boxes.BoxIterator() {
		cell := &Cell{
			Grid: &grid.Grid,
			BBox: childBox.BBox,
			x:    childBox.I,
			y:    childBox.J,
		}
		grid.PopulateCellCurveFragments(cell)
		cells[cellCoord{childBox.I, childBox.J}] = cell
	}
	grid.nX = boxes.NX
	grid.nY = boxes.NY
	grid.cells = cells
	return grid
}

func findInterpolatedTValue(a, b, threshold float64) float64 {
	if a == b {
		// both endpoints are the same, default to 0.5 for consistency
		return 0.5
	}
	if a < b {
		width := b - a
		return (threshold - a) / width
	} else {
		width := a - b
		return (a - threshold) / width
	}
}

func checkTValue(a float64) {
	if (a < 0.0) || (a > 1.0) {
		panic(fmt.Sprintf("incorrect t value %.2f", a))
	}
}

func addWNConnection(cell *Cell, wT, nT float64) {
	checkTValue(wT)
	checkTValue(nT)
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
	checkTValue(wT)
	checkTValue(sT)
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
	checkTValue(eT)
	checkTValue(nT)
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
	checkTValue(eT)
	checkTValue(sT)
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
	checkTValue(wT)
	checkTValue(eT)
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
		return
	}
	if (p00 && !p01 && !p10 && !p11) || (!p00 && p01 && p10 && p11) {
		// only top left corner is lit
		addWNConnection(cell, wT, nT)
		return
	}
	if (!p00 && p01 && !p10 && !p11) || (p00 && !p01 && p10 && p11) {
		// only bottom left corner is lit
		addWSConnection(cell, wT, sT)
		return
	}
	if (!p00 && !p01 && p10 && !p11) || (p00 && p01 && !p10 && p11) {
		// only top right corner is lit
		addENConnection(cell, eT, nT)
		return
	}
	if (!p00 && !p01 && !p10 && p11) || (p00 && p01 && p10 && !p11) {
		// only bottom right corner is lit
		addESConnection(cell, eT, sT)
		return
	}
	if (p00 == p01) && (p10 == p11) {
		// the line is vertical
		addNSConnection(cell, nT, sT)
		return
	}
	if (p00 == p10) && (p01 == p11) {
		// the line is horizontal
		addWEConnection(cell, wT, eT)
		return
	}
	// the only remainder is the x-pattern, the saddle point option. Check the centerpoint.
	side := cell.Width()
	diagonal := primitives.Vector{X: side * .5, Y: side * .5}
	val := g.source.GetValue(cell.UpperLeft.Add(diagonal))
	pCenter := val > g.threshold
	// fmt.Printf("p00 %v, pCenter %v\n", p00, pCenter)
	if pCenter == p00 {
		// fmt.Printf("NE and SW\n")
		// NE and SW
		addWEConnection(cell, wT, eT)
		addWSConnection(cell, wT, sT)
	} else {
		// NW and SE
		addWNConnection(cell, wT, nT)
		addESConnection(cell, eT, sT)
	}

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
