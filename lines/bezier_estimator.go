package lines

import "fmt"

func getLengthEstimate(o LengthEstimator, segments int) float64 {
	distance := 0.0
	start := o.At(0)
	for i := range segments {
		t := float64(i+1) / float64(segments)
		end := o.At(t)
		dist := end.Subtract(start).Len()
		distance += dist
		start = end
	}
	return distance
}

func estimateLength(o LengthEstimator, acc float64) float64 {
	if o.BBox().IsEmpty() {
		fmt.Printf("chunk %s has empty bbox %s\n", o, o.BBox())
		return 0
	}
	oldLen := 0.0
	n := 1
	// double number of segments until length stabilizes below acc, or if there are more than 100 segments
	for {
		// evaluate on an odd number of segments
		// this is due to no change when going from 1->2
		nn := n + 1

		newLen := getLengthEstimate(o, nn)
		if (newLen - oldLen) < acc {
			return newLen
		}
		if n > 100 {
			return newLen
		}
		oldLen = newLen
		n = n * 2
	}
	return 0
}
