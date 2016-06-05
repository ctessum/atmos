package plumerise

import "testing"

func TestFindStackLayer(t *testing.T) {
	// Layer heights are staggered.
	var layerHeights = []float64{0, 10, 20, 30, 40}

	type test struct {
		height float64
		layer  int
		err    error
	}
	var tests = []test{
		test{
			height: 0,
			layer:  0,
		},
		test{
			height: 9,
			layer:  0,
		},
		test{
			height: 15,
			layer:  1,
		},
		test{
			height: 35,
			layer:  3,
		},
		test{
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
