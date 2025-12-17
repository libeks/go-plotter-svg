package maths

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	PrecisionThreshold = 0.1
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

// given an interval [a,b], find in relative terms where t lies on that range
// It is assumed that a <= t <= b.
// if tPrime := ReverseInterpolatedTValue(a,b,t), then t == Interpolate(a,b tPrime)
func ReverseInterpolatedTValue(a, b, threshold float64) float64 {
	if a == b {
		// both endpoints are the same, default to 0.5 for consistency
		return 0.5
	}
	if a < b {
		width := b - a
		return (threshold - a) / width
	} else {
		width := a - b
		return (a - threshold) / width
	}
}

// given an interval [a,b], find in relative terms where t lies on that range
// It is assumed that a <= t <= b.
// if tPrime := ReverseInterpolatedTValue(a,b,t), then t == Interpolate(a,b tPrime)
func ReverseInterpolatedTValueFailure(a, b, threshold float64) float64 {
	if threshold < a || threshold > b {
		panic(fmt.Sprintf("Reverse Interpolation with incorrect threshold ( must have %.1f <= %.1f <= %.1f)", a, threshold, b))
	}
	return ReverseInterpolatedTValue(a, b, threshold)
}

// https://stackoverflow.com/a/59299881
// go doesn't do the expected thing for modding, since (-1%5) = -1, but we want to get 4 (to wrap around the index)
func Mod(a, b int) int {
	return (a%b + b) % b
}
