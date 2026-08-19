// Package soh implements classical State-of-Health estimation for batteries:
// capacity-ratio health, cycle-life remaining fraction, and a linear
// degradation model. No machine learning or external libraries are used.
package soh

import (
	"errors"
	"math"
)

// SOH returns capacity health as currentCap/initialCap, clamped to [0,1].
// It returns an error if initialCap <= 0.
func SOH(initialCap, currentCap float64) (float64, error) {
	if initialCap <= 0 {
		return 0, errors.New("initialCap must be positive")
	}
	ratio := currentCap / initialCap
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return ratio, nil
}

// CycleLife returns the remaining life fraction = 1 - throughputAh/ratedThroughputAh,
// clamped to [0,1]. It returns an error if ratedThroughputAh <= 0.
func CycleLife(throughputAh, ratedThroughputAh float64) (float64, error) {
	if ratedThroughputAh <= 0 {
		return 0, errors.New("ratedThroughputAh must be positive")
	}
	frac := 1 - throughputAh/ratedThroughputAh
	if frac < 0 {
		frac = 0
	}
	if frac > 1 {
		frac = 1
	}
	return frac, nil
}

// DegradationModel returns capacity retention after `cycles` full cycles:
//
//	max(0, 1 - 0.0002*cycles)
func DegradationModel(cycles int) float64 {
	ret := 1 - 0.0002*float64(cycles)
	return math.Max(0, ret)
}
