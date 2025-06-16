package foldable

import (
	"math"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/fonts"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

// Cube does a cube with a certain side length.
func Cube(b box.Box, side float64) []lines.LineLike {
	sq := Square(side)
	c := NewFace(sq).WithFlap(3).WithFace( // #0
		1,
		NewFace(sq).WithFace( // #1
			1,
			NewFace(sq).WithFace( // #2
				1,
				NewFace(sq).WithFlap(0).WithFlap(2), // #3
				3,
			),
			3,
		).WithFace(
			0,
			NewFace(sq).WithFlap(1).WithFlap(3), // #4
			2,
		).WithFace(
			2,
			NewFace(sq).WithFlap(1).WithFlap(3), // #5
			0,
		),
		3,
	)
	return c.Render(primitives.Point{X: b.X, Y: b.Y + side}, 0)
}

func RightTrianglePrism(b box.Box, height, leg1, leg2 float64) []lines.LineLike {
	a := math.Sqrt(leg1*leg1 + leg2*leg2)
	c := NewFace(Rectangle(leg1, height)).WithFlap(3).WithFace( // #0
		1,
		NewFace(Rectangle(leg2, height)).WithFace( // #1
			1,
			NewFace(Rectangle(a, height)).WithFlap(0).WithFlap(2), // #2
			3,
		).WithFace(
			0,
			NewFace(Shape{
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
			}).WithFlap(2), // #3
			1,
		).WithFace(
			2,
			NewFace(Shape{
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
			}).WithFlap(2), // #4
			0,
		),
		3,
	)
	return c.Render(primitives.Point{X: b.X, Y: b.Y + leg1}, 0)
}

func ShapeTester(b box.Box, side float64) []lines.LineLike {
	sq := Square(side)
	tri := EquiTriangle(side)

	c := NewFace(sq).WithFlap(3).
		WithFace(
			0,
			NewFace(tri).WithFlap(0).WithFace(1, NewFace(sq).WithFlap(3).WithFace(0, NewFace(tri).WithFlap(1), 0), 2),
			2,
		)
	return c.Render(primitives.Point{X: b.X, Y: b.Y + side*2}, 0)
}

func Rhombicuboctahedron(b box.Box, side float64) []lines.LineLike {
	sq := Square(side)
	tri := EquiTriangle(side)

	c := NewFace(sq).WithFlap(3).
		WithFace( // #0
			0,
			NewFace(sq).WithSmallFlap(1).
				WithFace( // #1
					0,
					NewFace(sq).WithFlap(1).WithFlap(3), // #2
					2,
				),
			2,
		).
		WithFace(
			2,
			NewFace(sq).
				WithFace(
					2,
					NewFace(sq).WithFlap(1).WithFlap(3),
					0,
				),
			0,
		).WithFace(
		1,
		NewFace(sq).
			WithFace(
				0,
				NewFace(tri).WithSmallFlap(1),
				2,
			).
			WithFace(
				2,
				NewFace(tri).WithSmallFlap(1),
				0,
			).
			WithFace(
				1,
				NewFace(sq).
					WithFace(
						0,
						NewFace(sq).WithSmallFlap(1),
						2,
					).
					WithFace(
						2,
						NewFace(sq).WithSmallFlap(1),
						0,
					).
					WithFace(
						1,
						NewFace(sq).
							WithFace(
								0,
								NewFace(tri).WithSmallFlap(1),
								2,
							).
							WithFace(
								2,
								NewFace(tri).WithSmallFlap(1),
								0,
							).
							WithFace(
								1,
								NewFace(sq).
									WithFace(
										0,
										NewFace(sq).WithFlap(0).WithSmallFlap(1),
										2,
									).
									WithFace(
										2,
										NewFace(sq).WithFlap(2).WithSmallFlap(1),
										0).
									WithFace(
										1,
										NewFace(sq).
											WithFace(
												0,
												NewFace(tri).WithSmallFlap(1),
												2,
											).
											WithFace(
												2,
												NewFace(tri).WithSmallFlap(1),
												0,
											).
											WithFace(
												1,
												NewFace(sq).
													WithFace(
														0,
														NewFace(sq).WithSmallFlap(1),
														2,
													).
													WithFace(
														2,
														NewFace(sq).WithSmallFlap(1),
														0,
													).
													WithFace(
														1,
														NewFace(sq).WithFlap(1).
															WithFace(
																0,
																NewFace(tri).WithSmallFlap(1),
																2,
															).
															WithFace(
																2,
																NewFace(tri).WithSmallFlap(1),
																0,
															),
														3),
												3),
										3),
								3),
						3),
				3),
		3).
		WithFace(
			2,
			NewFace(sq).WithSmallFlap(1).
				WithFace(
					2,
					NewFace(sq).WithFlap(1).WithFlap(3),
					0,
				),
			0,
		)
	return c.Render(primitives.Point{X: b.X - 2_000, Y: b.Y + side*2}, 0)
}

func RhombicuboctahedronWithoutCorners(b box.Box, side float64) FoldablePattern {
	sq := Square(side)
	// tri := EquiTriangle(side)

	c := NewFace(sq).WithFlap(3).
		WithFace( // #0
			0,
			NewFace(sq).WithSmallFlap(1).
				WithFace( // #1
					0,
					NewFace(sq).WithFlap(1).WithFlap(3), // #2
					2,
				),
			2,
		).
		WithFace(
			2,
			NewFace(sq).
				WithFace(
					2,
					NewFace(sq).WithFlap(1).WithFlap(3),
					0,
				),
			0,
		).WithFace(
		1,
		NewFace(sq).
			// WithFace(
			// 	0,
			// 	NewFace(tri).WithSmallFlap(1),
			// 	2,
			// ).
			// WithFace(
			// 	2,
			// 	NewFace(tri).WithSmallFlap(1),
			// 	0,
			// ).
			WithFace(
				1,
				NewFace(sq).
					WithFace(
						0,
						NewFace(sq).WithSmallFlap(1),
						2,
					).
					WithFace(
						2,
						NewFace(sq).WithSmallFlap(1),
						0,
					).
					WithFace(
						1,
						NewFace(sq).
							// WithFace(
							// 	0,
							// 	NewFace(tri).WithSmallFlap(1),
							// 	2,
							// ).
							// WithFace(
							// 	2,
							// 	NewFace(tri).WithSmallFlap(1),
							// 	0,
							// ).
							WithFace(
								1,
								NewFace(sq).
									WithFace(
										0,
										NewFace(sq).WithFlap(0).WithSmallFlap(1),
										2,
									).
									WithFace(
										2,
										NewFace(sq).WithFlap(2).WithSmallFlap(1),
										0).
									WithFace(
										1,
										NewFace(sq).
											// WithFace(
											// 	0,
											// 	NewFace(tri).WithSmallFlap(1),
											// 	2,
											// ).
											// WithFace(
											// 	2,
											// 	NewFace(tri).WithSmallFlap(1),
											// 	0,
											// ).
											WithFace(
												1,
												NewFace(sq).
													WithFace(
														0,
														NewFace(sq).WithSmallFlap(1),
														2,
													).
													WithFace(
														2,
														NewFace(sq).WithSmallFlap(1),
														0,
													).
													WithFace(
														1,
														NewFace(sq).WithFlap(1),
														// WithFace(
														// 	0,
														// 	NewFace(tri).WithSmallFlap(1),
														// 	2,
														// ).
														// WithFace(
														// 	2,
														// 	NewFace(tri).WithSmallFlap(1),
														// 	0,
														// ),
														3),
												3),
										3),
								3),
						3),
				3),
		3).
		WithFace(
			2,
			NewFace(sq).WithSmallFlap(1).
				WithFace(
					2,
					NewFace(sq).WithFlap(1).WithFlap(3),
					0,
				),
			0,
		)
	return FoldablePattern{
		Edges:       c.Render(primitives.Point{X: b.X - 2_000, Y: b.Y + side*2}, 0),
		Annotations: fonts.RenderText(b, "ABC").CharCurves,
	}
}

// CutCube is a cube with a triangular prism missing.
func CutCube(b box.Box, side float64, cutRatio float64) []lines.LineLike {
	sq := Square(side)
	a := math.Sqrt(1 + cutRatio*cutRatio)
	c := NewFace(sq).WithFlap(3).WithFace( // #0
		1,
		NewFace(sq).WithFace( // #1
			1,
			NewFace(Shape{
				Edges: []Edge{
					{
						Vector: primitives.Vector{X: 1 - cutRatio, Y: 0}.Mult(side),
					},
					{
						Vector: primitives.Vector{X: 0, Y: 1}.Mult(side),
					},
					{
						Vector: primitives.Vector{X: -(1 - cutRatio), Y: 0}.Mult(side),
					},
					{
						Vector: primitives.Vector{X: 0, Y: -1}.Mult(side),
					},
				},
			}).WithFace( // #2
				1,
				NewFace(Shape{
					Edges: []Edge{
						{
							Vector: primitives.Vector{X: a, Y: 0}.Mult(side),
						},
						{
							Vector: primitives.Vector{X: 0, Y: 1}.Mult(side),
						},
						{
							Vector: primitives.Vector{X: -a, Y: 0}.Mult(side),
						},
						{
							Vector: primitives.Vector{X: 0, Y: -1}.Mult(side),
						},
					},
				}).WithFlap(0).WithFlap(2), // #3
				3,
			),
			3,
		).WithFace(
			0,
			NewFace(Shape{
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
			}).WithFlap(1).WithFlap(3), // #4
			2,
		).WithFace(
			2,
			NewFace(Shape{
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
			}).WithFlap(1).WithFlap(3), // #5
			0,
		),
		3,
	)
	return c.Render(primitives.Point{X: b.X, Y: b.Y + side}, 0)
}

// ManualCube is deprecated, use Cube instead, it's more generic
func ManualCube(b box.Box, side float64) []lines.LineLike {
	start := primitives.Point{X: b.X, Y: b.Y + side}
	lns := []lines.LineLike{}
	// draws the cube as follows:
	//
	//     +---+
	//  /-\| 4 |/-\ /-\
	// +---+---+---+---+\
	// | 0 | 1 | 2 | 3 ||
	// +---+---+---+---+/
	//  \-/| 5 |\-/ \-/
	//     +---+

	// face 0
	l := lines.NewPath(start)
	end := start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: -side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 1
	start = primitives.Point{X: b.X + side, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 2
	start = primitives.Point{X: b.X + side*2, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 3
	start = primitives.Point{X: b.X + side*3, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 4
	start = primitives.Point{X: b.X + side, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: 0, Y: -side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 5
	start = primitives.Point{X: b.X + side, Y: b.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: -side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 04
	start = primitives.Point{X: b.X, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 24
	start = primitives.Point{X: b.X + side*2, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 34
	start = primitives.Point{X: b.X + side*3, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 05
	start = primitives.Point{X: b.X, Y: b.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 25
	start = primitives.Point{X: b.X + side*2, Y: b.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 35
	start = primitives.Point{X: b.X + side*3, Y: b.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 30
	start = primitives.Point{X: b.X + side*4, Y: b.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side - flapWidth*2})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	return lns
}
