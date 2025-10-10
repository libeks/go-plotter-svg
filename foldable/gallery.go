package foldable

import (
	"math"

	"github.com/libeks/go-plotter-svg/fonts"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

// Cube does a cube with a certain side length.
func Cube(b primitives.BBox, side float64) []lines.LineLike {
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
	return c.Render(primitives.Point{X: b.UpperLeft.X, Y: b.UpperLeft.Y + side}, 0).Lines
}

func CubeID(b primitives.BBox, side float64) FoldablePattern {
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
			{
				FaceA:          "A",
				FaceB:          "B",
				EdgeAID:        1,
				EdgeBID:        3,
				ConnectionType: FaceConnection,
			},
			{
				FaceA:          "B",
				FaceB:          "C",
				EdgeAID:        0,
				EdgeBID:        2,
				ConnectionType: FaceConnection,
			},
			{
				FaceA:          "B",
				FaceB:          "D",
				EdgeAID:        2,
				EdgeBID:        0,
				ConnectionType: FaceConnection,
			},
			{
				FaceA:          "B",
				FaceB:          "E",
				EdgeAID:        1,
				EdgeBID:        3,
				ConnectionType: FaceConnection,
			},
			{
				FaceA:          "E",
				FaceB:          "F",
				EdgeAID:        1,
				EdgeBID:        3,
				ConnectionType: FaceConnection,
			},
			{
				FaceA:          "A",
				FaceB:          "F",
				EdgeAID:        3,
				EdgeBID:        1,
				ConnectionType: FlapConnection,
			},
			{
				FaceA:          "C",
				FaceB:          "A",
				EdgeAID:        3,
				EdgeBID:        0,
				ConnectionType: FlapConnection,
			},
			{
				FaceA:          "D",
				FaceB:          "A",
				EdgeAID:        3,
				EdgeBID:        2,
				ConnectionType: FlapConnection,
			},
			{
				FaceA:          "C",
				FaceB:          "E",
				EdgeAID:        1,
				EdgeBID:        0,
				ConnectionType: FlapConnection,
			},
			{
				FaceA:          "D",
				FaceB:          "E",
				EdgeAID:        1,
				EdgeBID:        2,
				ConnectionType: FlapConnection,
			},
			{
				FaceA:          "F",
				FaceB:          "C",
				EdgeAID:        0,
				EdgeBID:        0,
				ConnectionType: FlapConnection,
			},
			{
				FaceA:          "F",
				FaceB:          "D",
				EdgeAID:        2,
				EdgeBID:        2,
				ConnectionType: FlapConnection,
			},
		},
	)
	return c.Render(b)
}

func RightTrianglePrism(b primitives.BBox, height, leg1, leg2 float64) []lines.LineLike {
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
	return c.Render(b.UpperLeft.Add(primitives.Vector{X: 0, Y: leg1}), 0).Lines
}

func ShapeTester(b primitives.BBox, side float64) []lines.LineLike {
	sq := Square(side)
	tri := EquiTriangle(side)

	c := NewFace(sq).WithFlap(3).
		WithFace(
			0,
			NewFace(tri).WithFlap(0).WithFace(1, NewFace(sq).WithFlap(3).WithFace(0, NewFace(tri).WithFlap(1), 0), 2),
			2,
		)
	return c.Render(primitives.Point{X: b.UpperLeft.X, Y: b.UpperLeft.Y + side*2}, 0).Lines
}

func Rhombicuboctahedron(b primitives.BBox, side float64) []lines.LineLike {
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
	return c.Render(primitives.Point{X: b.UpperLeft.X - 2_000, Y: b.UpperLeft.Y + side*2}, 0).Lines
}

func RhombicuboctahedronWithoutCorners(b primitives.BBox, side float64) FoldablePattern {
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
		Edges:       c.Render(primitives.Point{X: b.UpperLeft.X - 2_000, Y: b.UpperLeft.Y + side*2}, 0).Lines,
		Annotations: fonts.RenderText(b, "ABC").CharCurves,
	}
}

// CutCube is a cube with a triangular prism missing.
func CutCube(b primitives.BBox, side float64, cutRatio float64) []lines.LineLike {
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
	return c.Render(primitives.Point{X: b.UpperLeft.X, Y: b.UpperLeft.Y + side}, 0).Lines
}

// ManualCube is deprecated, use Cube instead, it's more generic
func ManualCube(b primitives.BBox, side float64) []lines.LineLike {
	start := primitives.Point{X: b.UpperLeft.X, Y: b.UpperLeft.Y + side}
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
	start = primitives.Point{X: b.UpperLeft.X + side, Y: b.UpperLeft.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 2
	start = primitives.Point{X: b.UpperLeft.X + side*2, Y: b.UpperLeft.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 3
	start = primitives.Point{X: b.UpperLeft.X + side*3, Y: b.UpperLeft.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: -side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 4
	start = primitives.Point{X: b.UpperLeft.X + side, Y: b.UpperLeft.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: 0, Y: -side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// face 5
	start = primitives.Point{X: b.UpperLeft.X + side, Y: b.UpperLeft.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: 0, Y: side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: 0, Y: -side})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 04
	start = primitives.Point{X: b.UpperLeft.X, Y: b.UpperLeft.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 24
	start = primitives.Point{X: b.UpperLeft.X + side*2, Y: b.UpperLeft.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 34
	start = primitives.Point{X: b.UpperLeft.X + side*3, Y: b.UpperLeft.Y + side}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 05
	start = primitives.Point{X: b.UpperLeft.X, Y: b.UpperLeft.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 25
	start = primitives.Point{X: b.UpperLeft.X + side*2, Y: b.UpperLeft.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 35
	start = primitives.Point{X: b.UpperLeft.X + side*3, Y: b.UpperLeft.Y + side*2}
	l = lines.NewPath(start)
	end = start.Add(primitives.Vector{X: flapWidth, Y: flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: side - flapWidth*2, Y: 0})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	start, end = end, end.Add(primitives.Vector{X: flapWidth, Y: -flapWidth})
	l = l.AddPathChunk(lines.LineChunk{Start: start, End: end})
	lns = append(lns, l)

	// flap 30
	start = primitives.Point{X: b.UpperLeft.X + side*4, Y: b.UpperLeft.Y + side}
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
