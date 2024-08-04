package collections

import (
	"fmt"
	"math"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

type RisingSun struct {
	BaselineY       float64
	LineSpacing     float64
	MinTurnRadius   float64
	NLines          int
	Sun             objects.Circle
	SunPadding      float64
	NLinesAroundSun int
}

func (s *RisingSun) Render(b box.Box) []lines.LineLike {
	lns := []lines.LineLike{}
	nStraightLines := s.NLines - s.NLinesAroundSun
	for i := range nStraightLines {
		ii := float64(i)
		baselineY := s.BaselineY - ii*s.LineSpacing
		start := primitives.Point{X: b.X, Y: baselineY}
		end := primitives.Point{X: b.XEnd, Y: baselineY}
		lns = append(lns, lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: end}))
	}
	leftCenter := primitives.Point{
		X: s.Sun.Center.X - s.Sun.Radius - s.SunPadding - s.MinTurnRadius - float64(s.NLinesAroundSun)*s.LineSpacing,
		Y: s.BaselineY - float64(s.NLines)*s.LineSpacing - s.MinTurnRadius,
	}
	rightCenter := primitives.Point{
		X: s.Sun.Center.X + s.Sun.Radius + s.SunPadding + s.MinTurnRadius + float64(s.NLinesAroundSun)*s.LineSpacing,
		Y: s.BaselineY - float64(s.NLines)*s.LineSpacing - s.MinTurnRadius,
	}
	for i := range s.NLinesAroundSun {
		ii := float64(i)
		baselineY := s.BaselineY - (ii+float64(nStraightLines))*s.LineSpacing
		start := primitives.Point{X: b.X, Y: baselineY}
		end := primitives.Point{X: b.XEnd, Y: baselineY}
		// verticalClimb := 0.0
		smallCircleRadius := s.MinTurnRadius + float64(s.NLinesAroundSun-i)*s.LineSpacing
		sunRadius := s.Sun.Radius + s.SunPadding + ii*s.LineSpacing

		xStart := s.Sun.Center.X - sunRadius
		xEnd := s.Sun.Center.X + sunRadius
		fmt.Printf("SmallCircleRadius: %.1f\n", smallCircleRadius)
		fmt.Printf("sunRadius: %.1f\n", sunRadius)
		fmt.Printf("start %.1f, end %.1f\n", xStart, xEnd)
		ln := lines.NewPath(start).AddPathChunk(lines.LineChunk{Start: start, End: primitives.Point{X: leftCenter.X, Y: baselineY}})
		ln = ln.AddPathChunk(lines.CircleArcChunk(leftCenter, smallCircleRadius, math.Pi/2, 2*math.Pi, false)) // circle up from base
		ln = ln.AddPathChunk(lines.LineChunk{
			Start: primitives.Point{X: xStart, Y: baselineY - smallCircleRadius},
			End:   primitives.Point{X: xStart, Y: s.Sun.Center.Y},
		}) // line upwards
		ln = ln.AddPathChunk(lines.CircleArcChunk(s.Sun.Center, sunRadius, math.Pi, 0, true)) // semicircle around sun
		ln = ln.AddPathChunk(lines.LineChunk{
			Start: primitives.Point{X: xEnd, Y: s.Sun.Center.Y},
			End:   primitives.Point{X: xEnd, Y: baselineY - smallCircleRadius},
		}) // line downwards
		ln = ln.AddPathChunk(lines.CircleArcChunk(rightCenter, smallCircleRadius, math.Pi, math.Pi/2, false)) // circle back to base
		ln = ln.AddPathChunk(lines.LineChunk{
			Start: primitives.Point{X: rightCenter.X, Y: s.Sun.Center.Y},
			End:   end,
		}) // line to end
		lns = append(lns, ln)

	}
	return lns
}
