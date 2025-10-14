package lines

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/libeks/go-plotter-svg/primitives"
)

func TestLineIntersection(t *testing.T) {
	type testCase struct {
		Name      string
		LineA     Line
		LineB     Line
		ExpectNil bool
		Expect    primitives.Point
	}
	tests := []testCase{
		{
			Name:   "origin",
			LineA:  Line{P: primitives.Origin, V: primitives.Vector{X: 1, Y: 0}},
			LineB:  Line{P: primitives.Origin, V: primitives.Vector{X: 0, Y: 1}},
			Expect: primitives.Origin,
		},
		{
			Name:   "origin*100",
			LineA:  Line{P: primitives.Origin, V: primitives.Vector{X: 100, Y: 0}},
			LineB:  Line{P: primitives.Origin, V: primitives.Vector{X: 0, Y: 100}},
			Expect: primitives.Origin,
		},
		{
			Name:   "1_1",
			LineA:  Line{P: primitives.Point{X: 1, Y: 1}, V: primitives.Vector{X: 100, Y: 0}},
			LineB:  Line{P: primitives.Point{X: 1, Y: 1}, V: primitives.Vector{X: 0, Y: 100}},
			Expect: primitives.Point{X: 1, Y: 1},
		},
		{
			Name:   "1_1",
			LineA:  Line{P: primitives.Point{X: 100, Y: 100}, V: primitives.Vector{X: 100, Y: 0}},
			LineB:  Line{P: primitives.Point{X: 100, Y: 100}, V: primitives.Vector{X: 0, Y: 100}},
			Expect: primitives.Point{X: 100, Y: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := tt.LineA.Intersect(tt.LineB)
			fmt.Printf("Got %v\n", got)
			if (got == nil) != tt.ExpectNil {
				t.Fatalf("Expected %v, got %v", tt.ExpectNil, got)
			}
			if tt.ExpectNil {
				return
			}
			if diff := cmp.Diff(tt.Expect, *got); diff != "" {
				fmt.Printf("got (%.6f, %.6f), want (%.6f, %.6f)\n", got.X, got.Y, tt.Expect.X, tt.Expect.Y)
				t.Fatalf("Unexpected diff %v", diff)
			}
		})
	}
}
