package gocart

import (
	"math"
)

const (
	g  = 9.81 // m/s2; gravitational acceleration
	vK = 0.4  // von Karman's constant
)

//Compute the the Monin-Obhukov length.
//The direct computation of the Monin-Obhukov length is:
//
//			- Air density * Cp * T(surface air) * Ustar^3
//	OBK =   ----------------------------------------------
//				vK   * g  * Sensible Heat flux
//
//
//	Cp = 1000 J/kg/K    = specific heat at constant pressure
//	vK = 0.4            = von Karman's constant
func ObhukovLen(hflux, airden, Ts, ustar float64) (obk float64) {
	if math.Abs(hflux) <= 1.e-5 {
		obk = 1.0E5
	} else {
		obk = -airden * 1000.0 * Ts * math.Pow(ustar, 3.) /
			(vK * g * hflux)
	}
	return
}

func min(v1, v2 float64) float64 {
	if v1 < v2 {
		return v1
	} else {
		return v2
	}
}
func max(v1, v2 float64) float64 {
	if v1 > v2 {
		return v1
	} else {
		return v2
	}
}
