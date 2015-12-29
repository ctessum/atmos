package advect

import (
	"math"
	"testing"
)

func TestUpwind(t *testing.T) {
	const tolerance = 1.e-8
	type test struct {
		u, cm1, c, dx float64
		dc            float64 // The expected result
	}
	tests := []test{
		test{u: 10, cm1: 2, c: 1, dx: 1000, dc: 0.02},
		test{u: 10, cm1: 1, c: 2, dx: 1000, dc: 0.01},
		test{u: -10, cm1: 1, c: 2, dx: 1000, dc: -0.02},
		test{u: -10, cm1: 2, c: 1, dx: 1000, dc: -0.01},
	}
	for _, tt := range tests {
		result := UpwindFlux(tt.u, tt.cm1, tt.c, tt.dx)
		if different(result, tt.dc, tolerance) {
			t.Fatalf("u=%g, cm1=%g, c=%g, dx=%g, dc should equal %g but equals %g",
				tt.u, tt.cm1, tt.c, tt.dx, tt.dc, result)
		}
	}
}

func different(a, b, tolerance float64) bool {
	if 2*math.Abs(a-b)/math.Abs(a+b) > tolerance || math.IsNaN(a) || math.IsNaN(b) {
		return true
	}
	return false
}
