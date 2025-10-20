package lines

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/libeks/go-plotter-svg/primitives"
)

func TestLineIntersection(t *testing.T) {
	type testCase struct {
		name      string
		lineA     Line
		lineB     Line
		expectNil bool
		expect    primitives.Point
	}
	tests := []testCase{
		{
			name:   "origin",
			lineA:  Line{P: primitives.Origin, V: primitives.Vector{X: 1, Y: 0}},
			lineB:  Line{P: primitives.Origin, V: primitives.Vector{X: 0, Y: 1}},
			expect: primitives.Origin,
		},
		{
			name:   "origin*100",
			lineA:  Line{P: primitives.Origin, V: primitives.Vector{X: 100, Y: 0}},
			lineB:  Line{P: primitives.Origin, V: primitives.Vector{X: 0, Y: 100}},
			expect: primitives.Origin,
		},
		{
			name:   "1_1",
			lineA:  Line{P: primitives.Point{X: 1, Y: 1}, V: primitives.Vector{X: 100, Y: 0}},
			lineB:  Line{P: primitives.Point{X: 1, Y: 1}, V: primitives.Vector{X: 0, Y: 100}},
			expect: primitives.Point{X: 1, Y: 1},
		},
		{
			name:   "1_1",
			lineA:  Line{P: primitives.Point{X: 100, Y: 100}, V: primitives.Vector{X: 100, Y: 0}},
			lineB:  Line{P: primitives.Point{X: 100, Y: 100}, V: primitives.Vector{X: 0, Y: 100}},
			expect: primitives.Point{X: 100, Y: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lineA.Intersect(tt.lineB)
			fmt.Printf("Got %v\n", got)
			if (got == nil) != tt.expectNil {
				t.Fatalf("Expected %v, got %v", tt.expectNil, got)
			}
			if tt.expectNil {
				return
			}
			if diff := cmp.Diff(tt.expect, *got); diff != "" {
				fmt.Printf("got (%.6f, %.6f), want (%.6f, %.6f)\n", got.X, got.Y, tt.expect.X, tt.expect.Y)
				t.Fatalf("Unexpected diff %v", diff)
			}
		})
	}
}
