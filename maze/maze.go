package maze

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

type Direction string

func (d Direction) Opposite() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	return ""
}

const (
	Up    Direction = "up"
	Down  Direction = "down"
	Left  Direction = "left"
	Right Direction = "right"
)

var (
	AllDirections = []Direction{Up, Down, Left, Right}
)

type Cell struct {
	X           int
	Y           int
	visited     bool
	rendered    bool
	connections []Direction
	grid        *Grid
}

func (c *Cell) Direction(d Direction) *Cell {
	switch d {
	case Up:
		return c.grid.At(c.X, c.Y-1)
	case Down:
		return c.grid.At(c.X, c.Y+1)
	case Left:
		return c.grid.At(c.X-1, c.Y)
	case Right:
		return c.grid.At(c.X+1, c.Y)
	default:
		return nil
	}
}

func NewGrid(size int) *Grid {
	cells := make([]*Cell, size*size)
	grid := &Grid{
		size: size,
	}
	for x := range size {
		for y := range size {
			cells[x*size+y] = &Cell{
				X:           x,
				Y:           y,
				connections: []Direction{},
				grid:        grid,
			}
		}
	}
	grid.cells = cells
	return grid
}

type Grid struct {
	size  int
	cells []*Cell
}

func (g *Grid) At(x, y int) *Cell {
	if x < 0 || y < 0 {
		return nil
	}

	if x >= g.size || y >= g.size {
		return nil
	}
	return g.cells[x*g.size+y]
}
func (g *Grid) CellSideInBox(b box.Box) float64 {
	return float64(max(b.Width(), b.Height())) / float64(g.size)
}

func (g *Grid) CellCenterInBox(b box.Box, x, y int) primitives.Point {
	edgeSize := g.CellSideInBox(b)
	return b.NWCorner().Add(primitives.Vector{X: (float64(x) + 0.5) * edgeSize, Y: (float64(y) + 0.5) * edgeSize})
}

func NewMaze(size int) *Maze {
	grid := NewGrid(size)
	start := grid.At(0, 0)
	fmt.Printf("start %v\n", start)
	start.visited = true
	stack := []*Cell{start}
	// perform a depth-first search until end
	for len(stack) > 0 {
		// fmt.Printf("stack size %d\n", len(stack))
		current := stack[0]
		possibilities := []Direction{}
		for _, direction := range AllDirections {
			if slices.Contains(current.connections, direction) {
				continue
			}
			if neighbor := current.Direction(direction); neighbor != nil && !neighbor.visited {
				possibilities = append(possibilities, direction)
			}
		}
		// fmt.Printf("possibilities at %d %d: %v\n", current.X, current.Y, possibilities)
		if len(possibilities) == 0 {
			stack = stack[1:] // pop first item
		} else {
			dir := possibilities[rand.Intn(len(possibilities))]
			current.connections = append(current.connections, dir)
			next := current.Direction(dir)
			next.visited = true
			next.connections = append(next.connections, dir.Opposite())
			stack = append([]*Cell{next}, stack...)
		}
	}
	return &Maze{
		Grid: grid,
	}
}

type Maze struct {
	*Grid
}

type MazeRender struct {
	Path  []lines.LineLike
	Walls []lines.LineLike
}

func (m *Maze) Render(b box.Box) MazeRender {
	start := m.Grid.At(0, 0)
	stack := []*Cell{start}

	paths := []lines.LineLike{}
	walls := []lines.LineLike{}

	fmt.Printf("Calculating path...\n")
	for len(stack) > 0 {
		current := stack[0]
		centerPoint := m.Grid.CellCenterInBox(b, current.X, current.Y)

		for _, dir := range current.connections {
			next := current.Direction(dir)
			if !next.rendered {
				nextCenter := m.Grid.CellCenterInBox(b, next.X, next.Y)
				paths = append(paths, lines.NewPath(centerPoint).AddPathChunk(lines.LineChunk{Start: centerPoint, End: nextCenter}))
				stack = append(stack, next)
				next.rendered = true
			}

		}
		for _, dir := range AllDirections {
			// fmt.Printf("contains %v, %s\n", current.connections, dir)
			isWall := slices.Contains(current.connections, dir)
			walls = append(walls, renderCellWall(centerPoint, dir, isWall, m.Grid.CellSideInBox(b))...)
		}
		stack = stack[1:]
	}
	fmt.Printf("Done.\n")

	return MazeRender{
		Path:  paths,
		Walls: walls,
	}
}

func renderCellWall(center primitives.Point, d Direction, isConnection bool, side float64) []lines.LineLike {
	wallWidth := 0.1
	roadWidth := 0.5 - wallWidth
	lns := []lines.LineLike{}
	if isConnection {
		switch d {
		case Up:
			start := center.Add(primitives.Vector{X: -roadWidth, Y: -0.5}.Mult(side))
			end := center.Add(primitives.Vector{X: -roadWidth, Y: -roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
			start = center.Add(primitives.Vector{X: roadWidth, Y: -0.5}.Mult(side))
			end = center.Add(primitives.Vector{X: roadWidth, Y: -roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
		case Down:
			start := center.Add(primitives.Vector{X: -roadWidth, Y: 0.5}.Mult(side))
			end := center.Add(primitives.Vector{X: -roadWidth, Y: roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
			start = center.Add(primitives.Vector{X: roadWidth, Y: 0.5}.Mult(side))
			end = center.Add(primitives.Vector{X: roadWidth, Y: roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
		case Left:
			start := center.Add(primitives.Vector{X: -0.5, Y: -roadWidth}.Mult(side))
			end := center.Add(primitives.Vector{X: -roadWidth, Y: -roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
			start = center.Add(primitives.Vector{X: -0.5, Y: roadWidth}.Mult(side))
			end = center.Add(primitives.Vector{X: -roadWidth, Y: roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
		case Right:
			start := center.Add(primitives.Vector{X: 0.5, Y: -roadWidth}.Mult(side))
			end := center.Add(primitives.Vector{X: roadWidth, Y: -roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
			start = center.Add(primitives.Vector{X: 0.5, Y: roadWidth}.Mult(side))
			end = center.Add(primitives.Vector{X: roadWidth, Y: roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
		}
	} else {
		// no connection
		switch d {
		case Up:
			start := center.Add(primitives.Vector{X: -roadWidth, Y: -roadWidth}.Mult(side))
			end := center.Add(primitives.Vector{X: roadWidth, Y: -roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
		case Down:
			start := center.Add(primitives.Vector{X: -roadWidth, Y: roadWidth}.Mult(side))
			end := center.Add(primitives.Vector{X: roadWidth, Y: roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
		case Left:
			start := center.Add(primitives.Vector{X: -roadWidth, Y: -roadWidth}.Mult(side))
			end := center.Add(primitives.Vector{X: -roadWidth, Y: roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
		case Right:
			start := center.Add(primitives.Vector{X: roadWidth, Y: -roadWidth}.Mult(side))
			end := center.Add(primitives.Vector{X: roadWidth, Y: roadWidth}.Mult(side))
			lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
		}
	}
	return lns
}
