package seinfeld

import (
	"math"
	"testing"

	"github.com/ctessum/atmos/wesely1989"
)

// Test calculation of the mean free path
func TestMFP(t *testing.T) {
	T := 298.    // K
	P := 101325. // Pa
	Mu := mu(T)
	lambda := mfp(T, P, Mu)
	if different(lambda, 0.0651e-6, 0.001) {
		t.Fail()
		t.Logf("T=%g: %.3g, %.3g,", T, lambda, 0.0651e-6)
	}
}

// Test calculation of dynamic viscosity
func TestMu(t *testing.T) {
	// K
	Ts := []float64{100, 150, 200, 250, 300, 350, 400, 450, 500, 550, 600, 650,
		700, 750, 800, 850, 900, 950, 1000, 1100, 1200, 1300, 1400, 1500, 1600}
	// kg/(m s) x 10^-5 from http://www.engineeringtoolbox.com/air-absolute-kinematic-viscosity-d_601.html
	Mus := []float64{0.6924, 1.0283, 1.3289, 1.488, 1.983, 2.075, 2.286, 2.484,
		2.671, 2.848, 3.018, 3.177, 3.332, 3.481, 3.625, 3.765, 3.899, 4.023,
		4.152, 4.44, 4.69, 4.93, 5.17, 5.4, 5.63}

	for i, T := range Ts {
		// Allow up to 33% difference from engineering toolbox
		if different(Mus[i]*1.e-5, mu(T), 0.33) {
			t.Fail()
			t.Logf("T=%g: %v, %v,", T, Mus[i]*1.e-5, mu(T))
		}
	}
}
func TestCc(t *testing.T) {
	// Particle diameters [um]
	Dps := []float64{0.001, 0.002, 0.005, 0.01, 0.02, 0.05, 0.1, 0.2, 0.5,
		1.0, 2.0, 5.0, 10.0, 20.0, 50.0, 100.0}
	// Corresponding slip correction factors
	// (Seinfeld and Pandis, 2006 Table 9.3).
	Ccs := []float64{216, 108, 43.6, 22.2, 11.4, 4.95, 2.85, 1.865, 1.326,
		1.164, 1.082, 1.032, 1.016, 1.008, 1.003, 1.0016}

	T := 298.    // K
	P := 101325. // Pa

	for i, Dp := range Dps {
		Cc := cc(Dp*1.e-6, T, P, mu(T))
		if different(Ccs[i], Cc, 1.e-2) {
			t.Fail()
			t.Logf("Dp=%g: %v, %v,", Dp*1.e-6, Ccs[i], Cc)
		}
	}
}
func TestVs(t *testing.T) {
	// Particle diameters [um]
	Dps := []float64{0.01, 0.1, 1, 10} //, 100}
	// Corresponding settling velocities [cm/h]
	// from Seinfeld and Pandis (2006) Figure 9.6
	Vss := []float64{0.025, 0.35, 10.8, 1000} //, 17500}

	const T = 298.    // K
	const P = 101325. // Pa
	//const rhoParticle = 1830. // [kg/m3] Jacobson (2005) Ex. 13.5
	const rhoParticle = 1000.

	for i, Dp := range Dps {
		Mu := mu(T)
		Cc := cc(Dp*1.e-6, T, P, Mu)
		Vs := vs(Dp*1.e-6, rhoParticle, Cc, Mu)
		Vsms := Vss[i] / 3600. / 100. // convert to m/s
		if different(Vsms, Vs, 0.33) {
			t.Fail()
			t.Logf("Dp=%.3g: %.3g, %.3g,", Dp*1.e-6, Vsms, Vs)
		}
	}
}
func TestDryDepGas(t *testing.T) {
	z := 50.      // [m]
	zo := 0.04    // [m]
	ustar := 0.44 // [m/s]
	L := 0.
	T := 298.   // [K]
	rhoA := 1.2 // [kg/m3]
	G := 0.
	theta := 0.
	vd := DryDepGas(z, zo, ustar, L, T, rhoA, G, theta,
		wesely1989.No2Data, wesely1989.Midsummer,
		wesely1989.Urban, false, false, false, false)
	t.Logf("NO2 dry deposition is %.3g cm/s.\n"+
		"Compared to example of 0.1 cm/s in S&P Table 19.1", vd*100)
}
func TestDryDepParticle(t *testing.T) {
	z := 20.      // [m]
	zo := 0.02    // [m]
	ustar := 0.44 // [m/s]
	L := 0.       // [m]
	T := 298.     // [K]
	P := 101325.  // Pa
	rhoA := 1.2   // [kg/m3]
	rhoP := 1000. // [kg/m3]
	results := []float64{0.5, 0.012, 0.02, 11}
	for i, Dp := range []float64{1.e-8, 1.e-7, 1.e-6, 1.e-5} {
		vd := DryDepParticle(z, zo, ustar, L, Dp, T, P, rhoP, rhoA,
			Midsummer, Desert)
		t.Logf("Dp=%.3g, vd=%.3g cm/s, S&P fig.19.2=~%.3g cm/s",
			Dp*1.e6, vd*100, results[i])
	}

}

func different(a, b, tolerance float64) bool {
	if 2*math.Abs(b-a)/math.Abs(b+a) > tolerance {
		return true
	}
	return false
}
