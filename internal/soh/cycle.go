package soh

import "math"

// CycleRecord represents one charge-discharge cycle with its measured
// capacity and environmental conditions.
type CycleRecord struct {
	CycleNumber      int
	DischargeCapAh   float64
	ChargeCapAh      float64
	AvgTemperature   float64
	AvgCRate         float64
}

// CycleHistory holds a sequence of cycle records for trend analysis.
type CycleHistory struct {
	Records         []CycleRecord
	NominalCapacity float64
}

// NewCycleHistory creates a history tracker with the given nominal capacity.
func NewCycleHistory(nominalCap float64) *CycleHistory {
	return &CycleHistory{NominalCapacity: nominalCap}
}

// Add appends a cycle record.
func (h *CycleHistory) Add(r CycleRecord) {
	h.Records = append(h.Records, r)
}

// LatestSOH returns the SOH based on the most recent cycle's discharge
// capacity vs nominal.
func (h *CycleHistory) LatestSOH() float64 {
	if len(h.Records) == 0 {
		return 1.0
	}
	last := h.Records[len(h.Records)-1]
	soh := last.DischargeCapAh / h.NominalCapacity
	if soh > 1 {
		soh = 1
	}
	if soh < 0 {
		soh = 0
	}
	return soh
}

// TrendSlope computes the linear degradation rate (capacity loss per
// cycle) using least-squares regression on discharge capacity.
func (h *CycleHistory) TrendSlope() float64 {
	n := len(h.Records)
	if n < 2 {
		return 0
	}
	// Simple linear regression: y = a + b*x where x=cycle, y=capacity
	var sumX, sumY, sumXY, sumX2 float64
	for _, r := range h.Records {
		x := float64(r.CycleNumber)
		y := r.DischargeCapAh
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	nf := float64(n)
	denom := nf*sumX2 - sumX*sumX
	if denom == 0 {
		return 0
	}
	return (nf*sumXY - sumX*sumY) / denom
}

// PredictEOL estimates the cycle number at which capacity will drop below
// the given threshold fraction (e.g., 0.8 for 80% SOH). Returns -1 if
// the trend is non-negative (no degradation).
func (h *CycleHistory) PredictEOL(threshold float64) int {
	slope := h.TrendSlope()
	if slope >= 0 {
		return -1 // not degrading
	}
	// Starting capacity approximation
	if len(h.Records) == 0 {
		return -1
	}
	startCap := h.Records[0].DischargeCapAh
	targetCap := h.NominalCapacity * threshold
	if startCap <= targetCap {
		return 0
	}
	// cycles = (startCap - targetCap) / |slope|
	cycles := (startCap - targetCap) / math.Abs(slope)
	return int(math.Ceil(cycles))
}

// AverageEfficiency returns the mean coulombic efficiency across all
// cycles (discharge/charge ratio).
func (h *CycleHistory) AverageEfficiency() float64 {
	if len(h.Records) == 0 {
		return 0
	}
	sum := 0.0
	count := 0
	for _, r := range h.Records {
		if r.ChargeCapAh > 0 {
			sum += r.DischargeCapAh / r.ChargeCapAh
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// CapacityVariance returns the variance of discharge capacities (useful
// for detecting inconsistent cycling).
func (h *CycleHistory) CapacityVariance() float64 {
	n := len(h.Records)
	if n < 2 {
		return 0
	}
	mean := 0.0
	for _, r := range h.Records {
		mean += r.DischargeCapAh
	}
	mean /= float64(n)
	variance := 0.0
	for _, r := range h.Records {
		d := r.DischargeCapAh - mean
		variance += d * d
	}
	return variance / float64(n-1)
}
