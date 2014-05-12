package plumerise

import (
	"fmt"
	"math"
	"errors"
)

const (
	g = 9.80665 // m/s2
)

// CalcPlumeRise takes emissions stack height(m), diameter (m), temperature (K),
// and exit velocity (m/s) and calculates the k index of the equivalent
// emissions height after accounting for plume rise.
// Additional required inputs are model layer heights (staggered grid; layerHeights [m]),
// temperature at each layer [K] (unstaggered grid),
// wind speed at each layer [m/s] (unstaggered grid),
// stability class (sClass [0 or 1], unstaggered grid),
// and stability parameter (s1 [unknown units], unstaggered grid).
// Uses the plume rise calculation: ASME (1973), as described in Sienfeld and Pandis,
// ``Atmospheric Chemistry and Physics - From Air Pollution to Climate Change
func PlumeRiseASME(stackHeight, stackDiam, stackTemp,
	stackVel float64,
	layerHeights, temperature, windSpeed, sClass, s1 []float64) (kPlume int, err error) {
	// Find K level of stack
	kStak := 0
	for layerHeights[kStak+1] < stackHeight {
		kStak++
		if kStak >= len(layerHeights)-2 {
			err = AboveModelTop
			return
		}
	}
	deltaH := 0. // Plume rise, (m).
	var calcType string

	airTemp := temperature[kStak]
	windSpd := windSpeed[kStak]

	if (stackTemp-airTemp) < 50. &&
		stackVel > windSpd && stackVel > 10. {
		// Plume is dominated by momentum forces
		calcType = "Momentum"

		deltaH = stackDiam * math.Pow(stackVel, 1.4) / math.Pow(windSpd, 1.4)

	} else { // Plume is dominated by buoyancy forces

		// Bouyancy flux, m4/s3
		F := g * (stackTemp - airTemp) / stackTemp * stackVel *
			math.Pow(stackDiam/2, 2)

		if sClass[kStak] > 0.5 { // stable conditions
			calcType = "Stable"

			deltaH = 29. * math.Pow(
				F/s1[kStak], 0.333333333) /
				math.Pow(windSpd, 0.333333333)

		} else { // unstable conditions
			calcType = "Unstable"

			deltaH = 7.4 * math.Pow(F*math.Pow(stackHeight, 2.),
				0.333333333) / windSpd

		}
	}
	if math.IsNaN(deltaH) {
		err = fmt.Errorf("plume height == NaN\n"+
			"calcType: %v, deltaH: %v, stackDiam: %v,\n"+
			"stackVel: %v, windSpd: %v, stackTemp: %v,\n"+
			"airTemp: %v, stackHeight: %v\n",
			calcType, deltaH, stackDiam, stackVel,
			windSpd, stackTemp, airTemp, stackHeight)
		return
	}

	plumeHeight := stackHeight + deltaH

	// Find K level of plume. If the plume rises above the top model
	// layer, return the top model layer.
	for kPlume = 0; layerHeights[kPlume+1] < plumeHeight; kPlume++ {
		if kPlume >= len(layerHeights)-2 {
			break
		}
	}
	return
}

var AboveModelTop  = errors.New("stack height > top of grid")
