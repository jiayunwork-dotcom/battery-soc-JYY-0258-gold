package ocv

// HysteresisModel tracks the direction-dependent voltage offset that
// exists between charge and discharge OCV curves in real cells.
type HysteresisModel struct {
	// HalfGap is half the voltage gap between charge and discharge OCV
	// at any given SOC.
	HalfGap float64
	// DecayRate controls how quickly the hysteresis state decays toward
	// equilibrium during rest (per hour).
	DecayRate float64
	// State is the current hysteresis state in [-1, 1] where +1 means
	// fully charging and -1 means fully discharging.
	State float64
}

// DefaultHysteresis returns typical values for a LiFePO4 cell.
func DefaultHysteresis() *HysteresisModel {
	return &HysteresisModel{
		HalfGap:   0.02, // 20 mV half-gap
		DecayRate: 0.5,
		State:     0,
	}
}

// Update adjusts the hysteresis state based on current flow direction
// and elapsed time.
func (h *HysteresisModel) Update(current, dtH float64) {
	if current > 0 {
		// Charging: state moves toward +1
		h.State += (1 - h.State) * h.DecayRate * dtH
	} else if current < 0 {
		// Discharging: state moves toward -1
		h.State += (-1 - h.State) * h.DecayRate * dtH
	} else {
		// Rest: state decays toward 0
		h.State *= (1 - h.DecayRate*dtH)
	}
	// Clamp state
	if h.State > 1 {
		h.State = 1
	}
	if h.State < -1 {
		h.State = -1
	}
}

// VoltageOffset returns the current hysteresis voltage offset to be added
// to the base OCV.
func (h *HysteresisModel) VoltageOffset() float64 {
	return h.State * h.HalfGap
}

// EffectiveOCV returns the base OCV adjusted for hysteresis.
func (h *HysteresisModel) EffectiveOCV(baseOCV float64) float64 {
	return baseOCV + h.VoltageOffset()
}

// Reset sets the hysteresis state to zero (equilibrium).
func (h *HysteresisModel) Reset() {
	h.State = 0
}

// FullChargeState sets state to +1 (fully charged direction).
func (h *HysteresisModel) FullChargeState() {
	h.State = 1
}

// FullDischargeState sets state to -1 (fully discharged direction).
func (h *HysteresisModel) FullDischargeState() {
	h.State = -1
}
