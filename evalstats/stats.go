// Package evalstats calculates metrics for evaluating physical models.
package evalstats

import (
	"math"

	"github.com/gonum/floats"
)

// MFB calculates the mean fractional bias of b against a. It assumes a and
// b are the same length.
func MFB(a, b []float64) float64 {
	r := 0.
	for i, v1 := range a {
		v2 := b[i]
		r += 2 * (v2 - v1) / (v1 + v2)
	}
	return r / float64(len(a))
}

// MFBWeighted calculates the mean fractional bias of b against a, weighted by w.
// It assumes a, b, and w are the same length.
func MFBWeighted(a, b, w []float64) float64 {
	r := 0.
	for i, v1 := range a {
		v2 := b[i]
		r += 2 * (v2 - v1) / (v1 + v2) * w[i]
	}
	return r / floats.Sum(w)
}

// MFE calculates the mean fractional error of b against a. It assumes a and
// b are the same length.
func MFE(a, b []float64) float64 {
	r := 0.
	for i, v1 := range a {
		v2 := b[i]
		r += 2 * math.Abs(v2-v1) / math.Abs(v1+v2)
	}
	return r / float64(len(a))
}

// MFEWeighted calculates the mean fractional error of b against a, weighted by
// w. It assumes a, b, and w are the same length.
func MFEWeighted(a, b, w []float64) float64 {
	r := 0.
	for i, v1 := range a {
		v2 := b[i]
		r += 2 * math.Abs(v2-v1) / math.Abs(v1+v2) * w[i]
	}
	return r / floats.Sum(w)
}

// MB calculates the mean bias of b against a. It assumes a and
// b are the same length.
func MB(a, b []float64) float64 {
	r := 0.
	for i, v1 := range a {
		v2 := b[i]
		r += (v2 - v1)
	}
	return r / float64(len(a))
}

// MBWeighted calculates the mean bias of b against a, weighted by w.
// It assumes a, b, and w are the same length.
func MBWeighted(a, b, w []float64) float64 {
	r := 0.
	for i, v1 := range a {
		v2 := b[i]
		r += (v2 - v1) * w[i]
	}
	return r / floats.Sum(w)
}

// ME calculates the mean error of b against a. It assumes a and
// b are the same length.
func ME(a, b []float64) float64 {
	r := 0.
	for i, v1 := range a {
		v2 := b[i]
		r += math.Abs(v2 - v1)
	}
	return r / float64(len(a))
}

// MEWeighted calculates the mean error of b against a, weighted by w.
// It assumes a, b, and w are the same length.
func MEWeighted(a, b, w []float64) float64 {
	r := 0.
	for i, v1 := range a {
		v2 := b[i]
		r += math.Abs(v2-v1) * w[i]
	}
	return r / floats.Sum(w)
}

// MR calculates the mean ratio of b:a. It assumes a and
// b are the same length.
func MR(a, b []float64) float64 {
	r := 0.
	for i, v1 := range a {
		r += b[i] / v1
	}
	return r / float64(len(a))
}

// MRWeighted calculates the mean ratio of b:a, weighted by w. It assumes a,
// b, and w are the same length.
func MRWeighted(a, b, w []float64) float64 {
	r := 0.
	for i, v1 := range a {
		r += b[i] / v1 * w[i]
	}
	return r / floats.Sum(w)
}
