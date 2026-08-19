// Package profile defines standard charge/discharge test profiles used
// for battery characterisation and SOC algorithm validation. Profiles
// are sequences of current samples that represent real-world usage patterns.
package profile

import "battery-soc/internal/coulomb"

// CC generates a constant-current profile: n steps of dtH hours at the
// given current (positive = charge, negative = discharge).
func CC(current, dtH float64, steps int) []coulomb.CurrentSample {
	samples := make([]coulomb.CurrentSample, steps)
	for i := range samples {
		samples[i] = coulomb.CurrentSample{DT: dtH, Current: current}
	}
	return samples
}

// CCCV generates a CC-CV (constant current then constant voltage) charge
// profile. The CC phase runs at ccCurrent until socThreshold is reached
// (approximated by step count), then the CV phase tapers current linearly
// to zero over cvSteps.
func CCCV(ccCurrent float64, ccSteps, cvSteps int, dtH float64) []coulomb.CurrentSample {
	samples := make([]coulomb.CurrentSample, 0, ccSteps+cvSteps)
	// CC phase
	for i := 0; i < ccSteps; i++ {
		samples = append(samples, coulomb.CurrentSample{DT: dtH, Current: ccCurrent})
	}
	// CV phase: linear taper
	for i := 0; i < cvSteps; i++ {
		frac := 1 - float64(i+1)/float64(cvSteps)
		samples = append(samples, coulomb.CurrentSample{DT: dtH, Current: ccCurrent * frac})
	}
	return samples
}

// PulsedDischarge generates a pulsed discharge profile alternating between
// high current and rest periods.
func PulsedDischarge(highCurrent, dtH float64, pulses int, restSteps int) []coulomb.CurrentSample {
	var samples []coulomb.CurrentSample
	for p := 0; p < pulses; p++ {
		// Discharge pulse
		samples = append(samples, coulomb.CurrentSample{DT: dtH, Current: -highCurrent})
		// Rest
		for r := 0; r < restSteps; r++ {
			samples = append(samples, coulomb.CurrentSample{DT: dtH, Current: 0})
		}
	}
	return samples
}

// DriveCycle generates a simplified drive cycle with acceleration,
// cruise, and regenerative braking phases repeated for the given number
// of cycles.
func DriveCycle(maxCurrent, dtH float64, cycles int) []coulomb.CurrentSample {
	var samples []coulomb.CurrentSample
	for c := 0; c < cycles; c++ {
		// Acceleration (high discharge)
		for i := 0; i < 3; i++ {
			samples = append(samples, coulomb.CurrentSample{DT: dtH, Current: -maxCurrent})
		}
		// Cruise (moderate discharge)
		for i := 0; i < 5; i++ {
			samples = append(samples, coulomb.CurrentSample{DT: dtH, Current: -maxCurrent * 0.3})
		}
		// Regen braking (charge)
		for i := 0; i < 2; i++ {
			samples = append(samples, coulomb.CurrentSample{DT: dtH, Current: maxCurrent * 0.2})
		}
	}
	return samples
}

// Rest generates a rest period with zero current.
func Rest(dtH float64, steps int) []coulomb.CurrentSample {
	samples := make([]coulomb.CurrentSample, steps)
	for i := range samples {
		samples[i] = coulomb.CurrentSample{DT: dtH, Current: 0}
	}
	return samples
}

// Concat joins multiple profiles into one sequence.
func Concat(profiles ...[]coulomb.CurrentSample) []coulomb.CurrentSample {
	var total int
	for _, p := range profiles {
		total += len(p)
	}
	out := make([]coulomb.CurrentSample, 0, total)
	for _, p := range profiles {
		out = append(out, p...)
	}
	return out
}

// TotalCharge computes the net charge transfer in Ah for a profile.
func TotalCharge(samples []coulomb.CurrentSample) float64 {
	total := 0.0
	for _, s := range samples {
		total += s.Current * s.DT
	}
	return total
}

// TotalEnergy estimates energy in Wh given a constant nominal voltage.
func TotalEnergy(samples []coulomb.CurrentSample, nominalV float64) float64 {
	energy := 0.0
	for _, s := range samples {
		energy += s.Current * s.DT * nominalV
	}
	return energy
}

// Duration returns the total time span of the profile in hours.
func Duration(samples []coulomb.CurrentSample) float64 {
	total := 0.0
	for _, s := range samples {
		total += s.DT
	}
	return total
}
