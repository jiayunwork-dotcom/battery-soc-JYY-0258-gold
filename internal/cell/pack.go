package cell

import (
	"fmt"

	"battery-soc/internal/coulomb"
)

// Pack models a series-connected battery pack of identical cells.
type Pack struct {
	Cells  []*Cell
	Series int
}

// NewPack creates a pack with n series-connected cells using the same config.
func NewPack(n int, cfg Config) (*Pack, error) {
	if n < 1 {
		return nil, fmt.Errorf("pack: need at least 1 cell, got %d", n)
	}
	cells := make([]*Cell, n)
	for i := 0; i < n; i++ {
		c, err := New(cfg)
		if err != nil {
			return nil, fmt.Errorf("pack: cell %d: %w", i, err)
		}
		cells[i] = c
	}
	return &Pack{Cells: cells, Series: n}, nil
}

// Step advances all cells by one time step with the same current.
func (p *Pack) Step(current, dtH float64) {
	for _, c := range p.Cells {
		c.Step(current, dtH)
	}
}

// PackVoltage returns the total pack voltage (sum of cell voltages).
func (p *Pack) PackVoltage() float64 {
	total := 0.0
	for _, c := range p.Cells {
		total += c.State().Voltage
	}
	return total
}

// MinSOC returns the lowest SOC among all cells (pack SOC is limited by
// the weakest cell).
func (p *Pack) MinSOC() float64 {
	min := 100.0
	for _, c := range p.Cells {
		if c.State().SOC < min {
			min = c.State().SOC
		}
	}
	return min
}

// MaxSOC returns the highest SOC among all cells.
func (p *Pack) MaxSOC() float64 {
	max := 0.0
	for _, c := range p.Cells {
		if c.State().SOC > max {
			max = c.State().SOC
		}
	}
	return max
}

// Imbalance returns the SOC spread (max - min) across cells.
func (p *Pack) Imbalance() float64 {
	return p.MaxSOC() - p.MinSOC()
}

// RunProfile applies a current profile to all cells and returns the pack
// voltage after each step.
func (p *Pack) RunProfile(samples []coulomb.CurrentSample) []float64 {
	voltages := make([]float64, 0, len(samples))
	for _, s := range samples {
		p.Step(s.Current, s.DT)
		voltages = append(voltages, p.PackVoltage())
	}
	return voltages
}

// SetTemperatures sets different temperatures for each cell (simulates
// thermal gradients in a pack).
func (p *Pack) SetTemperatures(temps []float64) {
	for i, c := range p.Cells {
		if i < len(temps) {
			c.SetTemperature(temps[i])
		}
	}
}
