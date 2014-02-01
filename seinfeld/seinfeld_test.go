package seinfeld

import (
	"testing"
)

const (
	P   = 1.   // [atm] Pressure
	T   = 298. // [K] Temperature
	ppb = 1.e-9
)

func TestAqueous(t *testing.T) {
	H2O2 := 1. // [ppb]
	pHarr := []float64{1., 2., 3., 4., 5., 6., 7., 8.}
	wL := 0.000001
	for _, pH := range pHarr {
		k := SulfurH2O2aqueousOxidationRate(H2O2, pH, T, P, wL)
		t.Logf("pH=%v, so2 ox=%.1f%% SO2(g) h-1 (g water/m3 air)-1 "+
			"(compared to 700%% in S&P near eq. 7.84)\n", pH, k*100*3600/(wL*1.e6))
	}
}
