package foldable

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

type Shape struct {
	// Shape is an abstract polygon without a positioning in space
	// Vectors in clockwise order describing the shape
	// by convention, the first vector should be horizontal
	// by convention, edges should appear in clockwise order
	Edges []Edge
}

// check that the shape is closed
func (s Shape) Verify() bool {
	v := primitives.Vector{X: 0, Y: 0}
	for _, edge := range s.Edges {
		v = v.Add(edge.Vector)
	}
	if v.Len() > 0 {
		panic(fmt.Sprintf("shape is not closed, %v of overlap", v))
		// return false
	}
	return true
}

// GetEdgeAngle returns the angle that the i-th edge makes with the initial edge
func (s Shape) GetEdgeAngle(i int) (float64, primitives.Vector) {
	if i >= len(s.Edges) {
		panic("edge index too high")
	}
	v := primitives.Vector{X: 0, Y: 0}
	for j, edge := range s.Edges {
		v = v.Add(edge.Vector)
		if i == j {
			return edge.Atan(), v
		}
	}
	return 0, primitives.Vector{X: 0, Y: 0}

}

func Square(side float64) Shape {
	return Shape{
		Edges: []Edge{
			{
				Vector: primitives.Vector{X: 1, Y: 0}.Mult(side),
			},
			{
				Vector: primitives.Vector{X: 0, Y: 1}.Mult(side),
			},
			{
				Vector: primitives.Vector{X: -1, Y: 0}.Mult(side),
			},
			{
				Vector: primitives.Vector{X: 0, Y: -1}.Mult(side),
			},
		},
	}
}

func Rectangle(a, b float64) Shape {
	return Shape{
		Edges: []Edge{
			{
				Vector: primitives.Vector{X: a, Y: 0},
			},
			{
				Vector: primitives.Vector{X: 0, Y: b},
			},
			{
				Vector: primitives.Vector{X: -a, Y: 0},
			},
			{
				Vector: primitives.Vector{X: 0, Y: -b},
			},
		},
	}
}

func EquiTriangle(side float64) Shape {
	return Shape{
		Edges: []Edge{
			{
				Vector: primitives.Vector{X: 0.5, Y: -math.Sqrt(3) / 2}.Mult(side),
			},
			{
				Vector: primitives.Vector{X: 0.5, Y: math.Sqrt(3) / 2}.Mult(side),
			},
			{
				Vector: primitives.Vector{X: -1, Y: 0}.Mult(side),
			},
		},
	}
}

func PolygonToShape(p objects.Polygon) Shape {
	fmt.Printf("polygon %v\n", p)
	edges := []Edge{}
	for i, point := range p.Points {
		if i == 0 {
			continue
		}
		edges = append(edges, Edge{Vector: point.Subtract(p.Points[i-1])})
	}
	edges = append(edges, Edge{Vector: p.Points[0].Subtract(p.Points[len(p.Points)-1])})
	return Shape{Edges: edges}
}
