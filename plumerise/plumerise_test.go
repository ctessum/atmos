package plumerise

import (
	"errors"
	"math"
	"testing"
)

func TestFindStackLayer(t *testing.T) {
	// Layer heights are staggered.
	var layerHeights = []float64{0, 10, 20, 30, 40}
	// go test -cover
	type test struct {
		height float64
		layer  int
		err    error
	}
	var tests = []test{
		{
			height: 0,
			layer:  0,
		},
		{
			height: 9,
			layer:  0,
		},
		{
			height: 15,
			layer:  1,
		},
		{
			height: 35,
			layer:  3,
		},
		{
			height: 45,
			layer:  3,
			err:    ErrAboveModelTop,
		},
	}
	for _, tt := range tests {
		layer, err := findLayer(layerHeights, tt.height)
		if layer != tt.layer {
			t.Errorf("height %g should be layer %d but is layer %d", tt.height, tt.layer, layer)
		}
		if err != tt.err {
			t.Errorf("height %g error should be %v but is %v", tt.height, tt.err, err)
		}
	}
}

func TestCalcDeltaH(t *testing.T) {
	// Layer heights are staggered.
	var temperature = []float64{50, 10, 15, 15}
	var windSpeed = []float64{10, 12.5, 15, 14}
	var sClass = []float64{0.25, 0.75, 0.4, 0.6}
	var s1 = []float64{0.2, 0.5, 1.0, math.NaN()}
	var stackHeight float64 = 100
	var stackTemp float64 = 80
	var stackVel float64 = 20
	var stackDiam float64 = 10

	// go test -cover
	type test struct {
		stackLayer int
		err        error
		deltaH     float64
	}
	var tests = []test{
		{
			stackLayer: 0,
			err:        nil,
			deltaH:     26.3901,
		},
		{
			stackLayer: 1,
			err:        nil,
			deltaH:     255.8218,
		},
		{
			stackLayer: 2,
			err:        nil,
			deltaH:     168.4916,
		},
		{
			stackLayer: 3,
			err: errors.New("plume height == NaN\n" +
				"deltaH: NaN, stackDiam: 10,\n" +
				"stackVel: 20, windSpd: 14, stackTemp: 80,\n" +
				"airTemp: 15, stackHeight: 100\n"),
			deltaH: 0,
		},
	}
	for _, tt := range tests {
		deltaH, err := calcDeltaH(tt.stackLayer, temperature, windSpeed, sClass, s1,
			stackHeight, stackTemp, stackVel, stackDiam)
		if math.Abs(deltaH-tt.deltaH) > 0.0001 {
			t.Errorf("stackLayer %v should have deltaH %v, but is %v", tt.stackLayer, tt.deltaH, deltaH)
		}
		if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
			t.Errorf("stackLayer %v error should be %v but is %v", tt.stackLayer, tt.err, err)
		} else if err != nil && tt.err == nil {
			t.Errorf("stackLayer %v error should be %v but is %v", tt.stackLayer, tt.err, err)
		} else if err == nil && tt.err != nil {
			t.Errorf("stackLayer %v error should be %v but is %v", tt.stackLayer, tt.err, err)
		}
	}
}

func TestCalcDeltaHPrecomputed(t *testing.T) {
	// Layer heights are staggered.
	var temperature = []float64{80, 80, 100, 20, 20, 20, 20}
	var windSpeed = []float64{10, 10, 40, 40, 40, 40, 40}
	var sClass = []float64{1, 1, 1, 1, 1, 0.25, 0.25}
	var s1 = []float64{1, 1, 1, 1, 1, 1, 1}
	var windSpeedMinusOnePointFour = []float64{math.NaN(), 10, 10, 1, 1, 1, 1}
	var windSpeedMinusThird = []float64{1, 1, 1, math.NaN(), 1, 1, 1}
	var windSpeedInverse = []float64{1, 1, 1, 1, 1, math.NaN(), 1}
	var stackHeight float64 = 100
	var stackTemp float64 = 100
	var stackVel float64 = 20
	var stackDiam float64 = 10

	// go test -cover
	type test struct {
		stackLayer int
		err        error
		deltaH     float64
	}
	var tests = []test{
		{
			stackLayer: 0,
			err: errors.New("plumerise: momentum-dominated deltaH is NaN. " +
				"stackDiam: 10, stackVel: 20, windSpeedMinusOnePointFour: NaN"),
			deltaH: 28782.0921,
		},
		{
			stackLayer: 1,
			err:        nil,
			deltaH:     6628.9080,
		},
		{
			stackLayer: 2,
			err:        nil,
			deltaH:     0,
		},
		{
			stackLayer: 3,
			err: errors.New("plumerise: stable bouyancy-dominated deltaH is NaN. " +
				"F: 6537.766666666667, s1: 1, windSpeedMinusThird: NaN"),
			deltaH: 9734.7946,
		},
		{
			stackLayer: 4,
			err:        nil,
			deltaH:     542.2602,
		},
		{
			stackLayer: 5,
			err: errors.New("plumerise: unstable bouyancy-dominated deltaH is NaN. " +
				"F: 6537.766666666667, stackHeight: 100, windSpeedInverse: NaN"),
			deltaH: 9734.7946,
		},
		{
			stackLayer: 6,
			err:        nil,
			deltaH:     2981.0884,
		},
	}
	for _, tt := range tests {
		deltaH, err := calcDeltaHPrecomputed(tt.stackLayer, temperature, windSpeed, sClass, s1,
			stackHeight, stackTemp, stackVel, stackDiam, windSpeedMinusOnePointFour,
			windSpeedMinusThird, windSpeedInverse)
		if math.Abs(deltaH-tt.deltaH) > 0.0001 {
			t.Errorf("stackLayer %v should have deltaH %v, but is %v", tt.stackLayer, tt.deltaH, deltaH)
		}
		if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
			t.Errorf("stackLayer %v error should be %v but is %v", tt.stackLayer, tt.err, err)
		} else if err != nil && tt.err == nil {
			t.Errorf("stackLayer %v error should be %v but is %v", tt.stackLayer, tt.err, err)
		} else if err == nil && tt.err != nil {
			t.Errorf("stackLayer %v error should be %v but is %v", tt.stackLayer, tt.err, err)
		}
	}
}

func TestASME(t *testing.T) {
	// Layer heights are staggered.
	var temperature = []float64{50, 10, 15, 15}
	var windSpeed = []float64{10, 12.5, 15, 14}
	var sClass = []float64{0.25, 0.75, 0.4, 0.6}
	var s1 = []float64{0.2, 0.5, 1.0, math.NaN()}
	var stackTemp float64 = 100
	var stackVel float64 = 20
	var stackDiam float64 = 10
	var layerHeights = []float64{0, 10, 20, 30, 40}

	// go test -cover
	type test struct {
		stackHeight float64
		err         error
		plumeLayer  int
		plumeHeight float64
	}
	var tests = []test{
		{
			stackHeight: 0,
			err:         nil,
			plumeLayer:  0,
			plumeHeight: 0,
		},
		{
			stackHeight: 10,
			err:         ErrAboveModelTop,
			plumeLayer:  3,
			plumeHeight: 56.3146,
		},
		{
			stackHeight: 20,
			err:         ErrAboveModelTop,
			plumeLayer:  3,
			plumeHeight: 278.2353,
		},
		{
			stackHeight: 30,
			err:         ErrAboveModelTop,
			plumeLayer:  3,
			plumeHeight: 106.6521,
		},
		{
			stackHeight: 40,
			err: errors.New("plume height == NaN\n" +
				"deltaH: NaN, stackDiam: 10,\n" +
				"stackVel: 20, windSpd: 14, stackTemp: 100,\n" +
				"airTemp: 15, stackHeight: 40\n"),
			plumeLayer:  0,
			plumeHeight: 0,
		},
		{
			stackHeight: 50,
			err:         ErrAboveModelTop,
			plumeLayer:  3,
			plumeHeight: 50,
		},
		{
			stackHeight: 60,
			err:         ErrAboveModelTop,
			plumeLayer:  3,
			plumeHeight: 60,
		},
	}
	for _, tt := range tests {
		plumeLayer, plumeHeight, err := ASME(tt.stackHeight, stackDiam, stackTemp,
			stackVel, layerHeights, temperature, windSpeed,
			sClass, s1)
		if plumeLayer != tt.plumeLayer {
			t.Errorf("%v should have plumeLayer %v, but is %v", tt.stackHeight, tt.plumeLayer, plumeLayer)
		}
		if math.Abs(plumeHeight-tt.plumeHeight) > 0.0001 {
			t.Errorf("%v should have plumeHeight %v, but is %v", tt.stackHeight, tt.plumeHeight, plumeHeight)
		}
		if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
			t.Errorf("stackHeight %v error should be %v but is %v", tt.stackHeight, tt.err, err)
		} else if err != nil && tt.err == nil {
			t.Errorf("stackHeight %v error should be %v but is %v", tt.stackHeight, tt.err, err)
		} else if err == nil && tt.err != nil {
			t.Errorf("stackHeight %v error should be %v but is %v", tt.stackHeight, tt.err, err)
		}
	}
}

func TestASMEPrecomputed(t *testing.T) {
	// Layer heights are staggered.
	var temperature = []float64{80, 80, 100, 20, 20, 20, 20}
	var windSpeed = []float64{10, 10, 40, 40, 40, 40, 40}
	var sClass = []float64{1, 1, 1, 1, 1, 0.25, 0.25}
	var s1 = []float64{1, 1, 1, 1, 1, 1, 1}
	var windSpeedMinusOnePointFour = []float64{math.NaN(), 10, 10, 1, 1, 1, 1}
	var windSpeedMinusThird = []float64{1, 1, 1, math.NaN(), 1, 1, 1}
	var windSpeedInverse = []float64{1, 1, 1, 1, 1, math.NaN(), 1}
	var stackTemp float64 = 100
	var stackVel float64 = 20
	var stackDiam float64 = 10
	var layerHeights = []float64{0, 10, 20, 30, 40, 50, 60}

	// go test -cover
	type test struct {
		stackHeight float64
		err         error
		plumeLayer  int
		plumeHeight float64
	}
	var tests = []test{
		{
			stackHeight: 0,
			err: errors.New("plumerise: momentum-dominated deltaH is NaN. " +
				"stackDiam: 10, stackVel: 20, windSpeedMinusOnePointFour: NaN"),
			plumeLayer:  0,
			plumeHeight: 0,
		},
		{
			stackHeight: 10,
			err: errors.New("plumerise: momentum-dominated deltaH is NaN. " +
				"stackDiam: 10, stackVel: 20, windSpeedMinusOnePointFour: NaN"),
			plumeLayer:  0,
			plumeHeight: 0,
		},
		{
			stackHeight: 20,
			err:         ErrAboveModelTop,
			plumeLayer:  5,
			plumeHeight: 6648.9080,
		},
		{
			stackHeight: 30,
			err:         nil,
			plumeLayer:  2,
			plumeHeight: 30,
		},
		{
			stackHeight: 40,
			err: errors.New("plumerise: stable bouyancy-dominated deltaH is NaN. " +
				"F: 6537.766666666667, s1: 1, windSpeedMinusThird: NaN"),
			plumeLayer:  0,
			plumeHeight: 0,
		},
		{
			stackHeight: 50,
			err:         ErrAboveModelTop,
			plumeLayer:  5,
			plumeHeight: 592.2602,
		},
		{
			stackHeight: 60,
			err: errors.New("plumerise: unstable bouyancy-dominated deltaH is NaN. " +
				"F: 6537.766666666667, stackHeight: 60, windSpeedInverse: NaN"),
			plumeLayer:  0,
			plumeHeight: 0,
		},
		{
			stackHeight: 70,
			err:         ErrAboveModelTop,
			plumeLayer:  5,
			plumeHeight: 70,
		},
	}
	for _, tt := range tests {
		plumeLayer, plumeHeight, err := ASMEPrecomputed(tt.stackHeight, stackDiam, stackTemp,
			stackVel, layerHeights, temperature, windSpeed,
			sClass, s1, windSpeedMinusOnePointFour, windSpeedMinusThird,
			windSpeedInverse)
		if plumeLayer != tt.plumeLayer {
			t.Errorf("stackHeight %v should have plumeLayer %v, but is %v", tt.stackHeight, tt.plumeLayer, plumeLayer)
		}
		if math.Abs(plumeHeight-tt.plumeHeight) > 0.0001 {
			t.Errorf("stackHeight %v should have plumeHeight %v, but is %v", tt.stackHeight, tt.plumeHeight, plumeHeight)
		}
		if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
			t.Errorf("stackHeight %v error should be %v but is %v", tt.stackHeight, tt.err, err)
		} else if err != nil && tt.err == nil {
			t.Errorf("stackHeight %v error should be %v but is %v", tt.stackHeight, tt.err, err)
		} else if err == nil && tt.err != nil {
			t.Errorf("stackHeight %v error should be %v but is %v", tt.stackHeight, tt.err, err)
		}
	}
}
