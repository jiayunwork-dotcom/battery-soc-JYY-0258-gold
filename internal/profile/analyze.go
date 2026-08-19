package profile

import (
	"math"

	"battery-soc/internal/coulomb"
)

// Analysis holds computed characteristics of a current profile.
type Analysis struct {
	Steps         int
	DurationH     float64
	NetChargeAh   float64
	PeakCharge    float64
	PeakDischarge float64
	MeanAbsCurrent float64
	RMSCurrent    float64
}

// Analyze computes key characteristics of a current profile.
func Analyze(samples []coulomb.CurrentSample) Analysis {
	if len(samples) == 0 {
		return Analysis{}
	}
	a := Analysis{Steps: len(samples)}
	var sumAbs, sumSq float64

	for _, s := range samples {
		a.DurationH += s.DT
		a.NetChargeAh += s.Current * s.DT
		abs := math.Abs(s.Current)
		sumAbs += abs
		sumSq += s.Current * s.Current
		if s.Current > a.PeakCharge {
			a.PeakCharge = s.Current
		}
		if s.Current < a.PeakDischarge {
			a.PeakDischarge = s.Current
		}
	}
	a.MeanAbsCurrent = sumAbs / float64(len(samples))
	a.RMSCurrent = math.Sqrt(sumSq / float64(len(samples)))
	return a
}

// MaxCRate returns the peak C-rate in a profile given the battery capacity.
func MaxCRate(samples []coulomb.CurrentSample, capacityAh float64) float64 {
	if capacityAh <= 0 {
		return 0
	}
	maxI := 0.0
	for _, s := range samples {
		abs := math.Abs(s.Current)
		if abs > maxI {
			maxI = abs
		}
	}
	return maxI / capacityAh
}
