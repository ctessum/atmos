package gocart

// Calculate oxidation of SO2 by hydrogen peroxide in clouds as adopted
// from the WRF/Chem file module_gocart_chem.F.
// Reactions are assumed to occur instantaneously.
// Inputs are fraction of grid cell containing clouds (cloudFrac),
// temperature (T [K]), and concentrations of SO2 and H2O2. Units of
// of SO2 and H2O2 are not important as long as both quantities have
// the same units. Outputs are the fractions of SO2 and H2O2 that have
// reacted (so2rxFrac, h2o2rxFrac).
func SulfurAqueousOxidationFraction(cloudFrac, T, so2, h2o2 float64) (
	so2rxFrac, h2o2rxFrac float64) {

	if cloudFrac > 0. && so2 > 0. && T > 258.0 { // Only happens above 258 K
		if so2 > h2o2 {
			so2rxFrac = cloudFrac * h2o2 / so2
			h2o2rxFrac = cloudFrac
		} else {
			so2rxFrac = cloudFrac
			h2o2rxFrac = cloudFrac * so2 / h2o2
		}
	}
	return
}
