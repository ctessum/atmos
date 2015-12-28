package wesely1989

import (
	"fmt"
	"math"
	"testing"
)

// Results from Wesely (1989) table 3; updated to values
// from Walmsley (1996) table 1.
var (
	SO2 = [][]float64{
		{130, 140, 160, 380, 1000, 100, 1200},
		{1400, 1400, 1400, 1400, 1500, 100, 1300},
		{1100, 1100, 1100, 1100, 1200, 90, 1000},
		{1000, 1000, 1000, 1000, 1100, 1100, 1100},
		{270, 290, 330, 620, 1100, 90, 1000}}

	O3 = [][]float64{
		{100, 110, 130, 320, 960, 960, 580},
		{430, 470, 520, 710, 1300, 950, 580},
		{390, 420, 460, 610, 960, 770, 510},
		{560, 620, 710, 1100, 3200, 3200, 3200},
		{180, 200, 230, 440, 950, 820, 530}}

	NO2 = [][]float64{
		{120, 130, 160, 480, 2900, 2700, 2300},
		{1900, 1900, 1900, 2000, 2700, 2500, 2200},
		{1700, 1700, 1800, 1900, 2400, 2300, 2000},
		{3900, 4000, 4100, 4500, 9999, 9999, 9999},
		{270, 290, 350, 850, 2500, 2300, 2000}}

	H2O2 = [][]float64{
		{90, 90, 110, 250, 640, 90, 80},
		{400, 430, 480, 650, 1100, 90, 90},
		{370, 390, 430, 550, 840, 90, 80},
		{400, 430, 470, 620, 1000, 1000, 1000},
		{160, 170, 200, 370, 750, 90, 80}}

	ALD = [][]float64{
		{330, 340, 370, 800, 9999, 9999, 9999},
		{9999, 9999, 9999, 9999, 9999, 9999, 9999},
		{9999, 9999, 9999, 9999, 9999, 9999, 9999},
		{9999, 9999, 9999, 9999, 9999, 9999, 9999},
		{520, 550, 630, 1700, 9999, 9999, 9999}}

	HCHO = [][]float64{
		{100, 110, 140, 450, 6700, 1400, 1400},
		{8700, 8700, 8700, 8700, 8700, 1400, 1400},
		{8300, 8300, 8300, 8300, 8400, 1400, 1400},
		{2900, 2900, 2900, 2900, 2900, 2900, 2900},
		{250, 270, 340, 1000, 7500, 1400, 1400}}

	OP = [][]float64{
		{120, 130, 160, 480, 2800, 2500, 2200},
		{1900, 1900, 1900, 2000, 2700, 2400, 2000},
		{1700, 1700, 1800, 1800, 2400, 2100, 1900},
		{3700, 3700, 3800, 4200, 8600, 8600, 8600},
		{270, 290, 350, 850, 2500, 2200, 1900}}

	PAA = [][]float64{
		{150, 160, 200, 580, 2800, 2400, 2000},
		{1900, 1900, 1900, 2000, 2700, 2200, 1900},
		{1700, 1700, 1700, 1800, 2400, 2000, 1800},
		{3400, 3400, 3500, 3800, 7200, 7200, 7200},
		{330, 350, 420, 960, 2400, 2100, 1800}}

	ORA = [][]float64{
		{30, 30, 30, 40, 50, 10, 10},
		{140, 140, 150, 170, 190, 10, 10},
		{130, 140, 140, 160, 180, 10, 10},
		{310, 340, 390, 550, 910, 910, 910},
		{60, 60, 70, 80, 90, 10, 10}}

	NH3 = [][]float64{
		{80, 80, 100, 320, 2700, 430, 430},
		{3400, 3400, 3400, 3400, 3400, 440, 440},
		{3000, 3000, 3000, 3000, 3100, 430, 430},
		{1500, 1500, 1500, 1500, 1500, 1500, 1500},
		{180, 200, 240, 680, 2800, 430, 430}}
	PAN = [][]float64{
		{190, 210, 250, 700, 2900, 2700, 2300},
		{1900, 1900, 1900, 2000, 2700, 2500, 2200},
		{1700, 1700, 1800, 1900, 2400, 2300, 2000},
		{3900, 4000, 4100, 4500, 9999, 9999, 9999},
		{410, 430, 510, 1100, 2500, 2300, 2000}}
	HNO2 = [][]float64{
		{110, 120, 140, 330, 950, 90, 90},
		{1000, 1000, 1000, 1100, 1400, 90, 90},
		{860, 860, 870, 910, 1100, 90, 90},
		{820, 830, 830, 870, 1000, 1000, 1000},
		{220, 240, 280, 530, 1000, 90, 90}}
)

func TestWesely(t *testing.T) {
	const iLandUse = 3                       // deciduous forest
	Ts := []float64{25, 10, 2, 0, 10}        // Surface Temperature [C]
	Garr := []float64{800, 500, 300, 100, 0} // Solar radiation [W m-2]
	Θ := 0.                                  // Slope [radians]

	polNames := []string{"SO2", "O3", "NO2", "H2O2", "ALD", "HCHO", "OP", "PAA", "ORA", "NH3", "PAN", "HNO2"}
	testData := [][][]float64{SO2, O3, NO2, H2O2, ALD, HCHO, OP, PAA, ORA, NH3, PAN, HNO2}
	gasData := []*GasData{So2Data, O3Data, No2Data, H2o2Data, AldData, HchoData, OpData, PaaData, OraData, Nh3Data, PanData, Hno2Data}

	for iPol, pol := range polNames {
		polData := testData[iPol]
		isSO2, isO3 := false, false
		if pol == "SO2" {
			isSO2 = true
		}
		if pol == "O3" {
			isO3 = true
		}
		for iSeason := SeasonCategory(0); iSeason < 5; iSeason++ {
			for ig, G := range Garr {
				r_c := SurfaceResistance(gasData[iPol], G, Ts[iSeason], Θ,
					iSeason, iLandUse, false, false, isSO2, isO3)
				if different(r_c, polData[iSeason][ig]) {
					fmt.Printf("%v, %v, %v: %.0f, %g\n", pol, iSeason, G, r_c, polData[iSeason][ig])
					t.Fail()
				}
			}
			r_c := SurfaceResistance(gasData[iPol], 0., Ts[iSeason], Θ,
				iSeason, iLandUse, false, true, isSO2, isO3) // dew
			if different(r_c, polData[iSeason][5]) {
				fmt.Printf("%v, %v, %v: %.0f, %g\n", pol, iSeason, "dew", r_c, polData[iSeason][5])
				t.Fail()
			}

			r_c = SurfaceResistance(gasData[iPol], 0., Ts[iSeason], Θ,
				iSeason, iLandUse, true, false, isSO2, isO3) // rain
			if different(r_c, polData[iSeason][6]) {
				fmt.Printf("%v, %v, %v: %.0f, %g\n", pol, iSeason, "rain", r_c, polData[iSeason][6])
				t.Fail()
			}
		}
	}
}

func different(a, b float64) bool {
	c := math.Abs(a - b)
	return c/b > .1 && c >= 11.
}
