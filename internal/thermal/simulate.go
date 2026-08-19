package thermal

// Simulation models the dynamic temperature evolution of a battery cell
// using a simplified lumped thermal model. It tracks internal heat
// generation and convective cooling to ambient.
type Simulation struct {
	CellTemp    float64 // current cell temperature (C)
	AmbientTemp float64 // environment temperature (C)
	ThermalMass float64 // J/C (thermal capacity of the cell)
	Conductance float64 // W/C (heat transfer coefficient to ambient)
	Params      Params
}

// NewSimulation creates a thermal simulation with initial conditions.
func NewSimulation(cellTemp, ambientTemp, thermalMass, conductance float64) *Simulation {
	return &Simulation{
		CellTemp:    cellTemp,
		AmbientTemp: ambientTemp,
		ThermalMass: thermalMass,
		Conductance: conductance,
		Params:      DefaultParams(),
	}
}

// Step advances the thermal simulation by dtH hours given the current
// flowing through resistance rOhms. It computes heat generation and
// cooling, then updates CellTemp.
func (s *Simulation) Step(current, rOhms, dtH float64) {
	// Heat generated (W) = I^2 * R * resistance_factor
	rFactor := s.Params.ResistanceFactor(s.CellTemp)
	heatW := current * current * rOhms * rFactor

	// Cooling (W) = conductance * (Tcell - Tambient)
	coolingW := s.Conductance * (s.CellTemp - s.AmbientTemp)

	// Net heat (W)
	netHeatW := heatW - coolingW

	// Temperature change: dT = (net_heat * dt_seconds) / thermal_mass
	dtSeconds := dtH * 3600
	if s.ThermalMass > 0 {
		s.CellTemp += (netHeatW * dtSeconds) / s.ThermalMass
	}
}

// IsOverheating returns true if the cell temperature exceeds the safe
// maximum.
func (s *Simulation) IsOverheating() bool {
	return s.CellTemp > s.Params.MaxTemp
}

// IsTooGold returns true if the cell temperature is below operating minimum.
func (s *Simulation) IsTooGold() bool {
	return s.CellTemp < s.Params.MinTemp
}

// Reset sets the cell temperature back to ambient.
func (s *Simulation) Reset() {
	s.CellTemp = s.AmbientTemp
}

// RunProfile simulates temperature evolution over a current profile and
// returns the temperature at each step.
func (s *Simulation) RunProfile(currents []float64, rOhms, dtH float64) []float64 {
	temps := make([]float64, len(currents))
	for i, current := range currents {
		s.Step(current, rOhms, dtH)
		temps[i] = s.CellTemp
	}
	return temps
}
