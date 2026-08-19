package kalman

// TunedEKF is an EKF with configurable process noise, measurement noise,
// and gain bounds. It provides more realistic estimation than the basic
// fixed-gain EKF by properly propagating uncertainty.
type TunedEKF struct {
	SOC float64
	P   float64 // estimate covariance
	Q   float64 // process noise variance
	R   float64 // measurement noise variance
}

// TuningParams controls the filter behaviour.
type TuningParams struct {
	InitSOC float64
	InitP   float64
	Q       float64 // process noise (higher = trusts model less)
	R       float64 // measurement noise (higher = trusts measurement less)
}

// DefaultTuning returns conservative tuning for a well-instrumented system.
func DefaultTuning() TuningParams {
	return TuningParams{
		InitSOC: 50.0,
		InitP:   1.0,
		Q:       0.01,
		R:       0.5,
	}
}

// NewTunedEKF creates a tuned Kalman filter.
func NewTunedEKF(params TuningParams) *TunedEKF {
	return &TunedEKF{
		SOC: params.InitSOC,
		P:   params.InitP,
		Q:   params.Q,
		R:   params.R,
	}
}

// Predict performs the time-update step using coulomb counting.
func (f *TunedEKF) Predict(current, dtH, capacityAh float64) {
	// State transition: SOC += (I * dt / Cap) * 100
	f.SOC += (current * dtH / capacityAh) * 100
	// Covariance prediction
	f.P += f.Q
}

// Update performs the measurement-update step given a voltage-derived SOC
// measurement z.
func (f *TunedEKF) Update(z float64) float64 {
	// Kalman gain
	K := f.P / (f.P + f.R)
	// State correction
	f.SOC += K * (z - f.SOC)
	// Covariance update
	f.P = (1 - K) * f.P
	// Clamp
	if f.SOC < 0 {
		f.SOC = 0
	}
	if f.SOC > 100 {
		f.SOC = 100
	}
	return f.SOC
}

// Step combines predict and update in one call.
func (f *TunedEKF) Step(z, current, dtH, capacityAh float64) float64 {
	f.Predict(current, dtH, capacityAh)
	return f.Update(z)
}

// Gain returns the current Kalman gain K = P / (P + R).
func (f *TunedEKF) Gain() float64 {
	return f.P / (f.P + f.R)
}

// Innovation returns the difference between measurement and prediction
// (useful for diagnostics).
func (f *TunedEKF) Innovation(z float64) float64 {
	return z - f.SOC
}
