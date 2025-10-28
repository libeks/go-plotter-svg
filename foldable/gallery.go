package foldable

import (
	"math"

	"github.com/libeks/go-plotter-svg/pen"
	"github.com/libeks/go-plotter-svg/primitives"
)

// Cube does a cube with a certain side length.
func Cube(b primitives.BBox, side float64) []FoldablePattern {
	sq := Square(side)
	c := NewCutOut(
		[]FaceID{
			faceID(sq, "A"),
			faceID(sq, "B"),
			faceID(sq, "C"),
			faceID(sq, "D"),
			faceID(sq, "E"),
			faceID(sq, "F"),
		},
		[]ConnectionID{
			link("A", "B", 1, 3),
			link("B", "C", 0, 2),
			link("B", "D", 2, 0),
			link("B", "E", 1, 3),
			link("E", "F", 1, 3),

			flap("A", "F", 3, 1),
			flap("C", "A", 3, 0),
			flap("D", "A", 3, 2),
			flap("C", "E", 1, 0),
			flap("D", "E", 1, 2),
			flap("F", "C", 0, 0),
			flap("F", "D", 2, 2),
		},
	)
	return c.Render(b)
}

func RightTrianglePrism(b primitives.BBox, height, leg1, leg2 float64) []FoldablePattern {
	a := math.Sqrt(leg1*leg1 + leg2*leg2)
	c := NewCutOut(
		[]FaceID{
			faceID(Rectangle(leg1, height), "A"),
			faceID(Rectangle(leg2, height), "B"),
			faceID(Rectangle(a, height), "C"),
			faceID(Shape{
				Edges: []Edge{
					{
						Vector: primitives.Vector{X: leg2, Y: leg1},
					},
					{
						Vector: primitives.Vector{X: -leg2, Y: 0},
					},
					{
						Vector: primitives.Vector{X: 0, Y: -leg1},
					},
				},
			}, "B+1"),
			faceID(Shape{
				Edges: []Edge{
					{
						Vector: primitives.Vector{X: leg2, Y: 0},
					},
					{
						Vector: primitives.Vector{X: -leg2, Y: leg1},
					},
					{
						Vector: primitives.Vector{X: 0, Y: -leg1},
					},
				},
			}, "B-1"),
		},
		[]ConnectionID{
			link("A", "B", 1, 3),
			link("B", "C", 1, 3),
			link("B", "B+1", 0, 1),
			link("B", "B-1", 2, 0),

			flap("C", "A", 1, 3),

			flap("B+1", "A", 2, 0),
			flap("B-1", "A", 2, 2),

			flap("C", "B+1", 0, 0),
			flap("C", "B-1", 2, 1),
		},
	)
	return c.Render(b)
}

func Rhombicuboctahedron(b primitives.BBox, side float64) []FoldablePattern {
	sq := Square(side)
	tri := EquiTriangle(side)
	c := NewCutOut(
		[]FaceID{
			faceID(sq, "A"),
			faceID(sq, "B"),
			faceID(sq, "C"),
			faceID(sq, "D"),
			faceID(sq, "E"),
			faceID(sq, "F"),
			faceID(sq, "G"),
			faceID(sq, "H"),
			faceID(sq, "A+1"),
			faceID(sq, "A+2"),
			faceID(sq, "A-1"),
			faceID(sq, "A-2"),
			faceID(tri, "B+1"),
			faceID(tri, "B-1"),
			faceID(sq, "C+1"),
			faceID(sq, "C-1"),
			faceID(tri, "D+1"),
			faceID(tri, "D-1"),
			faceID(sq, "E+1"),
			faceID(sq, "E-1"),
			faceID(tri, "F+1"),
			faceID(tri, "F-1"),
			faceID(sq, "G+1"),
			faceID(sq, "G-1"),
			faceID(tri, "H+1"),
			faceID(tri, "H-1"),
		},
		[]ConnectionID{
			link("A", "B", 1, 3),
			link("B", "C", 1, 3),
			link("C", "D", 1, 3),
			link("D", "E", 1, 3),
			link("E", "F", 1, 3),
			link("F", "G", 1, 3),
			link("G", "H", 1, 3),

			link("A", "A+1", 0, 2),
			link("A+1", "A+2", 0, 2),
			link("A", "A-1", 2, 0),
			link("A-1", "A-2", 2, 0),

			link("B", "B+1", 0, 0),
			link("B", "B-1", 2, 0),

			link("C", "C+1", 0, 2),
			link("C", "C-1", 2, 0),

			link("D", "D+1", 0, 0),
			link("D", "D-1", 2, 0),

			link("E", "E+1", 0, 2),
			link("E", "E-1", 2, 0),

			link("F", "F+1", 0, 0),
			link("F", "F-1", 2, 0),

			link("G", "G+1", 0, 2),
			link("G", "G-1", 2, 0),

			link("H", "H+1", 0, 0),
			link("H", "H-1", 2, 0),

			flap("A", "H", 3, 1),
			flap("A+2", "G+1", 3, 0),
			flap("A+2", "C+1", 1, 0),
			flap("A-2", "G-1", 3, 2),
			flap("A-2", "C-1", 1, 2),

			smallFlap("A+1", "B+1", 1, 1),
			smallFlap("A-1", "B-1", 1, 2),

			smallFlap("B+1", "C+1", 2, 3),
			smallFlap("B-1", "C-1", 1, 3),

			smallFlap("C+1", "D+1", 1, 1),
			smallFlap("C-1", "D-1", 1, 2),

			smallFlap("D+1", "E+1", 2, 3),
			smallFlap("D-1", "E-1", 1, 3),

			smallFlap("E+1", "F+1", 1, 1),
			smallFlap("E-1", "F-1", 1, 2),

			smallFlap("F+1", "G+1", 2, 3),
			smallFlap("F-1", "G-1", 1, 3),

			smallFlap("G+1", "H+1", 1, 1),
			smallFlap("G-1", "H-1", 1, 2),

			smallFlap("H+1", "A+1", 2, 3),
			smallFlap("H-1", "A-1", 1, 3),

			flap("E+1", "A+2", 0, 0),
			flap("E-1", "A-2", 2, 2),
		},
	)
	return c.Render(b)
}

// RhombicuboctahedronWithoutCorners is the same as RhombicuboctahedronID, but with the triangular corner pieces missing
// The idea is to have a right-angle corner inserts in each space, but this requires a disconneced foldable
func RhombicuboctahedronWithoutCorners(b primitives.BBox, side float64) []FoldablePattern {
	// TODO: Add in corner pieces
	// fillSpacing := 20.0
	triColor := "red"
	faceColor := "yellow"
	redPen := pen.BicIntensityBrushTip
	yellowPen := pen.BicIntensityBrushTip
	sq := Square(side)
	tr := Shape{
		Edges: []Edge{
			{
				Vector: primitives.Vector{X: side / 2, Y: -side / 2},
			},
			{
				Vector: primitives.Vector{X: side / 2, Y: side / 2},
			},
			{
				Vector: primitives.Vector{X: -side, Y: 0},
			},
		},
	}
	c := NewCutOut(
		[]FaceID{
			faceID(sq, "A").WithBrushFill(faceColor, yellowPen),
			faceID(sq, "B"),
			faceID(sq, "C").WithBrushFill(faceColor, yellowPen),
			faceID(sq, "D"),
			faceID(sq, "E").WithBrushFill(faceColor, yellowPen),
			faceID(sq, "F"),
			faceID(sq, "G").WithBrushFill(faceColor, yellowPen),
			faceID(sq, "H"),
			faceID(sq, "A+1"),
			faceID(sq, "A+2").WithBrushFill(faceColor, yellowPen),
			faceID(sq, "A-1"),
			faceID(sq, "A-2").WithBrushFill(faceColor, yellowPen),
			faceID(sq, "C+1"),
			faceID(sq, "C-1"),
			faceID(sq, "E+1"),
			faceID(sq, "E-1"),
			faceID(sq, "G+1"),
			faceID(sq, "G-1"),

			faceID(tr, "B+a").WithBrushFill(triColor, redPen),
			faceID(tr, "B+b").WithBrushFill(triColor, redPen),
			faceID(tr, "B+c").WithBrushFill(triColor, redPen),

			faceID(tr, "B-a").WithBrushFill(triColor, redPen),
			faceID(tr, "B-b").WithBrushFill(triColor, redPen),
			faceID(tr, "B-c").WithBrushFill(triColor, redPen),

			faceID(tr, "D+a").WithBrushFill(triColor, redPen),
			faceID(tr, "D+b").WithBrushFill(triColor, redPen),
			faceID(tr, "D+c").WithBrushFill(triColor, redPen),

			faceID(tr, "D-a").WithBrushFill(triColor, redPen),
			faceID(tr, "D-b").WithBrushFill(triColor, redPen),
			faceID(tr, "D-c").WithBrushFill(triColor, redPen),

			faceID(tr, "F+a").WithBrushFill(triColor, redPen),
			faceID(tr, "F+b").WithBrushFill(triColor, redPen),
			faceID(tr, "F+c").WithBrushFill(triColor, redPen),

			faceID(tr, "F-a").WithBrushFill(triColor, redPen),
			faceID(tr, "F-b").WithBrushFill(triColor, redPen),
			faceID(tr, "F-c").WithBrushFill(triColor, redPen),

			faceID(tr, "H+a").WithBrushFill(triColor, redPen),
			faceID(tr, "H+b").WithBrushFill(triColor, redPen),
			faceID(tr, "H+c").WithBrushFill(triColor, redPen),

			faceID(tr, "H-a").WithBrushFill(triColor, redPen),
			faceID(tr, "H-b").WithBrushFill(triColor, redPen),
			faceID(tr, "H-c").WithBrushFill(triColor, redPen),
		},
		[]ConnectionID{
			link("A", "B", 1, 3),
			link("B", "C", 1, 3),
			link("C", "D", 1, 3),
			link("D", "E", 1, 3),
			link("E", "F", 1, 3),
			link("F", "G", 1, 3),
			link("G", "H", 1, 3),

			link("A", "A+1", 0, 2),
			link("A+1", "A+2", 0, 2),
			link("A", "A-1", 2, 0),
			link("A-1", "A-2", 2, 0),

			link("C", "C+1", 0, 2),
			link("C", "C-1", 2, 0),

			link("E", "E+1", 0, 2),
			link("E", "E-1", 2, 0),

			link("G", "G+1", 0, 2),
			link("G", "G-1", 2, 0),

			link("B+a", "B+b", 0, 1),
			link("B+b", "B+c", 0, 1),

			flap("B+a", "B", 2, 0),
			flap("B+b", "A+1", 2, 3),
			flap("B+c", "C+1", 2, 1),
			flap("B+a", "B+c", 1, 0),

			link("B-a", "B-b", 1, 0),
			link("B-b", "B-c", 1, 0),

			flap("B-a", "B", 2, 2),
			flap("B-b", "A-1", 2, 3),
			flap("B-c", "C-1", 2, 1),
			flap("B-a", "B-c", 0, 1),

			link("D+a", "D+b", 0, 1),
			link("D+b", "D+c", 0, 1),

			flap("D+a", "D", 2, 0),
			flap("D+b", "C+1", 2, 3),
			flap("D+c", "E+1", 2, 1),
			flap("D+a", "D+c", 1, 0),

			link("D-a", "D-b", 1, 0),
			link("D-b", "D-c", 1, 0),

			flap("D-a", "D", 2, 2),
			flap("D-b", "C-1", 2, 3),
			flap("D-c", "E-1", 2, 1),
			flap("D-a", "D-c", 0, 1),

			link("F+a", "F+b", 0, 1),
			link("F+b", "F+c", 0, 1),

			flap("F+a", "F", 2, 0),
			flap("F+b", "E+1", 2, 3),
			flap("F+c", "G+1", 2, 1),
			flap("F+a", "F+c", 1, 0),

			link("F-a", "F-b", 1, 0),
			link("F-b", "F-c", 1, 0),

			flap("F-a", "F", 2, 2),
			flap("F-b", "E-1", 2, 3),
			flap("F-c", "G-1", 2, 1),
			flap("F-a", "F-c", 0, 1),

			link("H+a", "H+b", 0, 1),
			link("H+b", "H+c", 0, 1),

			flap("H+a", "H", 2, 0),
			flap("H+b", "G+1", 2, 3),
			flap("H+c", "A+1", 2, 1),
			flap("H+a", "H+c", 1, 0),

			link("H-a", "H-b", 1, 0),
			link("H-b", "H-c", 1, 0),

			flap("H-a", "H", 2, 2),
			flap("H-b", "G-1", 2, 3),
			flap("H-c", "A-1", 2, 1),
			flap("H-a", "H-c", 0, 1),

			flap("A", "H", 3, 1),
			flap("A+2", "G+1", 3, 0),
			flap("A+2", "C+1", 1, 0),
			flap("A-2", "G-1", 3, 2),
			flap("A-2", "C-1", 1, 2),

			flap("E+1", "A+2", 0, 0),
			flap("E-1", "A-2", 2, 2),
		},
	)
	return c.Render(b)
}

func CutCube(b primitives.BBox, side float64, cutRatio float64) []FoldablePattern {
	sq := Square(side)
	a := math.Sqrt(1 + cutRatio*cutRatio)
	c := NewCutOut(
		[]FaceID{
			faceID(sq, "A"),
			faceID(sq, "B"),
			faceID(Rectangle((1-cutRatio)*side, side), "C"),
			faceID(Rectangle(a*side, side), "D"),
			faceID(Shape{
				Edges: []Edge{
					{
						Vector: primitives.Vector{X: 1, Y: cutRatio}.Mult(side),
					},
					{
						Vector: primitives.Vector{X: 0, Y: 1 - cutRatio}.Mult(side),
					},
					{
						Vector: primitives.Vector{X: -1, Y: 0}.Mult(side),
					},
					{
						Vector: primitives.Vector{X: 0, Y: -1}.Mult(side),
					},
				},
			}, "B+1"),
			faceID(Shape{
				Edges: []Edge{
					{
						Vector: primitives.Vector{X: 1, Y: 0}.Mult(side),
					},
					{
						Vector: primitives.Vector{X: 0, Y: 1 - cutRatio}.Mult(side),
					},
					{
						Vector: primitives.Vector{X: -1, Y: cutRatio}.Mult(side),
					},
					{
						Vector: primitives.Vector{X: 0, Y: -1}.Mult(side),
					},
				},
			}, "B-1"),
		},
		[]ConnectionID{
			link("A", "B", 1, 3),
			link("B", "C", 1, 3),
			link("C", "D", 1, 3),
			link("B", "B+1", 0, 2),
			link("B", "B-1", 2, 0),

			flap("A", "D", 3, 1),
			flap("B+1", "A", 3, 0),
			flap("B-1", "A", 3, 2),

			flap("B+1", "C", 1, 0),
			flap("B-1", "C", 1, 2),

			flap("D", "B+1", 0, 0),
			flap("D", "B-1", 2, 2),
		},
	)
	return c.Render(b)
}
