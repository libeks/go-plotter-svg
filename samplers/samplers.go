package samplers

import (
	"math/rand"

	"github.com/libeks/go-plotter-svg/primitives"
)

type DataSource interface {
	GetValue(p primitives.Point) float64
}

type ConstantDataSource struct {
	Val float64
}

func (s ConstantDataSource) GetValue(p primitives.Point) float64 {
	return s.Val
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
