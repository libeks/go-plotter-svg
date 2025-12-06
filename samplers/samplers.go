package samplers

import (
	"math"
	"math/rand"

	"github.com/libeks/go-plotter-svg/primitives"
)

type DataSource interface {
	GetValue(p primitives.Point) float64
}

func Constant(val float64) constant {
	return constant{
		val: val,
	}
}

type constant struct {
	val float64
}

func (s constant) GetValue(p primitives.Point) float64 {
	return s.val
}

type RandomDataSource struct {
}

func (s RandomDataSource) GetValue(p primitives.Point) float64 {
	return rand.Float64()
}

type HighCenterRelativeDataSource struct {
	Scale float64
}

// assumes that point will be from a point in the bounding box -1..1
func (s HighCenterRelativeDataSource) GetValue(p primitives.Point) float64 {
	// distance to center
	dist := p.Subtract(primitives.Origin).Len()
	return 1 / (dist*s.Scale + 1) // should be in range (0,1]
}

type HighInCircleRelativeDataSource struct {
	Radius float64
}

// assumes that point will be from a point in the bounding box -1..1
func (s HighInCircleRelativeDataSource) GetValue(p primitives.Point) float64 {
	// distance to center
	dist := p.Subtract(primitives.Origin).Len()
	if dist < s.Radius {
		return 1.0
	}
	return 0.0
}

type InvertSampler struct {
	DataSource
}

func (s InvertSampler) GetValue(p primitives.Point) float64 {
	// distance to center
	return 1.0 - s.DataSource.GetValue(p)
}

type InsideCircleRelativeDataSource struct {
	Radius  float64
	Inside  float64
	Outside float64
}

// assumes that point will be from a point in the bounding box -1..1
func (s InsideCircleRelativeDataSource) GetValue(p primitives.Point) float64 {
	// distance to center
	dist := p.Subtract(primitives.Origin).Len()
	if dist < s.Radius {
		return s.Inside
	}
	return s.Outside
}

type InsideCircleSubDataSource struct {
	Radius  float64
	Inside  DataSource
	Outside DataSource
}

// assumes that point will be from a point in the bounding box -1..1
func (s InsideCircleSubDataSource) GetValue(p primitives.Point) float64 {
	// distance to center
	dist := p.Subtract(primitives.Origin).Len()
	if dist < s.Radius {
		return s.Inside.GetValue(p)
	}
	return s.Outside.GetValue(p)
}

type RandomChooser struct {
	Values []float64
}

func (s RandomChooser) GetValue(p primitives.Point) float64 {
	return s.Values[rand.Intn(len(s.Values))]
}

func PointDistance(p primitives.Point) pointDistance {
	return pointDistance{
		point: p,
	}
}

type pointDistance struct {
	point primitives.Point
}

// assumes that point will be from a point in the bounding box -1..1
func (c pointDistance) GetValue(p primitives.Point) float64 {
	return p.Subtract(c.point).Len()
}

type BooleanSwitcher struct {
	BooleanSource Booleaner
	WhenTrue      DataSource
	WhenFalse     DataSource
}

func (s BooleanSwitcher) GetValue(p primitives.Point) float64 {
	if s.BooleanSource.GetBool(p) {
		return s.WhenTrue.GetValue(p)
	}
	return s.WhenFalse.GetValue(p)
}

type AngleFromCenter struct {
	Center primitives.Point
}

func (c AngleFromCenter) GetValue(p primitives.Point) float64 {
	return c.Center.Subtract(p).Atan()
}

type TurnAngleByRightAngle struct {
	Center primitives.Point
}

func (c TurnAngleByRightAngle) GetValue(p primitives.Point) float64 {
	return c.Center.Subtract(p).Perp().Atan()
}

func Add(sources ...DataSource) DataSource {
	return addSlice{
		Samplers: sources,
	}
}

type addSlice struct {
	Samplers []DataSource
}

func (s addSlice) GetValue(p primitives.Point) float64 {
	total := 0.0
	for _, sampler := range s.Samplers {
		total += sampler.GetValue(p)
	}
	return total
}

// type Min struct {
// 	SamplerA DataSource
// 	SamplerB DataSource
// }

// func (s Min) GetValue(p primitives.Point) float64 {
// 	return min(s.SamplerA.GetValue(p), s.SamplerB.GetValue(p))
// }

func Min(samplers ...DataSource) minSlice {
	return minSlice{
		Samplers: samplers,
	}
}

type minSlice struct {
	Samplers []DataSource
}

func (s minSlice) GetValue(p primitives.Point) float64 {
	min := math.MaxFloat64
	for _, sampler := range s.Samplers {
		val := sampler.GetValue(p)
		if val < min {
			min = val
		}
	}
	return min
}

func Lambda(f func(p primitives.Point) float64) lambdaFn {
	return lambdaFn{function: f}
}

type lambdaFn struct {
	function func(p primitives.Point) float64
}

func (c lambdaFn) GetValue(p primitives.Point) float64 {
	return c.function(p)
}

func ScalarMultiple(sampler DataSource, scalar float64) scalarMultiple {
	return scalarMultiple{
		sampler: sampler,
		scalar:  scalar,
	}
}

type scalarMultiple struct {
	sampler DataSource
	scalar  float64
}

func (s scalarMultiple) GetValue(p primitives.Point) float64 {
	return s.sampler.GetValue(p) * s.scalar
}
