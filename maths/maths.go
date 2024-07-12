package maths

import "math"

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
