package foldable

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/pen"
	"github.com/libeks/go-plotter-svg/primitives"
	"github.com/libeks/go-plotter-svg/voronoi"
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

// RhombicuboctahedronWithoutCorners is the same as RhombicuboctahedronID, but with the triangular corner pieces missing
// The idea is to have a right-angle corner inserts in each space, but this requires a disconneced foldable
func RhombicuboctahedronWithoutCornersTricolor(b primitives.BBox, side float64) []FoldablePattern {
	// TODO: Add in corner pieces
	// fillSpacing := 20.0
	cyanColor := "cyan"
	magentaColor := "magenta"
	yellowColor := "yellow"
	cyanPen := pen.BicIntensityBrushTip
	magentaPen := pen.BicIntensityBrushTip
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
			faceID(sq, "A").WithBrushFill(cyanColor, cyanPen).WithBrushFill(magentaColor, magentaPen),
			faceID(sq, "B").WithBrushFill(magentaColor, magentaPen),
			faceID(sq, "C").WithBrushFill(magentaColor, magentaPen).WithBrushFill(yellowColor, yellowPen),
			faceID(sq, "D").WithBrushFill(magentaColor, magentaPen),
			faceID(sq, "E").WithBrushFill(cyanColor, cyanPen).WithBrushFill(magentaColor, magentaPen),
			faceID(sq, "F").WithBrushFill(magentaColor, magentaPen),
			faceID(sq, "G").WithBrushFill(magentaColor, magentaPen).WithBrushFill(yellowColor, yellowPen),
			faceID(sq, "H").WithBrushFill(magentaColor, magentaPen),
			faceID(sq, "A+1").WithBrushFill(cyanColor, cyanPen),
			faceID(sq, "A+2").WithBrushFill(cyanColor, cyanPen).WithBrushFill(yellowColor, yellowPen),
			faceID(sq, "A-1").WithBrushFill(cyanColor, cyanPen),
			faceID(sq, "A-2").WithBrushFill(cyanColor, cyanPen).WithBrushFill(yellowColor, yellowPen),
			faceID(sq, "C+1").WithBrushFill(yellowColor, yellowPen),
			faceID(sq, "C-1").WithBrushFill(yellowColor, yellowPen),
			faceID(sq, "E+1").WithBrushFill(cyanColor, cyanPen),
			faceID(sq, "E-1").WithBrushFill(cyanColor, cyanPen),
			faceID(sq, "G+1").WithBrushFill(yellowColor, yellowPen),
			faceID(sq, "G-1").WithBrushFill(yellowColor, yellowPen),

			faceID(tr, "B+a").WithBrushFill(magentaColor, magentaPen),
			faceID(tr, "B+b").WithBrushFill(cyanColor, cyanPen),
			faceID(tr, "B+c").WithBrushFill(yellowColor, yellowPen),

			faceID(tr, "B-a").WithBrushFill(magentaColor, magentaPen),
			faceID(tr, "B-b").WithBrushFill(cyanColor, cyanPen),
			faceID(tr, "B-c").WithBrushFill(yellowColor, yellowPen),

			faceID(tr, "D+a").WithBrushFill(magentaColor, magentaPen),
			faceID(tr, "D+b").WithBrushFill(yellowColor, yellowPen),
			faceID(tr, "D+c").WithBrushFill(cyanColor, cyanPen),

			faceID(tr, "D-a").WithBrushFill(magentaColor, magentaPen),
			faceID(tr, "D-b").WithBrushFill(yellowColor, yellowPen),
			faceID(tr, "D-c").WithBrushFill(cyanColor, cyanPen),

			faceID(tr, "F+a").WithBrushFill(magentaColor, magentaPen),
			faceID(tr, "F+b").WithBrushFill(cyanColor, cyanPen),
			faceID(tr, "F+c").WithBrushFill(yellowColor, yellowPen),

			faceID(tr, "F-a").WithBrushFill(magentaColor, magentaPen),
			faceID(tr, "F-b").WithBrushFill(cyanColor, cyanPen),
			faceID(tr, "F-c").WithBrushFill(yellowColor, yellowPen),

			faceID(tr, "H+a").WithBrushFill(magentaColor, magentaPen),
			faceID(tr, "H+b").WithBrushFill(yellowColor, yellowPen),
			faceID(tr, "H+c").WithBrushFill(cyanColor, cyanPen),

			faceID(tr, "H-a").WithBrushFill(magentaColor, magentaPen),
			faceID(tr, "H-b").WithBrushFill(yellowColor, yellowPen),
			faceID(tr, "H-c").WithBrushFill(cyanColor, cyanPen),
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

func VoronoiFoldable(b primitives.BBox) []FoldablePattern {
	bbox := primitives.BBox{UpperLeft: primitives.Origin, LowerRight: primitives.Point{X: 5000, Y: 5000}}
	edgeWidth := 400.0
	// points := []primitives.Point{
	// 	{X: 1000, Y: 1000},
	// 	{X: 1500, Y: 2000},
	// 	{X: 4500, Y: 1500},
	// 	{X: 3000, Y: 4500},
	// }
	// points := []primitives.Point{
	// 	{X: 2803.0, Y: 1061.1},
	// 	{X: 2288.1, Y: 3567.1},
	// 	{X: 2800.8, Y: 3797.2},
	// 	{X: 3377.6, Y: 1596.4},
	// }
	// points := []primitives.Point{
	// 	{3939.3, 3469.2},
	// 	{3860.3, 2980.4},
	// 	{1301.9, 2714.9},
	// 	{2320.7, 3376.5},
	// 	{2884.2, 2017.6},
	// 	{1692.4, 1053.9},
	// }
	// points := []primitives.Point{
	// 	{2629.9, 1649.7},
	// 	{3400.8, 2265.2},
	// 	{2385.6, 2281.2},
	// 	{1485.5, 1855.5},
	// 	{1806.5, 2754.3},
	// 	{2038.0, 3307.8},
	// }
	// points := []primitives.Point{
	// 	{1962.3, 3098.7},
	// 	{2962.1, 1710.6},
	// 	{688.7, 4222.8},
	// 	{2804.3, 1045.3},
	// 	{2904.4, 579.4},
	// }
	// points := []primitives.Point{
	// 	{1685.5, 1226.8},
	// 	{1396.8, 3673.7},
	// 	{3100.1, 2711.9},
	// 	{2300.4, 2249.0},
	// 	{3378.5, 712.4},
	// }
	// points := []primitives.Point{
	// 	{2026.8, 1272.8},
	// 	{3389.8, 2680.9},
	// 	{4302.4, 1969.8},
	// 	{4308.0, 3301.0},
	// 	{4449.9, 3944.1},
	// }
	// points := []primitives.Point{
	// 	{1577.2, 987.1},
	// 	{1932.1, 3904.1},
	// 	{759.6, 4256.4},
	// 	{1427.1, 2965.9},
	// 	{3026.7, 1864.5},
	// 	{1551.9, 2038.8},
	// 	{2963.5, 1649.6},
	// }
	nPoints := 7
	points := make([]primitives.Point, nPoints)
	fmt.Printf("points := []primitives.Point{\n")
	for i := range nPoints {
		points[i] = primitives.Point{
			X: rand.Float64()*4000 + 500.0,
			Y: rand.Float64()*4000 + 500.0,
		}
		fmt.Printf("    {%.1f, %.1f},\n", points[i].X, points[i].Y)
	}
	fmt.Printf("}\n")

	vor := voronoi.ComputeVoronoiConnections(bbox, points)
	faces := []FaceID{}
	connections := []ConnectionID{}
	faceNames := make([]string, len(vor.Polygons))
	edgesVisited := make(map[string]bool)
	minEdge := math.MaxFloat64
	for i, poly := range vor.Polygons {
		name := fmt.Sprintf("%d", i)
		shape := PolygonToShape(poly)
		for _, edge := range shape.Edges {
			length := edge.Len()
			if length < minEdge {
				minEdge = length
			}
		}
		faces = append(faces, faceID(shape, name))
		faceNames[i] = name
		for j := range len(poly.Points) {
			edgesVisited[fmt.Sprintf("%s-%d", name, j)] = false
		}
	}
	fmt.Printf("Min Edge is %.1f\n", minEdge)
	for _, conn := range vor.EdgeMap {
		connections = append(connections, ConnectionID{
			FaceA:          faceNames[conn.From.PolyIndex],
			FaceB:          faceNames[conn.To.PolyIndex],
			EdgeAID:        conn.From.EdgeIndex,
			EdgeBID:        conn.To.EdgeIndex,
			ConnectionType: DoubleConnection,
		})
		edgesVisited[fmt.Sprintf("%s-%d", faceNames[conn.From.PolyIndex], conn.From.EdgeIndex)] = true
		edgesVisited[fmt.Sprintf("%s-%d", faceNames[conn.To.PolyIndex], conn.To.EdgeIndex)] = true
	}
	for i, poly := range vor.Polygons {
		shape := PolygonToShape(poly)
		for j, edge := range shape.Edges {
			faceName := faceNames[i]
			edgeName := fmt.Sprintf("%s-%d", faceName, j)
			if visited := edgesVisited[edgeName]; !visited {
				// this is an outer edge
				edgeShape := Shape{
					Edges: []Edge{
						{edge.Mult(-1)},
						{edge.Perp().Unit().Mult(-edgeWidth)},
						edge,
						{edge.Perp().Unit().Mult(edgeWidth)},
					},
				}
				faces = append(faces, faceID(edgeShape, edgeName))
				// connect the edge to the original face
				connections = append(connections, ConnectionID{
					FaceA:          faceName,
					FaceB:          edgeName,
					EdgeAID:        j,
					EdgeBID:        0,
					ConnectionType: FaceConnection,
				})
				// mark the outer edge as not connected to anything
				connections = append(connections, ConnectionID{
					FaceA:          edgeName,
					FaceB:          "",
					EdgeAID:        2,
					EdgeBID:        -1,
					ConnectionType: NoneConnection,
				})
				nextEdgeName := fmt.Sprintf("%s-%d", faceName, maths.Mod(j+1, len(shape.Edges)))
				if visited := edgesVisited[nextEdgeName]; !visited {
					// if the next edge is also not connected to anything, add a flap between the edge elements
					connections = append(connections, ConnectionID{
						FaceA:          edgeName,
						FaceB:          nextEdgeName,
						EdgeAID:        3,
						EdgeBID:        1,
						ConnectionType: FlapConnection,
					})
				}
			}
		}
	}
	for _, conn := range vor.EdgeMap {
		faceA := faceNames[conn.From.PolyIndex]
		faceB := faceNames[conn.To.PolyIndex]
		faceAPreviousEdgeIndex := vor.Polygons[conn.From.PolyIndex].PreviousFaceIndex(conn.From.EdgeIndex)
		faceANextEdgeIndex := vor.Polygons[conn.From.PolyIndex].NextFaceIndex(conn.From.EdgeIndex)
		faceAPreviousEdge := fmt.Sprintf("%s-%d", faceA, faceAPreviousEdgeIndex)
		faceANextEdge := fmt.Sprintf("%s-%d", faceA, faceANextEdgeIndex)
		faceBPreviousEdgeIndex := vor.Polygons[conn.To.PolyIndex].PreviousFaceIndex(conn.To.EdgeIndex)
		faceBNextEdgeIndex := vor.Polygons[conn.To.PolyIndex].NextFaceIndex(conn.To.EdgeIndex)
		if visited := edgesVisited[faceAPreviousEdge]; !visited {
			// if faces are connected, and the previous edge is not connected to anything, add a flap between the edge elements
			connections = append(connections, ConnectionID{
				FaceA:          fmt.Sprintf("%s-%d", faceA, faceAPreviousEdgeIndex),
				FaceB:          fmt.Sprintf("%s-%d", faceB, faceBNextEdgeIndex),
				EdgeAID:        3,
				EdgeBID:        1,
				ConnectionType: FlapConnection,
			})
		}
		if visited := edgesVisited[faceANextEdge]; !visited {
			// if faces are connected, and the next edge is not connected to anything, add a flap between the edge elements
			connections = append(connections, ConnectionID{
				FaceA:          fmt.Sprintf("%s-%d", faceA, faceANextEdgeIndex),
				FaceB:          fmt.Sprintf("%s-%d", faceB, faceBPreviousEdgeIndex),
				EdgeAID:        1,
				EdgeBID:        3,
				ConnectionType: FlapConnection,
			})
		}
	}
	return NewCutOut(faces, connections).Render(b)
}
