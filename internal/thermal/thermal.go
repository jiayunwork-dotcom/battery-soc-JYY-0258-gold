// Package thermal provides temperature-dependent corrections for battery
// SOC and capacity estimation. Battery behaviour changes significantly
// with temperature: internal resistance rises at low temperatures and
// capacity fades at extremes.
package thermal

import "math"

// Params holds the thermal model parameters for a battery cell.
type Params struct {
	// NominalTemp is the reference temperature in Celsius (typically 25).
	NominalTemp float64
	// TempCoeffCapacity is the fractional capacity change per degree C
	// below NominalTemp (e.g., 0.005 means 0.5% loss per degree).
	TempCoeffCapacity float64
	// TempCoeffResistance is the fractional resistance increase per
	// degree C below NominalTemp.
	TempCoeffResistance float64
	// MinTemp is the lowest operating temperature (below this, capacity
	// is zero).
	MinTemp float64
	// MaxTemp is the highest safe operating temperature.
	MaxTemp float64
}

// DefaultParams returns typical Li-ion thermal parameters.
func DefaultParams() Params {
	return Params{
		NominalTemp:         25.0,
		TempCoeffCapacity:   0.005,
		TempCoeffResistance: 0.015,
		MinTemp:             -20.0,
		MaxTemp:             60.0,
	}
}

// CapacityFactor returns the fraction of nominal capacity available at
// the given temperature. Result is in [0, 1].
func (p *Params) CapacityFactor(tempC float64) float64 {
	if tempC <= p.MinTemp {
		return 0
	}
	if tempC >= p.NominalTemp {
		return 1.0
	}
	delta := p.NominalTemp - tempC
	factor := 1.0 - p.TempCoeffCapacity*delta
	if factor < 0 {
		return 0
	}
	return factor
}

// EffectiveCapacity returns the temperature-adjusted capacity in Ah.
func (p *Params) EffectiveCapacity(nominalAh, tempC float64) float64 {
	return nominalAh * p.CapacityFactor(tempC)
}

// ResistanceFactor returns the multiplicative factor for internal
// resistance at the given temperature (>= 1.0 for cold, 1.0 at nominal).
func (p *Params) ResistanceFactor(tempC float64) float64 {
	if tempC >= p.NominalTemp {
		return 1.0
	}
	delta := p.NominalTemp - tempC
	return 1.0 + p.TempCoeffResistance*delta
}

// VoltageDropCorrection returns the additional voltage drop due to
// temperature-increased resistance: I * R_nominal * (factor - 1).
func (p *Params) VoltageDropCorrection(current, rNominal, tempC float64) float64 {
	factor := p.ResistanceFactor(tempC)
	return math.Abs(current) * rNominal * (factor - 1)
}

// InRange returns true if the temperature is within operating limits.
func (p *Params) InRange(tempC float64) bool {
	return tempC >= p.MinTemp && tempC <= p.MaxTemp
}

// HeatGeneration estimates the heat generated (in Watts) from resistive
// losses: I^2 * R * factor.
func (p *Params) HeatGeneration(current, rNominal, tempC float64) float64 {
	factor := p.ResistanceFactor(tempC)
	return current * current * rNominal * factor
}

// ThermalRunaway returns true if temperature exceeds the critical
// threshold (typically 150C for Li-ion).
func ThermalRunaway(tempC, criticalTemp float64) bool {
	return tempC >= criticalTemp
}

// CoolingRate estimates the temperature change per hour given ambient
// temperature and a thermal mass constant (simplified Newton's law).
func CoolingRate(cellTemp, ambientTemp, thermalConstant float64) float64 {
	return -thermalConstant * (cellTemp - ambientTemp)
}
