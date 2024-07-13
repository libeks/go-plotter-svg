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
