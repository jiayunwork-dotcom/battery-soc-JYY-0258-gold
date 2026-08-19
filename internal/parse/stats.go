package parse

import (
	"math"

	"battery-soc/internal/coulomb"
)

// LogStats holds summary statistics about a current log.
type LogStats struct {
	Count       int
	TotalTimeH  float64
	TotalChargeAh float64
	MaxCurrent  float64
	MinCurrent  float64
	MeanCurrent float64
	RMSCurrent  float64
}

// ComputeStats calculates statistics from a set of current samples.
func ComputeStats(samples []coulomb.CurrentSample) LogStats {
	if len(samples) == 0 {
		return LogStats{}
	}
	s := LogStats{
		Count:      len(samples),
		MaxCurrent: samples[0].Current,
		MinCurrent: samples[0].Current,
	}
	sumI := 0.0
	sumI2 := 0.0
	for _, samp := range samples {
		s.TotalTimeH += samp.DT
		s.TotalChargeAh += samp.Current * samp.DT
		sumI += samp.Current
		sumI2 += samp.Current * samp.Current
		if samp.Current > s.MaxCurrent {
			s.MaxCurrent = samp.Current
		}
		if samp.Current < s.MinCurrent {
			s.MinCurrent = samp.Current
		}
	}
	s.MeanCurrent = sumI / float64(len(samples))
	s.RMSCurrent = math.Sqrt(sumI2 / float64(len(samples)))
	return s
}

// CRate returns the C-rate for a given current and capacity.
// C-rate = |current| / capacity.
func CRate(current, capacityAh float64) float64 {
	if capacityAh <= 0 {
		return 0
	}
	return math.Abs(current) / capacityAh
}

// EstimateTimeToFull estimates hours to reach 100% SOC from current SOC
// at the given constant charge current (simplified, ignores CV phase).
func EstimateTimeToFull(currentSOC, capacityAh, chargeCurrent float64) float64 {
	if chargeCurrent <= 0 || capacityAh <= 0 {
		return 0
	}
	remainingAh := capacityAh * (100 - currentSOC) / 100
	return remainingAh / chargeCurrent
}

// EstimateTimeToEmpty estimates hours to reach 0% SOC from current SOC
// at the given constant discharge current.
func EstimateTimeToEmpty(currentSOC, capacityAh, dischargeCurrent float64) float64 {
	if dischargeCurrent <= 0 || capacityAh <= 0 {
		return 0
	}
	remainingAh := capacityAh * currentSOC / 100
	return remainingAh / dischargeCurrent
}

// DetectRest identifies rest periods (current ~0) in the log and returns
// their start indices and durations.
type RestPeriod struct {
	StartIdx int
	Duration float64 // hours
}

// FindRestPeriods scans the log for consecutive samples where |current|
// is below the threshold. Returns all rest periods longer than minDuration.
func FindRestPeriods(samples []coulomb.CurrentSample, threshold, minDuration float64) []RestPeriod {
	var periods []RestPeriod
	var current *RestPeriod

	for i, s := range samples {
		isRest := math.Abs(s.Current) <= threshold
		if isRest {
			if current == nil {
				current = &RestPeriod{StartIdx: i}
			}
			current.Duration += s.DT
		} else {
			if current != nil && current.Duration >= minDuration {
				periods = append(periods, *current)
			}
			current = nil
		}
	}
	if current != nil && current.Duration >= minDuration {
		periods = append(periods, *current)
	}
	return periods
}
