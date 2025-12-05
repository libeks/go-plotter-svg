package samplers

import (
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/libeks/go-plotter-svg/primitives"
)

var (
	perlinAlpha = 2.0
	perlinBeta  = 2.0
	perlinN     = int32(10)
	perlinSeed  = int64(109)
)

type PerlinNoise struct {
	noise  *perlin.Perlin
	offset primitives.Vector
	scale  float64
}

func NewPerlinNoise(scale float64, offset primitives.Vector) PerlinNoise {
	return PerlinNoise{
		noise:  perlin.NewPerlinRandSource(perlinAlpha, perlinBeta, perlinN, rand.NewSource(perlinSeed)),
		scale:  scale,
		offset: offset,
	}
}

// // returns a value from -1 to 1, based on Perlin Noise
// func (p PerlinNoise) GetFrameValue(x, y, t float64) float64 {
// 	val := p.noise.Noise3D(x+p.offsetX, y+p.offsetY, t)
// 	return val
// }

func (n PerlinNoise) GetValue(p primitives.Point) float64 {
	p = p.Add(n.offset)
	val := n.noise.Noise2D(n.scale*p.X, n.scale*p.Y)
	return val
	// if rand.Float64() < 0.001 {
	// 	fmt.Printf("val %.3f\n", val)
	// }
	// return val
}
