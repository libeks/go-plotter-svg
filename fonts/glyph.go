package fonts

import (
	"fmt"

	"github.com/golang/freetype/truetype"
	"github.com/kintar/etxt/efixed"
	"golang.org/x/image/math/fixed"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
)

type Glyph struct {
	glyph truetype.GlyphBuf
}

func (g Glyph) Contours() [][]truetype.Point {
	start := 0
	retlist := make([][]truetype.Point, len(g.glyph.Ends))
	for i, end := range g.glyph.Ends {
		retlist[i] = g.glyph.Points[start:end]
		start = end
	}
	return retlist
}
func getWidthHeight(b fixed.Rectangle26_6) (float64, float64) {
	glyphWidth := b.Max.X - b.Min.X
	glyphHeight := b.Max.Y - b.Min.Y
	return efixed.ToFloat64(glyphWidth), efixed.ToFloat64(glyphHeight)
}

type ControlPoint struct {
	primitives.Point
	OnLine bool
}

func (g Glyph) GetControlPoints(b box.Box) []ControlPoint {
	w, h := getWidthHeight(g.glyph.Bounds)
	wRatio, hRatio := b.Width()/w, b.Height()/h
	r := min(wRatio, hRatio)
	pts := []ControlPoint{}
	for _, pt := range g.glyph.Points {
		pts = append(pts, ControlPoint{
			Point:  convertPoint(b, pt, r),
			OnLine: isPointOnLine(pt),
		})
	}
	return pts
}

// positions the point relative to the box and the scaling factor r
func convertPoint(b box.Box, pt truetype.Point, r float64) primitives.Point {
	return primitives.Point{
		X: b.X + r*efixed.ToFloat64(pt.X),
		Y: b.YEnd - r*efixed.ToFloat64(pt.Y),
	}
}

func isPointOnLine(pt truetype.Point) bool {
	return (pt.Flags & 1) > 0
}

// cycle through the points of the contour until the first point is on the line, for easier rendering
func optimizeContour(pts []truetype.Point) []truetype.Point {
	idx := 0
	for i, pt := range pts {
		if isPointOnLine(pt) {
			idx = i
			break
		}
	}
	return append(pts[idx:], pts[:idx]...)
}

func (g Glyph) GetCurves(b box.Box) []lines.LineLike {
	w, h := getWidthHeight(g.glyph.Bounds)
	wRatio, hRatio := b.Width()/w, b.Height()/h
	r := min(wRatio, hRatio)
	lns := []lines.LineLike{}
	for _, contour := range g.Contours() {
		contour = optimizeContour(contour)
		if len(contour) == 0 {
			fmt.Printf("Encountered empty contour!")
			continue
		}
		// is point on line?
		if !isPointOnLine(contour[0]) {
			panic("First point on contour is not on line")
		}
		start := convertPoint(b, contour[0], r)
		midpoint := primitives.Point{}
		l := lines.NewPath(start)
		onCurve := false
		idx := 1
		for idx < len(contour) {
			pt := contour[idx]
			cp := convertPoint(b, pt, r)
			if isPointOnLine(pt) {
				if !onCurve {
					l = l.AddPathChunk(lines.LineChunk{Start: start, End: cp})
				} else {
					l = l.AddPathChunk(lines.QuadraticBezierChunk{Start: start, P1: midpoint, End: cp})
				}
				onCurve = false
				start = cp
			} else {
				// current point is a control point
				if onCurve {
					// previous point was also a control point, so we need to chain points correctly
					// get midpoint between successive bezier control points
					c := primitives.Midpoint(midpoint, cp)
					fmt.Printf("Midpoint between %v and %v is %v\n", midpoint, cp, c)

					l = l.AddPathChunk(lines.QuadraticBezierChunk{Start: start, P1: midpoint, End: c})
					start = c
					midpoint = cp
				} else {
					midpoint = convertPoint(b, pt, r)
				}
				onCurve = true
			}
			idx += 1
			if idx == len(contour) {
				idx = 0
			}
			if idx == 1 {
				break
			}
		}
		lns = append(lns, l)
	}
	return lns
}

// Explanation of why fixed-point 26.6 is used for fonts: https://github.com/Kintar/etxt/blob/main/docs/fixed-26-6.md
