package collections

import (
	"fmt"
	"math"
	"math/rand"
	"slices"

	"github.com/libeks/go-plotter-svg/box"
	"github.com/libeks/go-plotter-svg/lines"
	"github.com/libeks/go-plotter-svg/maths"
	"github.com/libeks/go-plotter-svg/objects"
	"github.com/libeks/go-plotter-svg/primitives"
)

func LimitLinesToShape(ls []lines.Line, shape objects.Object) []lines.LineSegment {
	segments := []lines.LineSegment{}
	for _, line := range ls {
		segs := ClipLineToObject(line, shape)
		segments = append(segments, segs...)
	}
	return segments
}

func CircularLineField(n int, center primitives.Point) []lines.Line {
	linelikes := []lines.Line{}
	for i := range n {
		ii := float64(i)
		angle := math.Pi * ii / (float64(n))
		linelikes = append(linelikes, lines.Line{
			P: center,
			V: primitives.Vector{X: 100 * math.Sin(angle), Y: 100 * math.Cos(angle)},
		})
	}
	return linelikes
}

func LimitCirclesToShape(circles []objects.Circle, shape objects.Object) []lines.LineLike {
	segments := []lines.LineLike{}
	for _, circle := range circles {
		segs := ClipCircleToObject(circle, shape)
		segments = append(segments, segs...)
	}
	return segments
}

func ConcentricCircles(b box.Box, center primitives.Point, spacing float64) []objects.Circle {
	// find out the maximum radius
	maxDist := 0.0
	for _, p := range b.Corners() {
		v := p.Subtract(center)
		dist := v.Len()
		if dist > maxDist {
			maxDist = dist
		}
	}
	nCircles := maxDist / spacing
	circles := []objects.Circle{}
	for i := range int(nCircles) + 1 {
		ii := float64(i)
		circles = append(circles, objects.Circle{Center: center, Radius: spacing * ii})
	}
	return circles
}

// find out min and max line index. The 0th line goes through the origin (0,0)
func getLineIndexRange(b box.Box, perpVect primitives.Vector) (float64, float64) {
	iSlice := []float64{}
	wSq := perpVect.Dot(perpVect)
	for _, point := range []primitives.Point{b.NWCorner(), b.NECorner(), b.SWCorner(), b.SECorner()} {
		pVect := point.Subtract(primitives.Point{X: 0, Y: 0})
		i := pVect.Dot(perpVect) / wSq
		iSlice = append(iSlice, i)
	}
	minI := slices.Min(iSlice)
	maxI := slices.Max(iSlice)
	return minI, maxI
}

// LinearLineField returns a set of parallel lines, all oriented in the direction of angle relative to 0x axis,
// only returns the lines that would fall inside the box
func LinearLineField(b box.Box, angle float64, spacing float64) []lines.Line {
	// find out min and max line index. The 0th line goes through the origin (0,0)
	v := primitives.Vector{X: math.Cos(angle), Y: math.Sin(angle)}.Mult(spacing)
	w := v.Perp()
	minI, maxI := getLineIndexRange(b, w)
	ls := []lines.Line{}
	for i := int(minI) - 1; i <= int(maxI)+1; i++ {
		ls = append(ls, lines.Line{P: w.Mult(float64(i)).Point(), V: v})
	}
	return ls
}

// LinearDensityLineField returns a set of parallel lines, all oriented in the direction of angle relative to 0x axis,
// only returns the lines that would fall inside the box
// densityFn takes input in the range [0;1] and outputs positive values denoting the spacing to respect
// at every increment
func LinearDensityLineField(b box.Box, angle float64, densityFn func(float64) float64) []lines.Line {
	// find out min and max line index. The 0th line goes through the origin (0,0)
	v := primitives.Vector{X: math.Cos(angle), Y: math.Sin(angle)}
	w := v.Perp()
	minI, maxI := getLineIndexRange(b, w)
	ls := []lines.Line{}
	i := minI
	for i < maxI+1 {
		ls = append(ls, lines.Line{P: w.Mult(float64(i)).Point(), V: v})
		densityVal := densityFn((i - minI) / (maxI - minI))
		fmt.Printf("Densityval at %.1f %.1f \n", (i-minI)/(maxI-minI), densityVal)
		i += densityVal
	}
	return ls
}

type StrokeStrip struct {
	box     box.Box
	padding float64
	Direction
}

func (h StrokeStrip) String() string {
	return fmt.Sprintf("StrokeStrip %s padding %.1f", h.box, h.padding)
}

func (h StrokeStrip) Lines() []lines.LineLike {
	var nLines int
	if h.Direction.CardinalDirection == Horizontal {
		nLines = int(h.box.Height()/h.padding) + 1
	} else {
		nLines = int((h.box.Width())/h.padding) + 1
	}
	ls := make([]lines.LineLike, nLines)

	for i := range nLines {
		j := i
		if h.Direction.OrderDirection == AwayToHome {
			j = nLines - i - 1
		}
		reverse := (h.Direction.StrokeDirection == AwayToHome)
		if h.Direction.Connection == AlternatingDirection && (i%2 == 1) {
			reverse = !reverse
		}
		var line lines.LineSegment
		if h.Direction.CardinalDirection == Horizontal {
			vect := primitives.Vector{X: 0, Y: float64(j) * h.padding}
			line = lines.LineSegment{
				P1: h.box.NWCorner().Add(vect),
				P2: h.box.SWCorner().Add(vect),
			}
		} else {
			vect := primitives.Vector{X: float64(j) * h.padding, Y: 0}
			line = lines.LineSegment{
				P1: h.box.NWCorner().Add(vect),
				P2: h.box.NECorner().Add(vect),
			}
		}
		if reverse {
			ls[i] = line.Reverse()
		} else {
			ls[i] = line
		}
	}
	return ls
}

func ClipLineToObject(line lines.Line, obj objects.Object) []lines.LineSegment {
	ts := obj.IntersectTs(line)
	if len(ts) == 0 {
		return nil
	}
	if len(ts) == 1 {
		panic(fmt.Errorf("not sure what to do with only one intersection %v", ts))
	}
	slices.Sort(ts)
	segments := []lines.LineSegment{}
	for i, t1 := range ts {
		if i == len(ts)-1 {
			break
		}
		t2 := ts[(i + 1)]
		midT := maths.Average(t1, t2)
		midPoint := line.At(midT)
		if obj.Inside(midPoint) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			seg := lines.LineSegment{P1: p1, P2: p2}
			segments = append(segments, seg)
		}
	}
	return segments
}

func ClipLineSegmentToObject(ls lines.LineSegment, obj objects.Object) []lines.LineSegment {
	line := ls.Line()
	ts := append(obj.IntersectTs(line), 0.0, 1.0)
	ts = filterToRange(ts, 0.0, 1.0)
	slices.Sort(ts)
	segments := []lines.LineSegment{}
	for i, t1 := range ts {
		t2 := ts[(i+1)%len(ts)]
		midT := maths.Average(t1, t2)
		midPoint := line.At(midT)
		if obj.Inside(midPoint) {
			p1 := line.At(t1)
			p2 := line.At(t2)
			segments = append(segments, lines.LineSegment{
				P1: p1, P2: p2,
			})
		}
	}
	return segments
}

func ClipCircleToObject(c objects.Circle, obj objects.Object) []lines.LineLike {
	ts := obj.IntersectCircleTs(c)
	if len(ts) == 0 {
		return []lines.LineLike{c}
	}
	slices.Sort(ts)
	segments := []lines.LineLike{}
	for i, t1 := range ts {
		t2 := ts[(i+1)%len(ts)]
		midT := maths.Average(t1, t2)
		midPoint := c.At(midT)
		if obj.Inside(midPoint) {
			segments = append(segments, CircleArc(c, t1, t2))
		}
	}
	return segments
}

func filterToRange(vals []float64, min, max float64) []float64 {
	rets := []float64{}
	for _, val := range vals {
		if val <= max && val >= min {
			rets = append(rets, val)
		}
	}
	return rets
}

func ConcentricCirclesInCircle(circle objects.Circle, spacing float64) []lines.LineLike {
	lns := []lines.LineLike{}
	nLines := circle.Radius / spacing
	for i := range int(nLines) {
		radius := (float64(i) + 0.5) * spacing
		lns = append(lns, lines.FullCircle(circle.Center, radius))
	}
	return lns
}

func CircleArc(circle objects.Circle, t1 float64, t2 float64) lines.Path {
	p1 := circle.At(t1)
	path := lines.NewPath(p1).AddPathChunk(lines.CircleArcChunk(circle.Center, circle.Radius, t1, t2, false))
	return path
}

type SineDensity struct {
	Min    float64
	Max    float64
	Offset float64
	Cycles float64
}

func (d SineDensity) Density(a float64) float64 {
	theta := d.Cycles * (a + d.Offset) * math.Pi
	dRange := d.Max - d.Min
	return d.Min + dRange*(math.Sin(theta)+1)/2
}

func RandomlyAllocateSegments(segments [][]lines.LineLike, threshold float64) ([]lines.LineLike, []lines.LineLike) {
	layer1 := []lines.LineLike{}
	layer2 := []lines.LineLike{}
	for _, segs := range segments {
		if rand.Float64() > threshold {
			layer1 = append(layer1, segs...)
		} else {
			layer2 = append(layer2, segs...)
		}
	}
	return layer1, layer2
}
