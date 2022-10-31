package tests

import (
	"kablamo/lib/functions"
	"kablamo/lib/parallel"
	"kablamo/lib/serial"
	"testing"

	"github.com/stretchr/testify/assert"
)

//replaces every even number with number+99
// in every step for nSteps, a step sum is calculated as :
// 	if i%3 == 0, then nothing added to sum,
// 	   i%3==1, then i*4 is added to sum
//     i%3==2 then i*2 is added
//

func Test_Process(t *testing.T) {
	testcases := []struct {
		name        string
		start       int
		end         int
		nSteps      int
		expectedSum int32
		processFn   func(int, *[]int, func(int32, int32) int32) int32
	}{
		{
			// series: [1,101,201]
			// stepSum: [4,303] = 307
			// repeat 5 times: 307 * 4  = 1228
			name:        "Simple Parallel",
			start:       1,
			end:         200,
			nSteps:      5,
			expectedSum: 1228,
			processFn:   parallel.Process,
		},
		{
			// series: [1,101,201]
			// stepSum: [4,303] = 307
			// repeat 5 times: 307 * 4  = 1228
			name:        "Simple Serial",
			start:       1,
			end:         200,
			nSteps:      5,
			expectedSum: 1228,
			processFn:   serial.Process,
		},
	}
	for _, tc := range testcases {
		series := functions.Series(tc.start, tc.end)
		result := tc.processFn(tc.nSteps, series, functions.StepFunction)
		assert.Equal(t, tc.expectedSum, result)
	}
}
