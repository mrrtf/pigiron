package geo

import "math"

// EqualFloat checks if two float64 are equal within 1E-4 (on difference)
func EqualFloat(a, b float64) bool {
	return (math.Abs(a-b) < 1E-4) // 1E-4 cm = 1 micron
}

// IsStrictlyBelowFloat returns true if a < b
func IsStrictlyBelowFloat(a, b float64) bool {
	return (a < b) && !EqualFloat(a, b)
}

// IsInRangeFloat64 returns true if a<=x<=b
func IsInRangeFloat64(x, a, b float64) bool {
	return x >= a && x <= b
}

// EqualFloat64Slice checks if two float64 slices are equal
// (same elements in exact same order)
func EqualFloat64Slice(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if !EqualFloat(v, b[i]) {
			return false
		}
	}
	return true
}

func removeDuplicates(in []float64) []float64 {
	out := make([]float64, 0, len(in))
	for i := 0; i < len(in); i++ {
		already := false
		for j := 0; j < len(out); j++ {
			if EqualFloat(in[i], out[j]) {
				already = true
				break
			}
		}
		if !already {
			out = append(out, in[i])
		}
	}
	return out
}
