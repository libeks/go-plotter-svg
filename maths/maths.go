package maths

import (
	"math"
	"math/rand"
)

func Average(a, b float64) float64 {
	return (a + b) / 2
}

func AngleDifference(a2, a1 float64) float64 {
	diff := a2 - a1
	if diff < -math.Pi {
		diff = diff + math.Pi
	}
	if diff > math.Pi {
		diff = diff - math.Pi
	}
	return diff
}

func Quadratic(a, b, c float64) []float64 {
	discriminant := b*b - 4*a*c
	if discriminant < 0.0 {
		return nil
	}
	if discriminant == 0 {
		return []float64{
			-b / (2 * a),
		}
	}
	d := math.Sqrt(discriminant)
	return []float64{
		(-b - d) / (2 * a),
		(-b + d) / (2 * a),
	}
}

func SumFloats(l []float64) float64 {
	total := 0.0
	for _, v := range l {
		total += v
	}
	return total
}

func RandRangeMinusPlusOne() float64 {
	return 2 * (rand.Float64() - 0.5)
}

func RandInRange(min, max float64) float64 {
	return (max-min)*rand.Float64() + min
}

// interpolate between a,b, with t in range [0,1]/ t=0 => a, t=1 => b
func Interpolate(a, b, t float64) float64 {
	return (b-a)*t + a
}
