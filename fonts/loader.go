package fonts

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/kintar/etxt/efixed"
	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/primitives"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	*truetype.Font
}

type Glyph struct {
	// contours [][]truetype.Point
	// bounds   fixed.Rectangle26_6
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
			Point: primitives.Point{
				X: b.X + r*efixed.ToFloat64(pt.X),
				Y: b.YEnd - r*efixed.ToFloat64(pt.Y),
			}, OnLine: isPointOnLine(pt),
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
				if onCurve {
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

func LoadFont(filename string) (*Font, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}
	return &Font{
		Font: font,
	}, nil
}

func (f *Font) LoadGlyph(r rune) (Glyph, error) {
	index := f.Index(r)
	glyph := truetype.GlyphBuf{}
	a, _ := efixed.FromFloat64(6000)
	err := glyph.Load(f.Font, a, index, font.HintingNone)
	if err != nil {
		return Glyph{}, err
	}
	fmt.Printf("Glyph %s\n", string(r))
	fmt.Printf("bounds %v\n", glyph.Bounds)
	fmt.Printf("Points %d, ends %d\n", len(glyph.Points), len(glyph.Ends))

	// nPoints := len(glyph.Points)
	for i, end := range glyph.Ends {
		fmt.Printf("Countour %d, %v\n", i, end)

		// for j :=
		// for _, point := range glypy.Points
	}
	// contours := partitionIntoContours(glyph.Ends, glyph.Points)

	return Glyph{
		glyph: glyph,
		// contours: contours,
		// bounds:   glyph.Bounds,
	}, nil
}

func partitionIntoContours(ends []int, points []truetype.Point) [][]truetype.Point {
	start := 0
	retlist := make([][]truetype.Point, len(ends))
	for i, end := range ends {
		retlist[i] = points[start:end]
		start = end + 1
	}
	return retlist
}
