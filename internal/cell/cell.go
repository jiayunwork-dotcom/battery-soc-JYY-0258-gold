// Package cell models a single battery cell with internal resistance,
// open-circuit voltage, temperature state, and coulomb-counted SOC. It
// ties together the coulomb, OCV, thermal, and SOH modules into one
// coherent simulation entity.
package cell

import (
	"errors"
	"fmt"

	"battery-soc/internal/coulomb"
	"battery-soc/internal/ocv"
	"battery-soc/internal/thermal"
)

// State holds the instantaneous state of a battery cell.
type State struct {
	SOC         float64 // percent [0,100]
	Voltage     float64 // terminal voltage (V)
	Temperature float64 // cell temperature (C)
	Current     float64 // last applied current (A)
	Capacity    float64 // effective capacity (Ah)
}

// Config defines the static parameters of a cell.
type Config struct {
	NominalCapacity float64 // Ah
	InternalR       float64 // Ohms at nominal temperature
	OCVTable        *ocv.Table
	Thermal         thermal.Params
	InitialSOC      float64 // percent
	InitialTemp     float64 // Celsius
}

// DefaultConfig returns a reasonable default for an NMC cell.
func DefaultConfig() Config {
	return Config{
		NominalCapacity: 50.0,
		InternalR:       0.02,
		OCVTable:        ocv.DefaultNMC(),
		Thermal:         thermal.DefaultParams(),
		InitialSOC:      50.0,
		InitialTemp:     25.0,
	}
}

// Cell is the simulation entity.
type Cell struct {
	cfg   Config
	state State
}

// New creates a cell with the given configuration.
func New(cfg Config) (*Cell, error) {
	if cfg.NominalCapacity <= 0 {
		return nil, errors.New("cell: nominal capacity must be positive")
	}
	if cfg.InternalR < 0 {
		return nil, errors.New("cell: internal resistance must be non-negative")
	}
	if cfg.OCVTable == nil {
		return nil, errors.New("cell: OCV table is required")
	}
	if err := cfg.OCVTable.Validate(); err != nil {
		return nil, fmt.Errorf("cell: %w", err)
	}
	c := &Cell{
		cfg: cfg,
		state: State{
			SOC:         cfg.InitialSOC,
			Temperature: cfg.InitialTemp,
			Capacity:    cfg.NominalCapacity,
		},
	}
	c.state.Voltage = cfg.OCVTable.VoltageAtSOC(cfg.InitialSOC)
	return c, nil
}

// State returns a copy of the current cell state.
func (c *Cell) State() State { return c.state }

// Step advances the cell by one time step. dtH is in hours, current in
// Amps (positive = charge). It updates SOC via coulomb counting,
// adjusts capacity for temperature, and computes terminal voltage.
func (c *Cell) Step(current, dtH float64) {
	// Temperature-adjusted effective capacity
	effCap := c.cfg.Thermal.EffectiveCapacity(c.cfg.NominalCapacity, c.state.Temperature)
	if effCap <= 0 {
		effCap = 0.01 // prevent division by zero
	}
	c.state.Capacity = effCap

	// Coulomb counting
	c.state.SOC = coulomb.CoulombSOC(effCap, c.state.SOC, current, dtH)
	c.state.Current = current

	// OCV at current SOC
	ocvV := c.cfg.OCVTable.VoltageAtSOC(c.state.SOC)

	// Terminal voltage = OCV - I*R (discharge) or OCV + I*R (charge sign convention)
	rFactor := c.cfg.Thermal.ResistanceFactor(c.state.Temperature)
	iR := current * c.cfg.InternalR * rFactor
	c.state.Voltage = ocvV - iR // sign: positive current (charge) reduces terminal V shown
}

// SetTemperature updates the cell temperature for the next step.
func (c *Cell) SetTemperature(tempC float64) {
	c.state.Temperature = tempC
}

// RunProfile applies a sequence of current samples and returns the
// state after each step.
func (c *Cell) RunProfile(samples []coulomb.CurrentSample) []State {
	states := make([]State, 0, len(samples))
	for _, s := range samples {
		c.Step(s.Current, s.DT)
		states = append(states, c.state)
	}
	return states
}
