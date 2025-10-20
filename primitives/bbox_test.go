package primitives

import (
	// "fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBBoxIntersect(t *testing.T) {
	tests := []struct {
		name            string
		boxA            BBox
		boxB            BBox
		expectIsOverlap bool
		expectOverlap   BBox
	}{
		{
			name:            "no_overlap",
			boxA:            BBox{UpperLeft: Origin, LowerRight: Origin.Add(Vector{X: 100, Y: 100})},
			boxB:            BBox{UpperLeft: Origin.Add(Vector{X: 1000, Y: 1000}), LowerRight: Origin.Add(Vector{X: 1100, Y: 1100})},
			expectIsOverlap: false,
			expectOverlap:   BBox{},
		},
		{
			name:            "inside_the_other",
			boxA:            BBox{UpperLeft: Origin, LowerRight: Origin.Add(Vector{X: 1000, Y: 1000})},
			boxB:            BBox{UpperLeft: Origin.Add(Vector{X: 400, Y: 400}), LowerRight: Origin.Add(Vector{X: 600, Y: 600})},
			expectIsOverlap: true,
			expectOverlap:   BBox{UpperLeft: Origin.Add(Vector{X: 400, Y: 400}), LowerRight: Origin.Add(Vector{X: 600, Y: 600})},
		},
		{
			name:            "just_the_edge",
			boxA:            BBox{UpperLeft: Origin, LowerRight: Origin.Add(Vector{X: 1000, Y: 1000})},
			boxB:            BBox{UpperLeft: Origin.Add(Vector{X: 900, Y: 400}), LowerRight: Origin.Add(Vector{X: 2000, Y: 600})},
			expectIsOverlap: true,
			expectOverlap:   BBox{UpperLeft: Origin.Add(Vector{X: 900, Y: 400}), LowerRight: Origin.Add(Vector{X: 1000, Y: 600})},
		},
		{
			name:            "just_the_corner",
			boxA:            BBox{UpperLeft: Origin, LowerRight: Origin.Add(Vector{X: 1000, Y: 1000})},
			boxB:            BBox{UpperLeft: Origin.Add(Vector{X: 900, Y: 900}), LowerRight: Origin.Add(Vector{X: 2000, Y: 2000})},
			expectIsOverlap: true,
			expectOverlap:   BBox{UpperLeft: Origin.Add(Vector{X: 900, Y: 900}), LowerRight: Origin.Add(Vector{X: 1000, Y: 1000})},
		},
		{
			name:            "opposite_corner",
			boxA:            BBox{UpperLeft: Origin.Add(Vector{X: 0, Y: 900}), LowerRight: Origin.Add(Vector{X: 1000, Y: 1900})},
			boxB:            BBox{UpperLeft: Origin.Add(Vector{X: 900, Y: 0}), LowerRight: Origin.Add(Vector{X: 1900, Y: 1000})},
			expectIsOverlap: true,
			expectOverlap:   BBox{UpperLeft: Origin.Add(Vector{X: 900, Y: 900}), LowerRight: Origin.Add(Vector{X: 1000, Y: 1000})},
		},
		{
			name:            "adjoining_zero_pixel_overlap",
			boxA:            BBox{UpperLeft: Origin, LowerRight: Origin.Add(Vector{X: 1000, Y: 1000})},
			boxB:            BBox{UpperLeft: Origin.Add(Vector{X: 1000, Y: 0}), LowerRight: Origin.Add(Vector{X: 2000, Y: 1000})},
			expectIsOverlap: true,
			expectOverlap:   BBox{UpperLeft: Origin.Add(Vector{X: 1000, Y: 0}), LowerRight: Origin.Add(Vector{X: 1000, Y: 1000})},
		},
		{
			name:            "adjoining_no_overlap",
			boxA:            BBox{UpperLeft: Origin, LowerRight: Origin.Add(Vector{X: 1000, Y: 1000})},
			boxB:            BBox{UpperLeft: Origin.Add(Vector{X: 1001, Y: 0}), LowerRight: Origin.Add(Vector{X: 2000, Y: 1000})},
			expectIsOverlap: false,
			expectOverlap:   BBox{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overlap, ok := tt.boxA.Intersect(tt.boxB)
			// fmt.Printf("Got %v\n", ok)
			if diff := cmp.Diff(tt.expectIsOverlap, ok); diff != "" {
				t.Fatalf("Unexpected diff %v", diff)
			}
			if diff := cmp.Diff(tt.expectOverlap, overlap); diff != "" {
				// fmt.Printf("got (%.6f, %.6f), want (%.6f, %.6f)\n", got.X, got.Y, tt.Expect.X, tt.Expect.Y)
				t.Fatalf("Unexpected diff %v", diff)
			}
		})
	}
}
