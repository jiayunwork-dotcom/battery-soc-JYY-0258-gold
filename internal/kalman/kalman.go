// Package kalman implements a minimal scalar extended Kalman filter (EKF) for
// SOC estimation. No external libraries or machine learning are used.
package kalman

// fixedGain is the constant correction gain applied toward the measurement z.
const fixedGain = 0.3

// EKF is a scalar extended Kalman filter for SOC.
type EKF struct {
	SOC float64 // current SOC estimate in [0,100]
	P   float64 // scalar covariance estimate
}

// NewEKF initializes the filter with a starting SOC estimate and unit covariance.
func NewEKF(initSOC float64) *EKF {
	return &EKF{SOC: initSOC, P: 1.0}
}

// Step predicts SOC by coulomb counting, then corrects toward the voltage-derived
// measurement z with a fixed gain, clamping to [0,100]. It updates k.SOC and k.P
// and returns the corrected SOC.
func (k *EKF) Step(z, current, dtH, capacityAh float64) float64 {
	pred := k.SOC + (current*dtH/capacityAh)*100
	corrected := pred + fixedGain*(z-pred)
	if corrected < 0 {
		corrected = 0
	}
	if corrected > 100 {
		corrected = 100
	}
	k.P = (1-fixedGain)*k.P + 0.05
	k.SOC = corrected
	return corrected
}
