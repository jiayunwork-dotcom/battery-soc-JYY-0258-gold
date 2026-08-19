package ocv

import "math"

// CalibrationPoint is one measured (SOC, Voltage) pair used to build or
// refine an OCV table.
type CalibrationPoint struct {
	SOC     float64
	Voltage float64
}

// CalibrateTable builds a new OCV table from calibration measurements.
// Points are sorted by SOC and validated.
func CalibrateTable(points []CalibrationPoint) (*Table, error) {
	if len(points) < 2 {
		return nil, ErrTooFewPoints
	}
	// Sort by SOC
	sorted := make([]CalibrationPoint, len(points))
	copy(sorted, points)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].SOC < sorted[i].SOC {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	t := &Table{
		SOC:     make([]float64, len(sorted)),
		Voltage: make([]float64, len(sorted)),
	}
	for i, p := range sorted {
		t.SOC[i] = p.SOC
		t.Voltage[i] = p.Voltage
	}
	if err := t.Validate(); err != nil {
		return nil, err
	}
	return t, nil
}

// ErrTooFewPoints is returned when calibration has fewer than 2 points.
var ErrTooFewPoints = errorString("ocv: need at least 2 calibration points")

type errorString string

func (e errorString) Error() string { return string(e) }

// FitError computes the mean absolute error between a table's predictions
// and measured calibration points.
func FitError(t *Table, points []CalibrationPoint) float64 {
	if len(points) == 0 {
		return 0
	}
	sum := 0.0
	for _, p := range points {
		predicted := t.VoltageAtSOC(p.SOC)
		sum += math.Abs(predicted - p.Voltage)
	}
	return sum / float64(len(points))
}

// ResampleTable creates a new table with n evenly-spaced SOC points by
// interpolating from the source table. Useful for normalizing tables with
// irregular spacing.
func ResampleTable(src *Table, n int) *Table {
	if n < 2 {
		n = 2
	}
	t := &Table{
		SOC:     make([]float64, n),
		Voltage: make([]float64, n),
	}
	step := 100.0 / float64(n-1)
	for i := 0; i < n; i++ {
		soc := float64(i) * step
		t.SOC[i] = soc
		t.Voltage[i] = src.VoltageAtSOC(soc)
	}
	return t
}

// ScaleVoltage creates a new table with all voltages multiplied by a
// factor. Used for aging correction (aged cells have lower OCV).
func ScaleVoltage(src *Table, factor float64) *Table {
	t := &Table{
		SOC:     make([]float64, len(src.SOC)),
		Voltage: make([]float64, len(src.Voltage)),
	}
	copy(t.SOC, src.SOC)
	for i, v := range src.Voltage {
		t.Voltage[i] = v * factor
	}
	return t
}

// Monotonic checks that the OCV table voltages are non-decreasing (a
// physical requirement for most chemistries).
func Monotonic(t *Table) bool {
	for i := 1; i < len(t.Voltage); i++ {
		if t.Voltage[i] < t.Voltage[i-1] {
			return false
		}
	}
	return true
}
