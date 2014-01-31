package gocart

import (
	"math"
)

// Calculate particle terminal settling velocity as
// adopted from WRF/Chem file module_gocart_settling.F.
// Inputs are effective particle radius (Reff [m]),
// particle density (ρ [kg/m3],
// air temperature (T [K]) and air pressure (P [Pa]).
// Returns settling velocity (vs [m/s]).
func SettlingVelocity(Reff, ρ, T, P float64) (vs float64) {

	// Dynamic viscosity
	c_stokes := 1.458E-6 * math.Pow(T, 1.5) / (T + 110.4)

	// Mean free path as a function of pressure (mb) and
	// temperature (K)
	free_path := 1.1E-3 / P / math.Sqrt(T) // m

	// Slip Correction Factor
	c_cun := 1.0 + free_path/Reff*
		(1.257+0.4*math.Exp(-1.1*Reff/free_path))

	// Corrected dynamic viscosity (kg/m/s)
	viscosity := c_stokes / c_cun

	// Settling velocity
	vs = 2.0 / 9.0 * g * ρ * math.Pow(Reff, 2.) / viscosity

	return
}
