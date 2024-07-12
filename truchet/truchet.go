package truchet

import (
	"fmt"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/samplers"
)

func TruchetTilesWithStrikeThrough(b box.Box, dataSource samplers.DataSource) []*Curve {
	val := dataSource.GetValue(b.Center())
	if val < 0.4 {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: West,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: East,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
		}
	} else if val > 0.6 {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: East,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: West,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
		}
	} else {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: LineOver,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: West,
						midpoint: 0.5,
					},
					{
						endpoint: East,
						midpoint: 0.5,
					},
				},
				CurveType: LineUnder,
			},
		}
	}
}

func TruchetTiles(b box.Box, dataSource samplers.DataSource) []*Curve {
	val := dataSource.GetValue(b.Center())
	if val < 0.5 {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: West,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: East,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
		}
	} else {
		return []*Curve{
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: North,
						midpoint: 0.5,
					},
					{
						endpoint: East,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
			{
				endpoints: []EndpointMidpoint{
					{
						endpoint: West,
						midpoint: 0.5,
					},
					{
						endpoint: South,
						midpoint: 0.5,
					},
				},
				CurveType: CircleSegment,
			},
		}
	}
}

type Winding int

const (
	Clockwise Winding = iota
	CounterClockwise
	Straight
	StraightUnder
	Undefined
)

func (w Winding) String() string {
	return []string{"Clockwise", "CounterClockwise", "Straight", "StraightUnder", "Undefined"}[w]
}

type NWSE int

const (
	North NWSE = iota
	West
	South
	East
)

func (d NWSE) Opposite() NWSE {
	switch d {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	case West:
		return East
	default:
		panic(fmt.Errorf("direction %s doesn't have an opposite", d))
	}
}

func (d NWSE) Winding(next NWSE) Winding {
	switch d {
	case North:
		switch next {
		case West:
			return CounterClockwise
		case South:
			return Straight
		case East:
			return Clockwise
		}
	case East:
		switch next {
		case North:
			return CounterClockwise
		case West:
			return Straight
		case South:
			return Clockwise
		}
	case South:
		switch next {
		case East:
			return CounterClockwise
		case North:
			return Straight
		case West:
			return Clockwise
		}
	case West:
		switch next {
		case South:
			return CounterClockwise
		case East:
			return Straight
		case North:
			return Clockwise
		}
	}
	return Undefined
}

func (d NWSE) String() string {
	return []string{"North", "West", "South", "East"}[d]
	// i := North
	// str := ""
	// for _, val := range []string{"North", "West", "South", "East"} {
	// 	if (d & i) > 0 {
	// 		str += val
	// 	}
	// 	i = i << 1
	// }
	// return str
}

type CurveType int

const (
	StraightLine CurveType = iota
	CircleSegment
	LineOver
	LineUnder
)

func (c CurveType) String() string {
	return []string{"StraightLine", "CircleSegment", "LineOver", "LineUnder"}[c]
}

type EndpointMidpoint struct {
	endpoint NWSE
	midpoint float64
}

func (e EndpointMidpoint) String() string {
	return fmt.Sprintf("%s %.1f", e.endpoint, e.midpoint)
}

type Curve struct {
	*Cell
	// endpoints NWSE
	// midpoints [2]float64 // in the range of [0;1], in increasing order of N,W,S,E
	endpoints []EndpointMidpoint
	CurveType
	visited bool
}

func (c *Curve) String() string {
	return fmt.Sprintf("Curve at %s with endpoints %v", c.Cell, c.endpoints)
}

func (c *Curve) XMLChunk(from NWSE) lines.PathChunk {
	if !c.HasEndpoint(from) {
		panic(fmt.Errorf("curve %s doesn't have endpoint %s", c, from))
	}
	to := c.GetOtherDirection(from)
	if to == nil {
		panic("No 'to' direction")
	}
	mTo := c.GetMidpoint(*to)
	mFrom := c.GetMidpoint(from)
	startPoint := c.Cell.At(from, *mFrom)
	endPoint := c.Cell.At(*to, *mTo)
	radius := c.Cell.Box.Width() / 2
	winding := from.Winding(*to)
	switch winding {
	case Straight:
		if c.CurveType == LineOver {
			fmt.Printf("doing line over %s\n", c)
			return lines.LineChunk{
				End: endPoint,
			}
		} else if c.CurveType == LineUnder {
			fmt.Printf("doing line under %s\n", c)
			return lines.LineGapChunk{
				Start:        startPoint,
				GapSizeRatio: 0.5,
				End:          endPoint,
			}
		} else {
			fmt.Printf("curve type %s\n", c.CurveType)
		}
	case StraightUnder:
		fmt.Printf("doing line under %s\n", c)
		return lines.LineGapChunk{
			Start:        startPoint,
			GapSizeRatio: 0.5,
			End:          endPoint,
		}
	case Clockwise:
		return lines.CircleArcChunk{
			Radius:      radius,
			IsClockwise: false, // Truchet circle arcs swing the other direction from winding
			IsLong:      false,
			End:         endPoint,
		}
	case CounterClockwise:
		return lines.CircleArcChunk{
			Radius:      radius,
			IsClockwise: true, // Truchet circle arcs swing the other direction from winding
			IsLong:      false,
			End:         endPoint,
		}
	case Undefined:
		fmt.Printf("winding is undefined: %s\n", winding)
		return lines.LineChunk{
			End: endPoint,
		}
	default:
		fmt.Printf("winding is %s\n", winding)
		return lines.LineChunk{
			End: endPoint,
		}
	}
	fmt.Printf("not even default: winding is %s\n", winding)

	return lines.LineChunk{
		End: endPoint,
	}
}

func (c *Curve) GetMidpoint(endpoint NWSE) *float64 {
	for _, pt := range c.endpoints {
		if pt.endpoint == endpoint {
			return &pt.midpoint
		}
	}
	return nil
}

func (c Curve) HasEndpoint(endpoint NWSE) bool {
	for _, pt := range c.endpoints {
		if pt.endpoint == endpoint {
			return true
		}
	}
	return false
}

func (c *Curve) GetOtherDirection(endpoint NWSE) *NWSE {
	var other *NWSE
	found := false
	for _, pt := range c.endpoints {
		if pt.endpoint == endpoint {
			found = true
		} else {
			other = &pt.endpoint
		}
	}
	if found {
		return other
	}
	return nil
}

// interpolate between a,b, with t in range [0,1]/ t=0 => a, t=1 => b
func interpolate(a, b, t float64) float64 {
	return (b-a)*t + a
}

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
		return primitives.Point{X: interpolate(c.Box.X, c.Box.XEnd, t), Y: c.Box.Y}
	case West:
		return primitives.Point{X: c.Box.X, Y: interpolate(c.Box.Y, c.Box.YEnd, t)}
	case South:
		return primitives.Point{X: interpolate(c.Box.X, c.Box.XEnd, t), Y: c.Box.YEnd}
	case East:
		return primitives.Point{X: c.Box.XEnd, Y: interpolate(c.Box.Y, c.Box.YEnd, t)}
	default:
		panic(fmt.Errorf("got composite direction %d", direction))
	}
}

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

type cellCoord struct {
	x int
	y int
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
