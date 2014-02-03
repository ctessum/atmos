package gocart

import (
	"math"
)

//*********************************************************************
//  (1) Aerosodynamic resistance Ra and sublayer resistance Rb.       *
//                                                                    *
//  The Reynolds number REYNO diagnoses whether a surface is          *
//  aerodynamically rough (REYNO > 10) or smooth.  Surface is         *
//  rough in all cases except over water with low wind speeds.        *
//                                                                    *
//  For gas species over land and ice (REYNO >= 10) and for aerosol   *
//  species for all surfaces:                                         *
//                                                                    *
//      Ra = 1./VT          (VT from GEOS Kzz at L=1, m/s).           *
//                                                                    *
//  The following equations are from Walcek et al, 1986:              *
//                                                                    *
//  For gas species when REYNO < 10 (smooth), Ra and Rb are combined  *
//  as Ra:                                                            *
//                                                                    *
//      Ra = { ln(ku* z1/Dg) - Sh } / ku*           eq.(13)           *
//                                                                    *
//   where z1 is the altitude at the center of the lowest model layer *
//               (CZ);                                                *
//            Sh is a stability correction function;                  *
//            k  is the von Karman constant (0.4, vK);                *
//            u* is the friction velocity (USTAR).                    *
//                                                                    *
//   Sh is computed as a function of z1 and L       eq ( 4) and (5)): *
//                                                                    *
//    0 < z1/L <= 1:     Sh = -5 * z1/L                               *
//    z1/L < 0:          Sh = exp{ 0.598 + 0.39*ln(E) - 0.09(ln(E))^2}*
//                       where E = min(1,-z1/L) (Balkanski, thesis).  *
//                                                                    *
//   For gas species when REYNO >= 10,                                *
//                                                                    *
//      Rb = 2/ku* (Dair/Dg)**(2/3)                 eq.(12)           *
//      where Dg is the gas diffusivity, and                          *
//            Dair is the air diffusivity.                            *
//                                                                    *
//  For aerosol species, Rb is combined with surface resistance as Rs.*
//                                                                    *
//*********************************************************************

// Calculate GOCART dry deposition for particles as adopted from WRF/Chem
// file module_gocart_drydep.F
// Inputs are the Monin-Obhukov length (obk [m]),
// friction velocity (ustar [m/s]),
// air temperature (T [K]),
// planetary boundary layer height (pblz; m),
// surface roughness length (z0; m), particle radius (r [m]),
// particle density (ρp [kg/m3]), and ambient pressure (P [Pa]).
// Returns dry deposition velocity (vd; m/s).
// This calculation differs from the original gocart calculation in that
// it includes the settling velocity as in Seinfeld and Pandis Eq 19.7
func ParticleDryDep(obk, ustar, T, pblz, z0, r, ρp, P float64) (
	vd float64) {

	ra := calcRa(obk, z0, ustar)
	rs := calcRs(obk, ustar, pblz)
	vs := SettlingVelocity(r, ρp, T, P)

	// Total resistance = Ra + Rs.
	// Set a minimum value for DVEL
	// MIN(vd_aerosol) = 1.0e-4 m/s
	vd = max(1./(ra+rs+ra*rs*vs)+vs, 1.0E-4)
	return
}

// Calculate GOCART dry deposition for gases as adopted from WRF/Chem
// file module_gocart_drydep.F
// Inputs are the Monin-Obhukov length (obk, [m]),
// friction velocity (ustar [m/s]),
// planetary boundary layer height (pblz [m]),
// surface roughness length (z0 [m]), and ratio of H2O to gas-of-interest
// diffusivities Dratio (Dratio [m2/s]).
// Returns dry deposition velocity (vd [m/s]).
func GasDryDep(obk, ustar, pblz, z0, Dratio float64) (vd float64) {

	ra := calcRa(obk, z0, ustar)
	rb := calcRb(ustar, Dratio)
	rs := calcRs(obk, ustar, pblz)

	// Total resistance = Ra + Rb + Rs.
	//  Aerosol species, Rs here is the combination of Rb and Rs.
	// Set a minimum value for DVEL
	// MIN(VdSO2)      = 2.0e-3 m/s  over ice
	//                 = 3.0e-3 m/s  over land
	vd = max(1./(ra+rb+rs), 3.0E-3)
	return
}

// Calculate aerosodynamic resistance. Assume that reynolds number
// is always >10.
func calcRa(obk, z0, ustar float64) (Ra float64) {
	cz := 2. // Surface layer height (I think)
	frac := cz / obk
	if frac > 1.0 {
		frac = 1.0
	}
	var psi_h float64
	if frac > 0.0 && frac <= 1.0 {
		psi_h = -5.0 * frac
	} else if frac < 0.0 {
		eps := min(1.0, -frac)
		logmfrac := math.Log(eps)
		psi_h = math.Exp(0.598 + 0.39*logmfrac - 0.09*math.Pow(logmfrac, 2.))
	}
	Ra = (math.Log(cz/z0) - psi_h) / (vK * ustar)
	return
}

// Ratios H2O diffusivity to other gas diffusivity from
// Seinfeld and Pandis table 19.4
var DratioForRb = map[string]float64{"SO2": 1.89, "O3": 1.63,
	"NO2": 1.6, "NO": 1.29, "H2O2": 1.37, "NH3": 0.97, "HCHO": 1.29}

// Calculate quasi-laminar resistance
func calcRb(ustar, Dratio float64) (Rb float64) {
	Rb = 2. / vK / ustar * math.Pow(Dratio, 0.66666667)
	return
}

// Calculate sublayer resistance
func calcRs(obk, ustar, pblz float64) (Rs float64) {
	vds := 0.002 * ustar
	if obk < 0.0 {
		vds = vds * (1.0 + math.Pow(-300.0/obk, 0.6667))
	}

	czh := pblz / obk
	if czh < -30.0 {
		vds = 0.0009 * ustar * math.Pow(-czh, 0.6667)
	}

	// --Set Vds to be less than VDSMAX (entry in input file divided --
	//   by 1.E4). VDSMAX is taken from Table 2 of Walcek et al. [1986].
	//   Invert to get corresponding R
	Rs = 1.0 / min(vds, 2.0e-3)

	// ------ Set max and min values for bulk surface resistances ------
	Rs = max(1.0, min(Rs, 9.9990e+3))
	return
}
