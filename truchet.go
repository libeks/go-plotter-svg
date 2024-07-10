package main

import (
	"fmt"
	"math/rand"
)

func truchetTiles(box Box, dataSource DataSource) []*Curve {
	val := dataSource.GetValue(box.Center())
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

type NWSE int

const (
	North NWSE = 0x1
	West  NWSE = 0x2
	South NWSE = 0x4
	East  NWSE = 0x8
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
		panic(fmt.Errorf("Direction %s doesn't have an opposite", d))
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
	i := North
	str := ""
	for _, val := range []string{"North", "West", "South", "East"} {
		if (d & i) > 0 {
			str += val
		}
		i = i << 1
	}
	return str
}

type CurveType int

const (
	StraightLine CurveType = iota
	CircleSegment
	OverCurve
	UnderCurve
)

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

func (c *Curve) XMLChunk(from NWSE) PathChunk {
	if !c.HasEndpoint(from) {
		panic(fmt.Errorf("curve %s doesn't have endpoint %s", c, from))
	}
	to := c.GetOtherDirection(from)
	if to == nil {
		// return ""
		panic("No to direction")
	}
	// mFrom := c.GetMidpoint(from)
	mTo := c.GetMidpoint(*to)
	endPoint := c.Cell.At(*to, *mTo)
	radius := c.Cell.Box.Width() / 2
	winding := from.Winding(*to)
	switch winding {
	case Straight:
		return LineChunk{
			endpoint: endPoint,
		}
		// return fmt.Sprintf("L %.1f %.1f", endPoint.x, endPoint.y)
	case Clockwise:
		return CircleArcChunk{
			radius:      radius,
			isClockwise: false, // Truchet circle arcs swing the other direction from winding
			isLong:      false,
			endpoint:    endPoint,
		}
		// return fmt.Sprintf("A %.1f %.1f 0 0 %d %.1f %.1f", radius, radius, 0, endPoint.x, endPoint.y)
	case CounterClockwise:
		return CircleArcChunk{
			radius:      radius,
			isClockwise: true, // Truchet circle arcs swing the other direction from winding
			isLong:      false,
			endpoint:    endPoint,
		}
		// return fmt.Sprintf("A %.1f %.1f 0 0 %d %.1f %.1f", radius, radius, 1, endPoint.x, endPoint.y)
	}
	// return fmt.Sprintf("A %.1f %.1f 0 0 %d %.1f %.1f", radius, radius, swing, endPoint.x, endPoint.y)
	return LineChunk{
		endpoint: endPoint,
	}
	// return fmt.Sprintf("L %.1f %.1f", endPoint.x, endPoint.y)
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
	// fmt.Printf("Getting other direction for %s, from %s\n", c, endpoint)
	for _, pt := range c.endpoints {
		if pt.endpoint == endpoint {
			found = true
		} else {
			other = &pt.endpoint
		}
	}
	if found {
		// fmt.Printf("Other direction for %s, from %s is %s\n", c, endpoint, other)
		return other
	}
	// fmt.Printf("Other direction for %s, from %s is nil\n", c, endpoint)
	return nil
}

// interpolate between a,b, with t in range [0,1]/ t=0 => a, t=1 => b
func interpolate(a, b, t float64) float64 {
	return (b-a)*t + a
}

type Cell struct {
	*Grid
	Box
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

func (c *Cell) PopulateCurves(curveConverter func(box Box, dataSource DataSource) []*Curve, dataSource DataSource) {
	c.curves = curveConverter(c.Box, dataSource)
	for _, curve := range c.curves {
		curve.Cell = c
		curve.visited = false
		curve.CurveType = StraightLine
	}
}

func (c *Cell) At(direction NWSE, t float64) Point {
	switch direction {
	case North:
		return Point{interpolate(c.Box.x, c.Box.xEnd, t), c.Box.y}
	case West:
		return Point{c.Box.x, interpolate(c.Box.y, c.Box.yEnd, t)}
	case South:
		return Point{interpolate(c.Box.x, c.Box.xEnd, t), c.Box.yEnd}
	case East:
		return Point{c.Box.xEnd, interpolate(c.Box.y, c.Box.yEnd, t)}
	default:
		panic(fmt.Errorf("got composite direction %d", direction))
	}
}

func NewGrid(box Box, nx int, dataSource DataSource, curveConverter func(box Box, dataSource DataSource) []*Curve) *Grid {
	boxes := partitionIntoSquares(box, nx)
	cells := make(map[cellCoord]*Cell, len(boxes))
	grid := &Grid{}
	if len(boxes) != nx*nx {
		panic(fmt.Errorf("not right, want %d, got %d", nx*nx, len(boxes)))
	}
	for _, childBox := range boxes {
		cell := &Cell{
			Grid: grid,
			Box:  childBox.box,
			x:    childBox.i,
			y:    childBox.j,
		}
		cell.PopulateCurves(curveConverter, dataSource)
		cells[cellCoord{childBox.i, childBox.j}] = cell
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
	DataSource
}

func (g Grid) At(x, y int) *Cell {
	// fmt.Printf("Getting at %d %d\n", x, y)
	if x < 0 || x >= g.nY || y < 0 || y >= g.nX {
		// fmt.Printf("Nothing at %d %d\n", x, y)
		return nil
	}
	// fmt.Printf("Returning  %s\n", g.cells[cellCoord{x, y}])
	return g.cells[cellCoord{x, y}]
}

func (g Grid) GenerateCurve(cell *Cell, direction NWSE) LineLike {
	startPoint := cell.At(direction, 0.5)
	// instructions := []string{fmt.Sprintf("M %.1f %.1f", startPoint.x, startPoint.y)}
	path := NewPath(startPoint)
	for {
		// fmt.Printf("GenerateCurve %s %s\n", cell, direction)
		if !cell.IsDone() {
			curve, nextCell, nextDirection := cell.VisitFrom(direction) // *Curve, *Cell, *NWSE
			if curve != nil {
				path = path.AddPathChunk(curve.XMLChunk(direction))
				// instructions = append(instructions, curve.XML(direction))
				if nextCell == nil {
					// if len(instructions) > 1 {
					// return Path{s: strings.Join(instructions, " ")}
					// }
					return path
				}
				cell = nextCell
				direction = nextDirection.Opposite()
				// fmt.Printf("GenerateCurve next is %s %s\n", cell, direction)
			} else {
				// if len(instructions) > 1 {
				// return Path{s: strings.Join(instructions, " ")}
				// }
				// return Path{s: ""}
				return path
			}
		} else {
			// if len(instructions) > 1 {
			// 	return Path{s: strings.Join(instructions, " ")}
			// }
			// return Path{s: ""}
			return path
		}
	}
	// return nil
}

func (g Grid) GererateCurves() []LineLike {
	curves := []LineLike{}
	// start with perimeter
	// first from the top
	// fmt.Printf("Top row\n")
	for x := range g.nX {
		cell := g.At(x, 0)
		direction := North
		curves = append(curves, g.GenerateCurve(cell, direction))
	}
	// fmt.Printf("Right column\n")
	for y := range g.nY {
		cell := g.At(g.nX-1, y)
		direction := East
		curves = append(curves, g.GenerateCurve(cell, direction))
	}
	// fmt.Printf("Bottom row\n")
	for x := g.nX - 1; x >= 0; x-- {
		cell := g.At(x, g.nY-1)
		direction := South
		curves = append(curves, g.GenerateCurve(cell, direction))
	}
	// fmt.Printf("Left column\n")
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

type DataSource interface {
	GetValue(p Point) float64
}

type ConstantDataSource struct {
	val float64
}

func (s ConstantDataSource) GetValue(p Point) float64 {
	return s.val
}

type RandomDataSource struct {
}

func (s RandomDataSource) GetValue(p Point) float64 {
	return rand.Float64()
}
