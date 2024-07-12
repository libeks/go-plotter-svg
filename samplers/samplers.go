package samplers

import (
	"math/rand"

	"github.com/libeks/go-plotter-svg/primitives"
)

type DataSource interface {
	GetValue(p primitives.Point) float64
}

type ConstantDataSource struct {
	val float64
}

func (s ConstantDataSource) GetValue(p primitives.Point) float64 {
	return s.val
}

type RandomDataSource struct {
}

func (s RandomDataSource) GetValue(p primitives.Point) float64 {
	return rand.Float64()
}
