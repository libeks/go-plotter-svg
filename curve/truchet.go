package curve

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/samplers"
)

var (
	Truchet4NonCrossing = TruchetTileSet{
		Tiles:            TruchetPairs,
		EdgePointMapping: EndpointMapping4,
	}

	Truchet4Crossing = TruchetTileSet{
		Tiles:            TruchetUnderPairs,
		EdgePointMapping: EndpointMapping4,
	}

	Truchet6NonCrossingSide = TruchetTileSet{
		Tiles:            Truchet6Pairs,
		EdgePointMapping: EndpointMapping6Side,
	}
)

func NewTruchetGrid(b primitives.BBox, nx int, tileSet TruchetTileSet, tilePicker, edgeSource samplers.DataSource, curveMapper CurveMapper) *TruchetGrid {
	boxes := primitives.PartitionIntoSquares(b, nx)
	cells := make(map[cellCoord]*Cell, len(boxes))
	grid := &TruchetGrid{
		TruchetTileSet: tileSet,
	}
	if len(boxes) != nx*nx {
		panic(fmt.Errorf("not right, want %d, got %d", nx*nx, len(boxes)))
	}
	horPoints := tileSet.EdgePointMapping.getHorizontal()
	vertPoints := tileSet.EdgePointMapping.getVertical()

	grid.rowEdges = make(map[cellCoord]Edge, nx+1)
	for i := range nx + 1 { // for each of horizontal edges
		for j := range nx { // for each cell
			intersects := getSourcedIntersects(horPoints, edgeSource, float64(j)/float64(nx+1), float64(i)/float64(nx+1))
			grid.rowEdges[cellCoord{j, i}] = Edge{intersects} // flipped order is intentional
		}
	}
	grid.columnEdges = make(map[cellCoord]Edge, nx+1)
	for i := range nx + 1 { // for each of vertical edges
		for j := range nx { // for each cell
			intersects := getSourcedIntersects(vertPoints, edgeSource, float64(i)/float64(nx+1), float64(j)/float64(nx+1))
			grid.columnEdges[cellCoord{i, j}] = Edge{intersects}
		}
	}
	grid.Grid = Grid{
		nX:               nx,
		nY:               nx,
		cells:            cells,
		edgePointMapping: &tileSet.EdgePointMapping,
		curveMapper:      curveMapper,
	}
	for _, childBox := range boxes {
		cell := &Cell{
			Grid: &grid.Grid,
			BBox: childBox.BBox,
			x:    childBox.I,
			y:    childBox.J,
		}
		grid.PopulateCurves(cell, tilePicker)
		cells[cellCoord{childBox.I, childBox.J}] = cell
	}
	return grid
}

type TruchetGrid struct {
	Grid
	// edge containers, specifying the position of cell border points
	columnEdges map[cellCoord]Edge
	rowEdges    map[cellCoord]Edge
	TruchetTileSet
	// CurveMapper
	endpointWiggle samplers.DataSource
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

// getSourcedIntersections returns the intersection points for coordinates in unit square
func getSourcedIntersects(pointDef []connectionPair, edgeSource samplers.DataSource, xCoord, yCoord float64) []edgeMap {
	var intersects = make([]edgeMap, len(pointDef))
	spacing := 1 / float64(len(pointDef)+1)
	variance := 0.5 / float64(len(pointDef))
	for i, pt := range pointDef {
		center := spacing * float64(i+1)
		sourceVal := edgeSource.GetValue(primitives.Point{X: xCoord*2 - 1, Y: yCoord*2 - 1})
		valPlusMinus := sourceVal*2 - 1
		intersects[i] = edgeMap{
			point: pt,
			val:   center + valPlusMinus*variance/2,
		}
	}
	return intersects
}

func (g *TruchetGrid) generateCurves(c *Cell, tile tile) {
	edgePointMap := g.GetEdgePoints(c)
	curves := make([]*Curve, len(tile.pairs))
	for i, pair := range tile.pairs {
		a := pair.a
		aDir := g.TruchetTileSet.EdgePointMapping.getDirection(a)
		b := pair.b
		bDir := g.TruchetTileSet.EdgePointMapping.getDirection(b)
		curves[i] = &Curve{
			endpoints: []EndpointMidpoint{
				{
					endpoint: aDir,
					tValue:   edgePointMap[a],
				},
				{
					endpoint: bDir,
					tValue:   edgePointMap[b],
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

	vals := map[int]float64{}
	for _, edgePointMapping := range g.edgePointMapping.pairs {
		for _, endPointTuple := range edgePointMapping.bothEnds() {
			vals[endPointTuple.endpoint] = edges[endPointTuple.NWSE].GetPoint(endPointTuple.endpoint)
		}
	}
	return vals
}

func NewPair(a, b int) pair {
	return pair{a: a, b: b}
}

// A pair is a set of indices that define a curve fragment. 0-index means the connection has no other endpoint
type pair struct {
	a int
	b int
}

func (p pair) Other(q int) int {
	if p.a == q {
		return p.b
	} else if p.b == q {
		return p.a
	}
	return -1
}

// TruchetTileSet is a definition of what tiles are available for this Truchet configuration
type TruchetTileSet struct {
	Tiles            []tile
	EdgePointMapping edgePointMapping
}

// A tile consists of pairs of indexes for which edges should be connected
type tile struct {
	pairs []pair
}

func (t tile) Other(i int) int {
	for _, pair := range t.pairs {
		if other := pair.Other(i); other > 0 {
			return other
		}
	}
	return -1
}

// non-intersecting links for a 4-set, corresponds to Catalan number 2
var TruchetPairs = []tile{
	{
		pairs: []pair{
			NewPair(1, 2),
			NewPair(3, 4),
		},
	},
	{
		pairs: []pair{
			NewPair(1, 4),
			NewPair(2, 3),
		},
	},
}

// all links for 4-set, including a straight-through intersection in the middle
var TruchetUnderPairs = []tile{
	{
		pairs: []pair{
			NewPair(1, 2),
			NewPair(3, 4),
		},
	},
	{
		// cross
		pairs: []pair{
			NewPair(1, 3),
			NewPair(2, 4),
		},
	},
	{
		pairs: []pair{
			NewPair(1, 4),
			NewPair(2, 3),
		},
	},
}

// non-intersecting links for a 6-set, corresponds to Catalan number 3
var Truchet6Pairs = []tile{
	{
		// ()()()
		pairs: []pair{
			NewPair(1, 2),
			NewPair(3, 4),
			NewPair(5, 6),
		},
	},
	{
		// ()(())
		pairs: []pair{
			NewPair(1, 2),
			NewPair(3, 6),
			NewPair(4, 5),
		},
	},
	{
		// (())()
		pairs: []pair{
			NewPair(1, 4),
			NewPair(2, 3),
			NewPair(5, 6),
		},
	},
	{
		// ((()))
		pairs: []pair{
			NewPair(1, 6),
			NewPair(2, 5),
			NewPair(3, 4),
		},
	},
	{
		// (()())
		pairs: []pair{
			NewPair(1, 6),
			NewPair(2, 3),
			NewPair(4, 5),
		},
	},
}

var EndpointMapping4 = edgePointMapping{
	[]connectionPair{
		{
			a: connectionEnd{
				endpoint: 1,
				NWSE:     North,
			},
			b: connectionEnd{
				endpoint: 3,
				NWSE:     South,
			},
		},
		{
			a: connectionEnd{
				endpoint: 2,
				NWSE:     East,
			},
			b: connectionEnd{
				endpoint: 4,
				NWSE:     West,
			},
		},
	},
}

var EndpointMapping6Side = edgePointMapping{
	[]connectionPair{
		{
			a: connectionEnd{
				endpoint: 1,
				NWSE:     North,
			},
			b: connectionEnd{
				endpoint: 4,
				NWSE:     South,
			},
		},
		{
			a: connectionEnd{
				endpoint: 2,
				NWSE:     East,
			},
			b: connectionEnd{
				endpoint: 6,
				NWSE:     West,
			},
		},
		{
			a: connectionEnd{
				endpoint: 3,
				NWSE:     East,
			},
			b: connectionEnd{
				endpoint: 5,
				NWSE:     West,
			},
		},
	},
}

type EndpointMidpoint struct {
	endpoint connectionEnd
	tValue   float64
}

func (e EndpointMidpoint) String() string {
	return fmt.Sprintf("%s %.1f", e.endpoint, e.tValue)
}
