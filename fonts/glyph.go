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

// Glyph is the abstract glyph form for a specific pixel height
type Glyph struct {
	Rune    rune
	glyph   truetype.GlyphBuf
	bounds  fixed.Rectangle26_6
	advance fixed.Int26_6
}

// Char is a character rendered at a specific height to the page
type Char struct {
	Rune         rune
	Curves       []lines.LineLike
	Points       []ControlPoint
	BoundingBox  box.Box
	AdvanceWidth float64
}

func (c Char) Translate(v primitives.Vector) Char {
	newCurves := make([]lines.LineLike, len(c.Curves))
	for i, c := range c.Curves {
		newCurves[i] = c.Translate(v)
	}
	newPoints := make([]ControlPoint, len(c.Points))
	for i, c := range c.Points {
		newPoints[i] = c.Translate(v)
	}
	return Char{
		Rune:         c.Rune,
		Curves:       newCurves,
		Points:       newPoints,
		BoundingBox:  c.BoundingBox.Translate(v),
		AdvanceWidth: c.AdvanceWidth,
	}
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

func (c ControlPoint) Translate(v primitives.Vector) ControlPoint {
	c.Point = c.Point.Add(v)
	return c
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
	res := primitives.Point{
		X: b.UpperLeft.X + r*efixed.ToFloat64(pt.X),
		Y: b.LowerRight.Y - r*efixed.ToFloat64(pt.Y),
	}
	return res
}

// positions the point relative to scaling factor r, in abstract space
func convertStaticPoint(pt truetype.Point, r float64) primitives.Point {
	res := primitives.Point{
		X: r * efixed.ToFloat64(pt.X),
		Y: -r * efixed.ToFloat64(pt.Y),
	}
	return res
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

func (g Glyph) GetHeightCurves(h float64) Char {
	r := h / fontHeight
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
		start := convertStaticPoint(contour[0], r)
		midpoint := primitives.Point{}
		l := lines.NewPath(start)
		onCurve := false
		idx := 1
		for idx < len(contour) {
			pt := contour[idx]
			cp := convertStaticPoint(pt, r)
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

					l = l.AddPathChunk(lines.QuadraticBezierChunk{Start: start, P1: midpoint, End: c})
					start = c
					midpoint = cp
				} else {
					midpoint = convertStaticPoint(pt, r)
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
	minPoint := truetype.Point{
		X: g.bounds.Min.X,
		Y: g.bounds.Min.Y,
	}
	maxPoint := truetype.Point{
		X: g.bounds.Max.X,
		Y: g.bounds.Max.Y,
	}
	minP := convertStaticPoint(minPoint, r)
	maxP := convertStaticPoint(maxPoint, r)
	bbox := box.Box{
		BBox: primitives.BBox{
			UpperLeft: primitives.Point{
				X: minP.X,
				Y: maxP.Y,
			},
			LowerRight: primitives.Point{
				X: maxP.X,
				Y: minP.Y,
			},
		},
	}
	pts := []ControlPoint{}
	for _, pt := range g.glyph.Points {
		pts = append(pts, ControlPoint{
			Point:  convertStaticPoint(pt, r),
			OnLine: isPointOnLine(pt),
		})
	}
	return Char{
		Rune:         g.Rune,
		Curves:       lns,
		Points:       pts,
		BoundingBox:  bbox,
		AdvanceWidth: r * efixed.ToFloat64(g.advance),
	}
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
