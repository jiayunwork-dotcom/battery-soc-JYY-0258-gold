// Package ocv provides Open Circuit Voltage to SOC lookup tables and
// interpolation for lithium-ion battery cells. The OCV-SOC relationship
// is the thermodynamic foundation for voltage-based SOC estimation.
package ocv

import (
	"errors"
	"fmt"
	"math"
)

// Table is a piecewise-linear OCV-SOC lookup table. SOC values are in
// percent [0,100] and voltages in Volts. Both slices must be the same
// length and SOC must be strictly increasing.
type Table struct {
	SOC     []float64
	Voltage []float64
}

// DefaultLiFePO4 returns a typical LiFePO4 OCV-SOC table with 11 points.
func DefaultLiFePO4() *Table {
	return &Table{
		SOC:     []float64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		Voltage: []float64{2.50, 3.00, 3.20, 3.22, 3.24, 3.26, 3.28, 3.30, 3.32, 3.35, 3.65},
	}
}

// DefaultNMC returns a typical NMC (Nickel Manganese Cobalt) OCV-SOC table.
func DefaultNMC() *Table {
	return &Table{
		SOC:     []float64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		Voltage: []float64{3.00, 3.40, 3.55, 3.62, 3.68, 3.73, 3.80, 3.88, 3.98, 4.10, 4.20},
	}
}

// Validate checks that the table is well-formed: same length, >=2 points,
// SOC strictly increasing, voltages positive.
func (t *Table) Validate() error {
	if len(t.SOC) != len(t.Voltage) {
		return errors.New("ocv: SOC and Voltage slices must have equal length")
	}
	if len(t.SOC) < 2 {
		return errors.New("ocv: table must have at least 2 points")
	}
	for i := 1; i < len(t.SOC); i++ {
		if t.SOC[i] <= t.SOC[i-1] {
			return fmt.Errorf("ocv: SOC must be strictly increasing at index %d", i)
		}
	}
	for i, v := range t.Voltage {
		if v <= 0 {
			return fmt.Errorf("ocv: voltage must be positive at index %d", i)
		}
	}
	return nil
}

// VoltageAtSOC returns the interpolated OCV for the given SOC percentage.
// SOC is clamped to the table range.
func (t *Table) VoltageAtSOC(soc float64) float64 {
	if soc <= t.SOC[0] {
		return t.Voltage[0]
	}
	if soc >= t.SOC[len(t.SOC)-1] {
		return t.Voltage[len(t.Voltage)-1]
	}
	for i := 1; i < len(t.SOC); i++ {
		if soc <= t.SOC[i] {
			frac := (soc - t.SOC[i-1]) / (t.SOC[i] - t.SOC[i-1])
			return t.Voltage[i-1] + frac*(t.Voltage[i]-t.Voltage[i-1])
		}
	}
	return t.Voltage[len(t.Voltage)-1]
}

// SOCAtVoltage returns the interpolated SOC for the given OCV voltage.
// If the voltage is outside the table range, the nearest boundary SOC is
// returned.
func (t *Table) SOCAtVoltage(voltage float64) float64 {
	if voltage <= t.Voltage[0] {
		return t.SOC[0]
	}
	if voltage >= t.Voltage[len(t.Voltage)-1] {
		return t.SOC[len(t.SOC)-1]
	}
	for i := 1; i < len(t.Voltage); i++ {
		if voltage <= t.Voltage[i] {
			frac := (voltage - t.Voltage[i-1]) / (t.Voltage[i] - t.Voltage[i-1])
			return t.SOC[i-1] + frac*(t.SOC[i]-t.SOC[i-1])
		}
	}
	return t.SOC[len(t.SOC)-1]
}

// Hysteresis models the voltage hysteresis between charge and discharge
// OCV curves. The effective OCV is shifted by +/- halfGap depending on
// the direction of current flow.
func Hysteresis(baseVoltage, halfGap, current float64) float64 {
	if current > 0 {
		return baseVoltage + halfGap
	}
	if current < 0 {
		return baseVoltage - halfGap
	}
	return baseVoltage
}

// RMSE computes the root mean square error between predicted and actual
// voltage series (used for OCV table calibration).
func RMSE(predicted, actual []float64) (float64, error) {
	if len(predicted) != len(actual) {
		return 0, errors.New("ocv: RMSE input lengths differ")
	}
	if len(predicted) == 0 {
		return 0, errors.New("ocv: RMSE needs at least 1 point")
	}
	sum := 0.0
	for i := range predicted {
		d := predicted[i] - actual[i]
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(predicted))), nil
}
