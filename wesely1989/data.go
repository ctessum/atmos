package wesely1989

// r_i represents the minimum bulk canopy stomatal resistances for water vapor.
var r_i = [][]float64{
	{9999, 60, 120, 70, 130, 100, 9999, 9999, 80, 100, 150},
	{9999, 9999, 9999, 9999, 250, 500, 9999, 9999, 9999, 9999, 9999},
	{9999, 9999, 9999, 9999, 250, 500, 9999, 9999, 9999, 9999, 9999},
	{9999, 9999, 9999, 9999, 400, 800, 9999, 9999, 9999, 9999, 9999},
	{9999, 120, 240, 140, 250, 190, 9999, 9999, 160, 200, 300}}

// r_lu signifies leaf cuticles in healthy vegetation and otherwise
// the outer surfaces in the upper canopy.
var r_lu = [][]float64{
	{9999, 2000, 2000, 2000, 2000, 2000, 9999, 9999, 2500, 2000, 4000},
	{9999, 9000, 9000, 9000, 4000, 8000, 9999, 9999, 9000, 9000, 9000},
	{9999, 9999, 9000, 9000, 4000, 8000, 9999, 9999, 9000, 9000, 9000},
	{9999, 9999, 9999, 9999, 6000, 9000, 9999, 9999, 9000, 9000, 9000},
	{9999, 4000, 4000, 4000, 2000, 3000, 9999, 9999, 4000, 4000, 8000}}

// r_ac signifies transfer that depends only on canopy height and density.
var r_ac = [][]float64{
	{100, 200, 100, 2000, 2000, 2000, 0, 0, 300, 150, 200},
	{100, 150, 100, 1500, 2000, 1700, 0, 0, 200, 120, 140},
	{100, 10, 100, 1000, 2000, 1500, 0, 0, 100, 50, 120},
	{100, 10, 10, 1000, 2000, 1500, 0, 0, 50, 10, 50},
	{100, 50, 80, 1200, 2000, 1500, 0, 0, 200, 60, 120}}

// r_gs signifies uptake at the "ground" by soil, leaf litter, snow, water etc.
// 'S' and 'O' stand for SO2 and O3 respectively.
var r_gsS = [][]float64{
	{400, 150, 350, 500, 500, 100, 0, 1000, 0, 220, 400},
	{400, 200, 350, 500, 500, 100, 0, 1000, 0, 300, 400},
	{400, 150, 350, 500, 500, 200, 0, 1000, 0, 200, 400},
	{100, 100, 100, 100, 100, 100, 0, 1000, 100, 100, 50},
	{500, 150, 350, 500, 500, 200, 0, 1000, 0, 250, 400}}

var r_gsO = [][]float64{
	{300, 150, 200, 200, 200, 300, 2000, 400, 1000, 180, 200},
	{300, 150, 200, 200, 200, 300, 2000, 400, 800, 180, 200},
	{300, 150, 200, 200, 200, 300, 2000, 400, 1000, 180, 200},
	{600, 3500, 3500, 3500, 3500, 3500, 2000, 400, 3500, 3500, 3500},
	{300, 150, 200, 200, 200, 300, 2000, 400, 1000, 180, 200}}

// r_cl is meant to account for uptake pathways at the leaves, bark, etc.
// 'S' and 'O' stand for SO2 and O3 respectively.
var r_clS = [][]float64{
	{9999, 2000, 2000, 2000, 2000, 2000, 9999, 9999, 2500, 2000, 4000},
	{9999, 9000, 9000, 9000, 2000, 4000, 9999, 9999, 9000, 9000, 9000},
	{9999, 9999, 9000, 9000, 3000, 6000, 9999, 9999, 9000, 9000, 9000},
	{9999, 9999, 9999, 9000, 200, 400, 9999, 9999, 9000, 9999, 9000},
	{9999, 4000, 4000, 4000, 2000, 3000, 9999, 9999, 4000, 4000, 8000}}

var r_clO = [][]float64{
	{9999, 1000, 1000, 1000, 1000, 1000, 9999, 9999, 1000, 1000, 1000},
	{9999, 400, 400, 400, 1000, 600, 9999, 9999, 400, 400, 400},
	{9999, 1000, 400, 400, 1000, 600, 9999, 9999, 800, 600, 600},
	{9999, 1000, 1000, 400, 1500, 600, 9999, 9999, 800, 1000, 800},
	{9999, 1000, 500, 500, 1500, 700, 9999, 9999, 600, 800, 800}}

// Holder for gas properties from Wesely (1989) Table 2.
type GasData struct {
	Dh2oPerDx float64 // ratio of water to chemical-of-interest diffusivities [-]
	Hstar     float64 // effective Henry's law coefficient [M atm-1]
	Fo        float64 // reactivity factor [-]
}

// Properties of various gases from Wesely (1989) Table 2.
var (
	So2Data  = &GasData{1.9, 1.e5, 0}
	O3Data   = &GasData{1.6, 0.01, 1}
	No2Data  = &GasData{1.6, 0.01, 0.1}
	NoData   = &GasData{1.3, 2.e-3, 0}
	Hno3Data = &GasData{1.9, 1.e14, 0}
	H2o2Data = &GasData{1.4, 1.e5, 1}
	AldData  = &GasData{1.6, 15, 0}    // Acetaldehyde (aldehyde class)
	HchoData = &GasData{1.3, 6.e3, 0}  // Formaldehyde
	OpData   = &GasData{1.6, 240, 0.1} // Methyl hydroperoxide (organic peroxide class)
	PaaData  = &GasData{2.0, 540, 0.1} // Peroxyacetyl nitrate
	OraData  = &GasData{1.6, 4.e6, 0}  // Formic acid (organic acid class)
	Nh3Data  = &GasData{1.0, 2.e4, 0}
	PanData  = &GasData{2.6, 3.6, 0.1}  // Peroxyacetyl nitrate
	Hno2Data = &GasData{1.6, 1.e5, 0.1} // Nitrous acid
)
